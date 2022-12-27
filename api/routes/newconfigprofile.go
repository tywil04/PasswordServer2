package routes

import (
	"encoding/json"
	"net/http"

	psCrypto "PasswordServer2/lib/crypto"
	psDatabase "PasswordServer2/lib/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PostNewConfigProfile(w http.ResponseWriter, r *http.Request) {
	newConfigProfileParameters := struct {
		UserId               string
		MasterHash           []byte
		Config               psDatabase.ClientConfig
		ProtectedDatabaseKey psDatabase.ProtectedDatabaseKey
	}{}

	decoderError := json.NewDecoder(r.Body).Decode(&newConfigProfileParameters)

	// Create json encoder that will be used to write response
	jsonResponse := json.NewEncoder(w)

	// Errors relating to validation / items not being present
	validationErrors := []map[string]any{}

	if decoderError != nil {
		validationErrors = append(validationErrors, map[string]any{"Code": "10", "Message": "Unable to decode JSON body."})
	}

	if newConfigProfileParameters.UserId == "" {
		validationErrors = append(validationErrors, map[string]any{"Code": "11", "Message": "Required parameter 'UserId' not provided."})
	}

	if newConfigProfileParameters.MasterHash == nil {
		validationErrors = append(validationErrors, map[string]any{"Code": "12", "Message": "Required parameter 'MasterHash' not provided."})
	}

	var emptyConfig psDatabase.ClientConfig
	if newConfigProfileParameters.Config == emptyConfig {
		validationErrors = append(validationErrors, map[string]any{"Code": "13", "Message": "Required parameter 'Config' not provided."})
	}

	userId, _ := primitive.ObjectIDFromHex(newConfigProfileParameters.UserId)
	user := psDatabase.FindUserViaId(userId)
	if user == nil {
		validationErrors = append(validationErrors, map[string]any{"Code": "14", "Message": "There is no existing user with the userId provided."})
	}

	if len(validationErrors) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse.Encode(map[string]any{"Error(s)": validationErrors})
		return
	}

	// Use pbkdf2 to strengthen keys with a random salt
	strengthenedMasterHashSalt := psCrypto.RandomBytes(16)
	strengthenedMasterHashBytes := psCrypto.StrengthenMasterHash(newConfigProfileParameters.MasterHash, strengthenedMasterHashSalt)

	newMasterHash := psDatabase.NewMasterHash()
	newMasterHash.Hash = strengthenedMasterHashBytes
	newMasterHash.Salt = strengthenedMasterHashSalt

	newProtectedDatabaseKey := psDatabase.NewProtectedDatabaseKey()
	newProtectedDatabaseKey.Key = newConfigProfileParameters.ProtectedDatabaseKey.Key
	newProtectedDatabaseKey.Iv = newConfigProfileParameters.ProtectedDatabaseKey.Iv

	newStoredConfig := psDatabase.NewConfigProfile()
	newStoredConfig.MasterHash = newMasterHash
	newStoredConfig.ProtectedDatabaseKey = newProtectedDatabaseKey
	newStoredConfig.Config = newConfigProfileParameters.Config

	psDatabase.InsertConfigProfileIntoUser(user, newStoredConfig)

	w.WriteHeader(http.StatusOK)
	jsonResponse.Encode(map[string]any{"ConfigProfileCreated": true})
}
