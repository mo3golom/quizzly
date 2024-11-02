package link

import (
	"fmt"
	"net/http"
	"quizzly/pkg/variables"
)

type linkBuilder struct {
	baseLink string
	host     string
	useHTTPS bool
}

func newLinkBuilder(baseLink string) linkBuilder {
	return linkBuilder{baseLink: baseLink}
}

func (b linkBuilder) addHost(request ...*http.Request) linkBuilder {
	if len(request) > 0 {
		b.host = request[0].Host
	}

	return b
}

func (b linkBuilder) addHTTPS(variablesRepo variables.Repository) linkBuilder {
	b.useHTTPS = variablesRepo.GetString(variables.AppEnvironmentVariable) == string(variables.EnvironmentProd)

	return b
}

func (b linkBuilder) build() string {
	result := b.baseLink
	if b.host == "" {
		return result
	}

	result = fmt.Sprintf("%s%s", b.host, result)

	if b.useHTTPS {
		return fmt.Sprintf("https://%s", result)
	}

	return fmt.Sprintf("http://%s", result)
}
