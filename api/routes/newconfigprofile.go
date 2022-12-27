package routes

import (
	"encoding/json"
	"net/http"

	psCrypto "PasswordServer2/lib/crypto"
	psDatabase "PasswordServer2/lib/database"
)

func PostNewConfigProfile(w http.ResponseWriter, r *http.Request) {
	authenticated, user, _, _ := psCrypto.VerifySessionCookie(r)

	if authenticated {
		newConfigProfileParameters := struct {
			MasterHash           []byte
			ProtectedDatabaseKey psDatabase.ProtectedDatabaseKey
			Config               psDatabase.ClientConfig
		}{}

		decoderError := json.NewDecoder(r.Body).Decode(&newConfigProfileParameters)

		// Create json encoder that will be used to write response
		jsonResponse := json.NewEncoder(w)

		// Errors relating to validation / items not being present
		validationErrors := []map[string]any{}

		if decoderError != nil {
			validationErrors = append(validationErrors, map[string]any{"Code": "10", "Message": "Unable to decode JSON body."})
		}

		if newConfigProfileParameters.MasterHash == nil {
			validationErrors = append(validationErrors, map[string]any{"Code": "12", "Message": "Required parameter 'MasterHash' not provided."})
		}

		var emptyConfig psDatabase.ClientConfig
		if newConfigProfileParameters.Config == emptyConfig {
			validationErrors = append(validationErrors, map[string]any{"Code": "13", "Message": "Required parameter 'Config' not provided."})
		}

		if user == nil {
			validationErrors = append(validationErrors, map[string]any{"Code": "14", "Message": "There is no existing user with the userId provided."})
		}

		if len(validationErrors) != 0 {
			w.WriteHeader(http.StatusBadRequest)
			jsonResponse.Encode(map[string]any{"Error(s)": validationErrors})
			return
		}

		configId := psDatabase.CreateClientConfig(newConfigProfileParameters.Config)

		// Use pbkdf2 to strengthen keys with a random salt
		strengthenedMasterHashSalt := psCrypto.RandomBytes(16)
		strengthenedMasterHashBytes := psCrypto.StrengthenMasterHash(newConfigProfileParameters.MasterHash, strengthenedMasterHashSalt)

		newMasterHash := psDatabase.NewMasterHash()
		newMasterHash.Hash = strengthenedMasterHashBytes
		newMasterHash.Salt = strengthenedMasterHashSalt

		newProtectedDatabaseKey := psDatabase.NewProtectedDatabaseKey()
		newProtectedDatabaseKey.Key = newConfigProfileParameters.ProtectedDatabaseKey.Key
		newProtectedDatabaseKey.Iv = newConfigProfileParameters.ProtectedDatabaseKey.Iv

		newCredential := psDatabase.NewCredential()
		newCredential.MasterHash = newMasterHash
		newCredential.ProtectedDatabaseKey = newProtectedDatabaseKey
		newCredential.ClientConfigId = configId

		psDatabase.InsertCredentialIntoUser(user, newCredential)

		w.WriteHeader(http.StatusOK)
		jsonResponse.Encode(map[string]any{"ConfigProfileCreated": true})
	} else {
		w.WriteHeader(http.StatusNetworkAuthenticationRequired)
	}
}
