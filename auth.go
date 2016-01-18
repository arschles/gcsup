package main

import (
	"io/ioutil"
	"net/http"

	storage "code.google.com/p/google-api-go-client/storage/v1"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
)

func getAuthenticatedClient(ctx context.Context, jwtFileLocation string) (*http.Client, error) {
	data, err := ioutil.ReadFile(jwtFileLocation)
	if err != nil {
		return nil, err
	}

	jwtConf, err := google.JWTConfigFromJSON(data, storage.DevstorageFull_controlScope)
	if err != nil {
		return nil, err
	}
	return jwtConf.Client(ctx), nil
}
