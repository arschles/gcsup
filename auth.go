package main

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

func getAuthenticatedClient(ctx context.Context, jwtFileLocation string) (*http.Client, error) {
	data, err := ioutil.ReadFile(jwtFileLocation)
	if err != nil {
		return nil, err
	}

	jwtConf, err := google.JWTConfigFromJSON(data, storage.DevstorageFullControlScope)
	if err != nil {
		return nil, err
	}
	return jwtConf.Client(ctx), nil
}
