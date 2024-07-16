package handlers

import (
	"errors"
	"fmt"
	"github.com/a-h/templ"
	"github.com/gorilla/schema"
	"net/http"
	"quizzly/pkg/logger"
)

type (
	Handler[T any] interface {
		Handle(writer http.ResponseWriter, request *http.Request, in T) (templ.Component, error)
	}

	BadRequestErr struct {
		originalErr error
	}
)

func BadRequest(err error) *BadRequestErr {
	return &BadRequestErr{
		originalErr: err,
	}
}

func (err *BadRequestErr) Error() string {
	return err.originalErr.Error()
}

func Templ[T any](handler Handler[T], log logger.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		inStruct, err := parseIn[T](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("parse request error", err)
			return
		}

		component, err := handler.Handle(w, r, inStruct)
		var badRequestErr *BadRequestErr
		if errors.As(err, &badRequestErr) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("handle request error", err)
			return
		}

		templHandler := templ.Handler(component)
		templHandler.ServeHTTP(w, r)
	}
}

func parseIn[T any](r *http.Request) (T, error) {
	var inStruct T
	switch r.Method {
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			return inStruct, fmt.Errorf("parse form error: %v", err)
		}
		err = schema.NewDecoder().Decode(&inStruct, r.Form)
		if err != nil {
			return inStruct, fmt.Errorf("decode request body error: %v", err)
		}
	case http.MethodGet, http.MethodDelete:
		err := schema.NewDecoder().Decode(&inStruct, r.URL.Query())
		if err != nil {
			return inStruct, fmt.Errorf("decode url query error: %v", err)
		}
	default:
		return inStruct, fmt.Errorf("unsupported method: %v", r.Method)
	}

	return inStruct, nil
}
