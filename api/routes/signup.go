package routes

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	psCrypto "PasswordServer2/lib/crypto"
	psDatabase "PasswordServer2/lib/database"
)

func PostSignup(w http.ResponseWriter, r *http.Request) {
	signupParameters := map[string]string{}
	decoderError := json.NewDecoder(r.Body).Decode(&signupParameters)

	// Create json encoder that will be used to write response
	jsonResponse := json.NewEncoder(w)

	// Decoding stuff so validation errors can be checked
	masterHashBytes, decodingMasterHashBytesError := hex.DecodeString(signupParameters["MasterHash"])
	protectedDatabaseParts := strings.Split(signupParameters["ProtectedDatabaseKey"], ";")

	var decodedProtectedDatabaseKey, decodedProtectedDatabaseKeyIV []byte
	var decodingProtectedDatabaseKeyError, decodingProtectedDatabaseKeyIVError error

	if len(protectedDatabaseParts) == 2 {
		decodedProtectedDatabaseKey, decodingProtectedDatabaseKeyError = hex.DecodeString(protectedDatabaseParts[1])
		decodedProtectedDatabaseKeyIV, decodingProtectedDatabaseKeyIVError = hex.DecodeString(protectedDatabaseParts[0])
	}

	// Errors relating to validation / items not being present
	validationErrors := []map[string]any{}

	if decoderError != nil {
		validationErrors = append(validationErrors, map[string]any{"Code": "00", "Message": "Unable to decode JSON body."})
	}

	if signupParameters["Email"] == "" {
		validationErrors = append(validationErrors, map[string]any{"Code": "01", "Message": "Required parameter 'Email' not provided."})
	}

	if signupParameters["MasterHash"] == "" {
		validationErrors = append(validationErrors, map[string]any{"Code": "02", "Message": "Required parameter 'MasterHash' not provided."})
	}

	if signupParameters["ProtectedDatabaseKey"] == "" {
		validationErrors = append(validationErrors, map[string]any{"Code": "03", "Message": "Required parameter 'ProtectedDatabaseKey' not provided."})
	}

	if decodingMasterHashBytesError != nil {
		validationErrors = append(validationErrors, map[string]any{"Code": "04", "Message": "Unable to decode hex encoded parameter 'MasterHash'."})
	}

	if len(protectedDatabaseParts) != 2 {
		validationErrors = append(validationErrors, map[string]any{"Code": "05", "Message": "Unable to split parameter 'ProtectedDatabaseKey' into its IV and key."})
	}

	if decodingProtectedDatabaseKeyError != nil {
		validationErrors = append(validationErrors, map[string]any{"Code": "06", "Message": "Unable to decode hex encoded key from split parameter 'ProtectedDatabaseKey'."})
	}

	if decodingProtectedDatabaseKeyIVError != nil {
		validationErrors = append(validationErrors, map[string]any{"Code": "07", "Message": "Unable to decode hex encoded IV from split parameter 'ProtectedDatabaseKey'."})
	}

	if len(validationErrors) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse.Encode(map[string]any{"Error(s)": validationErrors})
		return
	}

	// Use pbkdf2 to strengthen keys with a random salt
	strengthenedMasterHashSalt := psCrypto.RandomBytes(16)
	strengthenedMasterHashBytes := psCrypto.StrengthenMasterHash(masterHashBytes, strengthenedMasterHashSalt)

	newMasterHash := psDatabase.NewMasterHash()
	newMasterHash.MasterHash = strengthenedMasterHashBytes
	newMasterHash.Salt = strengthenedMasterHashSalt

	newProtectedDatabaseKey := psDatabase.NewProtectedDatabaseKey()
	newProtectedDatabaseKey.ProtectedDatabaseKey = decodedProtectedDatabaseKey
	newProtectedDatabaseKey.Iv = decodedProtectedDatabaseKeyIV

	newUser := psDatabase.NewUser()
	newUser.Email = signupParameters["Email"]
	newUser.MasterHashes = []psDatabase.MasterHash{newMasterHash}
	newUser.ProtectedDatabaseKeys = []psDatabase.ProtectedDatabaseKey{newProtectedDatabaseKey}

	userId := psDatabase.CreateUser(newUser)

	w.WriteHeader(http.StatusOK)
	jsonResponse.Encode(map[string]any{"UserId": userId.Hex()})
}
