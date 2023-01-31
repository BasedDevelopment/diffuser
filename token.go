package main

import (
	"bytes"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

func getToken() {
	log.Info().Msg("Getting token")
	reqUrl := k.String("api.url") + "/login"

	v := url.Values{}
	v.Set("username", k.String("api.username"))
	v.Set("password", k.String("api.password"))

	r, _ := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBufferString(v.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(r)
	if err != nil {
		log.Error().Err(err).Msg("Error getting token")
		return
	}

	if resp.Status != "302 Found" {
		log.Error().Msg("Error getting token")
		return
	}

	for _, cookie := range resp.Cookies() {
		token = cookie.Value
	}

	if token == "" {
		log.Error().Msg("Error getting token")
		return
	}

	log.Info().Msg("Got token")
}
