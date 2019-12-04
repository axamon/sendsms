// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package easyapiclient permette un facile utlizzo
// delle API Easyapi di TIM.
package easyapiclient

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// TokenResponse contiene le risposte di EasyApi.
type TokenResponse struct {
	Token     string `json:"access_token"`
	Scope     string `json:"scope"`
	Tokentype string `json:"token_type"`
	Scadenza  int    `json:"expires_in"`
}

const easyapiGetTokenURL = "https://easyapi.telecomitalia.it:8248/token"

// RecuperaToken restituisce un token Easyapi valido.
func RecuperaToken(ctx context.Context, username, password string) (token string, err error) {

	// Recupera le credenziali per Easyapi.
	credenziali := username + ":" + password

	// Encoda le credenziali per passarle come http header.
	authenticator := base64.StdEncoding.EncodeToString([]byte(credenziali))

	// Accetta anche certificati https non validi.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Crea il cliet http.
	client := &http.Client{Transport: tr}

	// Corpo della chiamata web.
	body := strings.NewReader(`grant_type=client_credentials`)

	// Crea la request da inviare.
	req, err := http.NewRequestWithContext(ctx, "POST", easyapiGetTokenURL, body)
	if err != nil {
		err = fmt.Errorf("Errore creazione request: %v, %v", req, err)
	}

	// Aggiunge alla request l'autenticazione.
	req.Header.Set("Authorization", "Basic "+authenticator)

	// Aggiunge alla request gli header per passare le informazioni.
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "easyapi.telecomitalia.it:8248")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Length", "29")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	// Invia la request HTTP.
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("Errore: %v", err.Error())
	}

	// Se la http response ha un codice di errore esce.
	if resp.StatusCode > 299 {
		err = fmt.Errorf("Impossibile recuperare token http statuscode: %d", resp.StatusCode)
	}

	// Legge il body della risposta.
	bodyresp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("Error Impossibile leggere risposta client http: %s", err.Error())
	}

	// Come da specifica chiude il body della response.
	defer resp.Body.Close()

	// Crea variabile per archiviare i risulati.
	tokeninfo := new(TokenResponse)

	// Effettua l'unmashalling dei dati nella variabile.
	err = json.Unmarshal(bodyresp, &tokeninfo)
	if err != nil {
		err = fmt.Errorf("Errore nella scomposizione del json: %s", err.Error())
	}

	return tokeninfo.Token, err
}
