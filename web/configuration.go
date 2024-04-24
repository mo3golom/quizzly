package web

import (
	"fmt"
	"net/http"
	"os"
)

const (
	publicPath = "web/frontend/public"
)

func Configure() {
	_, err := os.Stat("./web/frontend/public")
	if os.IsNotExist(err) {
		panic(fmt.Sprintf("Directory '%s' not found.\n", "web"))
	}

	http.Handle("/", http.FileServer(http.Dir(publicPath)))
	routes()
}

func routes() {

}
