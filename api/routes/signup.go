package routes

import (
	"encoding/json"
	"net/http"

	psCrypto "PasswordServer2/lib/crypto"
	psDatabase "PasswordServer2/lib/database"
)

func PostSignup(w http.ResponseWriter, r *http.Request) {
	signupParameters := struct {
		Email                string
		MasterHash           []byte
		ProtectedDatabaseKey psDatabase.ProtectedDatabaseKey
		Config               psDatabase.ClientConfig
	}{}

	decoderError := json.NewDecoder(r.Body).Decode(&signupParameters)

	// Create json encoder that will be used to write response
	jsonResponse := json.NewEncoder(w)

	// Errors relating to validation / items not being present
	validationErrors := []map[string]any{}

	if decoderError != nil {
		validationErrors = append(validationErrors, map[string]any{"Code": "00", "Message": "Unable to decode JSON body."})
	}

	if signupParameters.Email == "" {
		validationErrors = append(validationErrors, map[string]any{"Code": "01", "Message": "Required parameter 'Email' not provided."})
	}

	if signupParameters.MasterHash == nil {
		validationErrors = append(validationErrors, map[string]any{"Code": "02", "Message": "Required parameter 'MasterHash' not provided."})
	}

	if signupParameters.ProtectedDatabaseKey.Iv == nil {
		validationErrors = append(validationErrors, map[string]any{"Code": "03", "Message": "Required parameter 'ProtectedDatabaseKey.Iv' not provided."})
	}

	if signupParameters.ProtectedDatabaseKey.Key == nil {
		validationErrors = append(validationErrors, map[string]any{"Code": "04", "Message": "Required parameter 'ProtectedDatabaseKey.Key' not provided."})
	}

	var emptyConfig psDatabase.ClientConfig
	if signupParameters.Config == emptyConfig {
		validationErrors = append(validationErrors, map[string]any{"Code": "05", "Message": "Required parameter 'Config' not provided."})
	}

	if psDatabase.UserEmailInUse(signupParameters.Email) {
		validationErrors = append(validationErrors, map[string]any{"Code": "06", "Message": "There is an existing user with the email provided."})
	}

	if len(validationErrors) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse.Encode(map[string]any{"Error(s)": validationErrors})
		return
	}

	// Use pbkdf2 to strengthen keys with a random salt
	strengthenedMasterHashSalt := psCrypto.RandomBytes(16)
	strengthenedMasterHashBytes := psCrypto.StrengthenMasterHash(signupParameters.MasterHash, strengthenedMasterHashSalt)

	newMasterHash := psDatabase.NewMasterHash()
	newMasterHash.Hash = strengthenedMasterHashBytes
	newMasterHash.Salt = strengthenedMasterHashSalt

	newProtectedDatabaseKey := psDatabase.NewProtectedDatabaseKey()
	newProtectedDatabaseKey.Key = signupParameters.ProtectedDatabaseKey.Key
	newProtectedDatabaseKey.Iv = signupParameters.ProtectedDatabaseKey.Iv

	newStoredConfig := psDatabase.NewConfigProfile()
	newStoredConfig.MasterHash = newMasterHash
	newStoredConfig.ProtectedDatabaseKey = newProtectedDatabaseKey
	newStoredConfig.Config = signupParameters.Config

	newUser := psDatabase.NewUser()
	newUser.Email = signupParameters.Email
	newUser.ConfigProfiles = []psDatabase.ConfigProfile{newStoredConfig}

	userId := psDatabase.CreateUser(newUser)

	w.WriteHeader(http.StatusOK)
	jsonResponse.Encode(map[string]any{"UserId": userId.Hex()})
}
