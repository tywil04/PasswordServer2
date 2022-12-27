package api

import (
	"PasswordServer2/api/routes"
)

// ALl paths are relative to /api/v1. This means a path of /test actually means /api/v1/test.
var Routes map[string]MethodMap = map[string]MethodMap{
	"/auth/signup": {
		Post: routes.PostSignup,
	},
	"/auth/signin": {
		Post: routes.PostSignin,
	},
	"/user/configprofiles/new": {
		Post: routes.PostSignin,
	},
}
