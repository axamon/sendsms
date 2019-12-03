// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

const easyapiMTURL = "https://easyapi.telecomitalia.it:8248/sms/v1/mt"

// inviaSms invia un sms al destinatario.
func inviaSms(ctx context.Context, token, shortnumber, cell, message string) error {
	// Modificato il contesto impostando un timout.
	ctxInvio, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	select {
	// Se la funzione impiega più del tempo previsto dal timeout esce con errore.
	case <-ctxInvio.Done():
		return ctxInvio.Err() // "context deadline exceeded"

	default:
		// modifica cell aggiungendo quanto richiesto da easyapi e il formato internazionale +39
		address := "tel:+39" + cell

		// effetuttua verifiche sui formati dati passati.
		err := verificheFormali(address, message, token)
		if err != nil {
			return err
		}

		// Crea la struttura necessaria per un nuovo sms.
		nuovoSMS := new(sms)

		nuovoSMS.Address = address
		nuovoSMS.Msgid = "9938"
		nuovoSMS.Notify = "Y"
		nuovoSMS.Validity = "00:03"
		nuovoSMS.Oadc = shortnumber
		nuovoSMS.Message = message

		//fmt.Println(nuovoSMS)

		// Effettua il marshalling dei campi sms in []byte.
		bodyreq, err := xml.Marshal(nuovoSMS)
		if err != nil {
			return fmt.Errorf("Impossibile parsare dati in xml: %s", err.Error())
		}

		// Accetta anche certificati https non validi.
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		// Crea il cliet http.
		client := &http.Client{Transport: tr}

		// Crea la request da inviare.
		req, err := http.NewRequestWithContext(ctxInvio, "POST", easyapiMTURL, bytes.NewBuffer(bodyreq))
		if err != nil {
			return fmt.Errorf("Errore creazione request: %v: %s", req, err.Error())
		}

		// fmt.Println(req)

		// Aggiunge alla request l'autenticazione.
		req.Header.Set("Authorization", "Bearer "+token)

		// Aggiunge alla request gli header per passare le informazioni.
		req.Header.Set("Content-Type", "application/xml")

		// Invia la request HTTP.
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("Errore nella richiesta http %s", err.Error())
		}

		// Body va chiuso come da specifica.
		defer resp.Body.Close()

		// Se la http response ha un codice di errore esce.
		if resp.StatusCode > 299 {
			return fmt.Errorf("Impossibile inviare sms. http statuscode: %d", resp.StatusCode)
		}

		// Legge il body della risposta.
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Impossibile leggere risposta. client http: %s", err.Error())
		}

		// 	fmt.Println(string(bodyresp))
		return err
	}
}

// isCell è il formato internazionale italiano dei cellulari.
var isCell = regexp.MustCompile(`(?m)^tel:\+39\d{9,12}$`)

// isToken è il formato che deve avere un token easyapi ben formattato.
var isToken = regexp.MustCompile(`(?m)[0-9a-z]{8,8}-[0-9a-z]{4,4}-[0-9a-z]{4,4}-[0-9a-z]{4,4}-[0-9a-z]{12,12}`)

func verificheFormali(address, message, token string) (err error) {

	switch {
	// Verifica il formato del cellulare.
	case !isCell.MatchString(address):
		err = fmt.Errorf("Formato del numero di cellulare non corretto: %s", cell)

		// Verifica che il messsaggio non super 160 caratteri.
	case len(message) > 160:
		err = fmt.Errorf("Messaggio troppo lungo: %d caratteri. Max 160 caratteri ammessi", len(message))
		// Verifica che il token sia nel formato corretto.
	case !isToken.MatchString(token):
		err = fmt.Errorf("Token non nel formato standard")
	}

	return err
}
