package routes

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"

	psCrypto "PasswordServer2/lib/crypto"
	psDatabase "PasswordServer2/lib/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	oldConfigs := make([]psDatabase.ClientConfig, len(userModel.ConfigProfiles))
	foundStoredConfig := psDatabase.ConfigProfile{}

	for index, profile := range userModel.ConfigProfiles {
		if profile.Config == signinParameters.Config {
			configFound = true
			foundStoredConfig = profile
			oldConfigs[index] = profile.Config
		}
	}

	if !configFound {
		w.WriteHeader(http.StatusOK)
		jsonResponse.Encode(map[string]any{
			"Authenticated":     false,
			"NewConfigRequired": true,
			"OldConfigs":        oldConfigs,
		})
		return
	}

	strengthenedMasterHashBytes := psCrypto.StrengthenMasterHash(signinParameters.MasterHash, foundStoredConfig.MasterHash.Salt)
	same := subtle.ConstantTimeCompare(foundStoredConfig.MasterHash.Hash, strengthenedMasterHashBytes) == 1

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
		"Authenticated":        same,
		"UserId":               user["_id"].(primitive.ObjectID).Hex(),
		"ProtectedDatabaseKey": foundStoredConfig.ProtectedDatabaseKey,
	})
}
