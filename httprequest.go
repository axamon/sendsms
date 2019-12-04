// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)


func httpRequest(ctx context.Context, url, method, token string, payload []byte) (result []byte, err error) {
	
	// Accetta anche certificati https non validi.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Crea il cliet http.
	client := &http.Client{Transport: tr}

	// Crea la request da inviare.
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(payload))
	if err != nil {
		err = fmt.Errorf("Errore creazione request: %v", req)
	}

	// Aggiunge alla request l'autenticazione.
	req.Header.Set("Authorization", "Bearer "+token)

	if method == "POST" {
		req.Header.Set("Content-Type", "application/xml")
	}

	// Aggiunge alla request gli header per passare le informazioni.
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Invia la request HTTP.
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("Errore %v", err.Error())
	}

	// Se la http response ha un codice di errore esce.
	if resp.StatusCode > 299 {
		err = fmt.Errorf("Errore %d", resp.StatusCode)
	}

	// Legge il body della risposta.
	bodyresp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("Impossibile leggere risposta client http: %v", err.Error())
	}

	// Come da specifiche va chiuso il body.
	defer resp.Body.Close()

return bodyresp, err
}