package main

import (
	"fmt"
	"net/http"
)

type jwtRoundTripper struct {
	token []byte
	cl    *http.Client
}

func (j jwtRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authentication", fmt.Sprintf("Bearer %s", string(j.token)))
	return j.cl.Do(req)
}
