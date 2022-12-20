package api

import (
	"net/http"
	"strings"
)

type MethodMap struct {
	Get     func(response http.ResponseWriter, request *http.Request)
	Post    func(response http.ResponseWriter, request *http.Request)
	Patch   func(response http.ResponseWriter, request *http.Request)
	Put     func(response http.ResponseWriter, request *http.Request)
	Delete  func(response http.ResponseWriter, request *http.Request)
	Trace   func(response http.ResponseWriter, request *http.Request)
	Options func(response http.ResponseWriter, request *http.Request)
	Connect func(response http.ResponseWriter, request *http.Request)
	Head    func(response http.ResponseWriter, request *http.Request)
}

func APIHandler(path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, path)
		_, ok := Routes[path]

		if ok {
			switch r.Method {
			case http.MethodGet:
				if Routes[path].Get != nil {
					Routes[path].Get(w, r)
				}
			case http.MethodPost:
				if Routes[path].Post != nil {
					Routes[path].Post(w, r)
				}
			case http.MethodPatch:
				if Routes[path].Patch != nil {
					Routes[path].Patch(w, r)
				}
			case http.MethodPut:
				if Routes[path].Put != nil {
					Routes[path].Put(w, r)
				}
			case http.MethodDelete:
				if Routes[path].Delete != nil {
					Routes[path].Delete(w, r)
				}
			case http.MethodTrace:
				if Routes[path].Trace != nil {
					Routes[path].Trace(w, r)
				}
			case http.MethodOptions:
				if Routes[path].Options != nil {
					Routes[path].Options(w, r)
				}
			case http.MethodConnect:
				if Routes[path].Connect != nil {
					Routes[path].Connect(w, r)
				}
			case http.MethodHead:
				if Routes[path].Head != nil {
					Routes[path].Head(w, r)
				}
			}
		}
	})
}
