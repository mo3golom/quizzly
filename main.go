package main

import (
	"fmt"
	"github.com/a-h/templ"
	"net/http"
	"quizzly/web"
	frontend "quizzly/web/frontend/templ"
)

func main() {
	web.Configure()

	component := frontend.HeaderComponent("John")

	http.Handle("/John", templ.Handler(component))

	fmt.Println("Listening on :3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
