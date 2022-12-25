package crypto

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	psDatabase "PasswordServer2/lib/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionCookie struct {
	UserId         string
	SessionTokenId int
}

func CreateSessionCookie(response http.ResponseWriter, user primitive.M) error {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	publicKey := &privateKey.PublicKey

	cookieExpiry := time.Now().Add(time.Hour)

	sessionToken := psDatabase.NewSessionToken()
	sessionToken.N = publicKey.N.Bytes()
	sessionToken.E = publicKey.E
	sessionToken.Expiry = cookieExpiry

	sessionTokenId := psDatabase.InsertSessionTokenIntoUser(user, sessionToken)

	sessionCookie := SessionCookie{
		UserId:         user["_id"].(primitive.ObjectID).Hex(),
		SessionTokenId: sessionTokenId,
	}

	jsonPayload := new(bytes.Buffer)
	json.NewEncoder(jsonPayload).Encode(sessionCookie)
	hashed := sha512.Sum512(jsonPayload.Bytes())

	signature, _ := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA512, hashed[:])

	encodedSignature := hex.EncodeToString(signature)
	encodedSessionCookie := hex.EncodeToString(jsonPayload.Bytes())

	cookie := http.Cookie{
		Name:     "SessionToken",
		Value:    encodedSessionCookie + "," + encodedSignature,
		Expires:  cookieExpiry,
		Secure:   os.Getenv("ENVIRONMENT") == "production",
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(response, &cookie)

	return nil
}

func VerifySessionCookie(request *http.Request) (bool, primitive.M, primitive.M, error) {
	cookie, cookieError := request.Cookie("SessionToken")

	if cookie == nil || cookie.Value == "" {
		return false, primitive.M{}, primitive.M{}, nil
	}

	if cookieError != nil {
		return false, primitive.M{}, primitive.M{}, cookieError
	}

	splitValue := strings.Split(cookie.Value, ",")
	jsonSessionCookie, _ := hex.DecodeString(splitValue[0])

	signature, _ := hex.DecodeString(splitValue[1])

	sessionCookie := SessionCookie{}
	json.NewDecoder(bytes.NewBuffer(jsonSessionCookie)).Decode(&sessionCookie)

	userId, _ := primitive.ObjectIDFromHex(sessionCookie.UserId)
	user := psDatabase.FindUserViaId(userId)

	sessionTokens := user["sessiontokens"].(primitive.A)
	sessionToken := sessionTokens[sessionCookie.SessionTokenId].(primitive.M)

	publicKey := rsa.PublicKey{
		N: new(big.Int).SetBytes(sessionToken["n"].([]byte)),
		E: sessionToken["e"].(int),
	}

	jsonPayload := new(bytes.Buffer)
	json.NewEncoder(jsonPayload).Encode(sessionCookie)
	hashed := sha512.Sum512(jsonPayload.Bytes())

	for i := 0; i < len(sessionTokens); i++ {
		// if the session has expired remove the token. we are lax about this because any cookie which is expired doesn't isnt valid which means the code cant even process it
		sessionToken := sessionTokens[i].(primitive.M)
		if sessionToken["expiry"].(time.Time).Before(time.Now()) {
			psDatabase.RemoveSessionTokenViaIdFromUser(user, i)
		}
	}

	if rsa.VerifyPKCS1v15(&publicKey, crypto.SHA512, hashed[:], signature) == nil {
		return true, user, sessionToken, nil
	}

	return false, primitive.M{}, primitive.M{}, nil
}

func ClearSessionCookie(response http.ResponseWriter, request *http.Request) (bool, error) {
	authenticated, user, sessionToken, _ := VerifySessionCookie(request)

	if authenticated {
		cookie := http.Cookie{
			Name:     "SessionToken",
			Value:    "",
			MaxAge:   -1,
			Secure:   os.Getenv("ENVIRONMENT") == "production",
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
		}

		http.SetCookie(response, &cookie)

		psDatabase.RemoveSessionTokenFromUser(user, sessionToken)

		return true, nil
	}

	return false, nil
}
