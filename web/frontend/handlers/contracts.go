package handlers

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/gorilla/schema"
	"net/http"
)

type (
	Handler[T any] interface {
		Handle(writer http.ResponseWriter, request *http.Request, in T) (templ.Component, error)
	}
)

func Wrapper[T any](handler Handler[T]) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var inStruct T
		switch r.Method {
		case http.MethodPost:
			err := r.ParseForm()
			if err != nil {
				http.Error(w, fmt.Errorf("parse form error: %v", err).Error(), http.StatusInternalServerError)
				return
			}
			err = schema.NewDecoder().Decode(&inStruct, r.Form)
			if err != nil {
				http.Error(w, fmt.Errorf("decode request body error: %v", err).Error(), http.StatusInternalServerError)
				return
			}
		case http.MethodGet:
			err := schema.NewDecoder().Decode(&inStruct, r.URL.Query())
			if err != nil {
				http.Error(w, fmt.Errorf("decode url query error: %v", err).Error(), http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, "unsupported method", http.StatusInternalServerError)
			return
		}

		component, err := handler.Handle(w, r, inStruct)
		if err != nil {
			http.Error(w, fmt.Errorf("component handlers error: %v", err).Error(), http.StatusInternalServerError)
			return
		}

		templHandler := templ.Handler(component)
		templHandler.ServeHTTP(w, r)
	}
}
