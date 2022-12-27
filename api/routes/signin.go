package routes

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"net/http"

	psCrypto "PasswordServer2/lib/crypto"
	psDatabase "PasswordServer2/lib/database"

	"go.mongodb.org/mongo-driver/bson"
)

func PostSignin(w http.ResponseWriter, r *http.Request) {
	signinParameters := struct {
		Email      string
		MasterHash []byte
		Config     psDatabase.ClientConfig
	}{}

	decoderError := json.NewDecoder(r.Body).Decode(&signinParameters)

	// Create json encoder that will be used to write response
	jsonResponse := json.NewEncoder(w)

	// Errors relating to validation / items not being present
	validationErrors := []map[string]any{}

	if decoderError != nil {
		validationErrors = append(validationErrors, map[string]any{"Code": "10", "Message": "Unable to decode JSON body."})
	}

	if signinParameters.Email == "" {
		validationErrors = append(validationErrors, map[string]any{"Code": "11", "Message": "Required parameter 'Email' not provided."})
	}

	if signinParameters.MasterHash == nil {
		validationErrors = append(validationErrors, map[string]any{"Code": "12", "Message": "Required parameter 'MasterHash' not provided."})
	}

	var emptyConfig psDatabase.ClientConfig
	if signinParameters.Config == emptyConfig {
		validationErrors = append(validationErrors, map[string]any{"Code": "13", "Message": "Required parameter 'Config' not provided."})
	}

	if !psDatabase.UserEmailInUse(signinParameters.Email) {
		validationErrors = append(validationErrors, map[string]any{"Code": "14", "Message": "There is no existing user with the email provided."})
	}

	if len(validationErrors) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse.Encode(map[string]any{"Error(s)": validationErrors})
		return
	}

	user := psDatabase.FindUserViaEmail(signinParameters.Email)
	userModel := psDatabase.ConvertPrimitiveUserToUserModel(user)

	configFound := false
	foundCredential := psDatabase.NewCredential()
	configs := []psDatabase.ClientConfig{}
	query, _ := psDatabase.ClientConfigs.Find(context.TODO(), bson.M{"masterkey": 1})
	query.All(context.TODO(), &configs)

	for _, credential := range userModel.Credentials {
		for _, config := range configs {
			if config == signinParameters.Config {
				configFound = true
			}

			if credential.ClientConfigId == config.Id {
				foundCredential = credential
			}
		}
	}

	if !configFound {
		psDatabase.CreateClientConfig(signinParameters.Config)

		existingConfig := psDatabase.NewClientConfig()
		psDatabase.ClientConfigs.FindOne(context.TODO(), bson.M{"masterkey": 1}).Decode(&existingConfig)

		w.WriteHeader(http.StatusOK)
		jsonResponse.Encode(map[string]any{
			"Authenticated":          false,
			"NewCredentialForConfig": true,
			"ExistingConfig":         existingConfig,
		})
		return
	}

	strengthenedMasterHashBytes := psCrypto.StrengthenMasterHash(signinParameters.MasterHash, foundCredential.MasterHash.Salt)
	same := subtle.ConstantTimeCompare(foundCredential.MasterHash.Hash, strengthenedMasterHashBytes) == 1

	if same {
		cookieError := psCrypto.CreateSessionCookie(w, user)

		if cookieError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			jsonResponse.Encode(map[string]any{})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse.Encode(map[string]any{
		"Authenticated":          same,
		"NewCredentialForConfig": false,
		"UserId":                 userModel.Id.Hex(),
		"ProtectedDatabaseKey":   foundCredential.ProtectedDatabaseKey,
	})
}
