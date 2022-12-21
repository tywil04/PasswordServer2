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
	validationError := ""

	if decoderError != nil {
		validationError = "Unable to decode JSON body."
	} else if signupParameters["Email"] == "" {
		validationError = "Required parameter 'Email' not provided."
	} else if signupParameters["MasterHash"] == "" {
		validationError = "Required parameter 'MasterHash' not provided."
	} else if signupParameters["ProtectedDatabaseKey"] == "" {
		validationError = "Required parameter 'ProtectedDatabaseKey' not provided."
	} else if decodingMasterHashBytesError != nil {
		validationError = "Unable to decode hex encoded parameter 'MasterHash'."
	} else if len(protectedDatabaseParts) != 2 {
		validationError = "Unable to split parameter 'ProtectedDatabaseKey' into its IV and key."
	} else if decodingProtectedDatabaseKeyError != nil {
		validationError = "Unable to decode hex encoded key from split parameter 'ProtectedDatabaseKey'."
	} else if decodingProtectedDatabaseKeyIVError != nil {
		validationError = "Unable to decode hex encoded IV from split parameter 'ProtectedDatabaseKey'."
	}

	if validationError != "" {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse.Encode(map[string]any{"Error": validationError})
		return
	}

	// Use pbkdf2 to strengthen keys with a random salt
	strengthenedMasterHashSalt := psCrypto.RandomBytes(16)
	strengthenedMasterHashBytes := psCrypto.StrengthenMasterHash(masterHashBytes, strengthenedMasterHashSalt)

	newUser := psDatabase.User{
		Email:                  signupParameters["Email"],
		MasterHash:             strengthenedMasterHashBytes,
		MasterHashSalt:         strengthenedMasterHashSalt,
		ProtectedDatabaseKey:   decodedProtectedDatabaseKey,
		ProtectedDatabaseKeyIV: decodedProtectedDatabaseKeyIV,
	}
	psDatabase.Database.Create(&newUser)

	w.WriteHeader(http.StatusOK)
	jsonResponse.Encode(map[string]any{"UserId": newUser.Id})
}
