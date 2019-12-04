// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"regexp"
	"time"
)

const easyapiMTURL = "https://easyapi.telecomitalia.it:8248/sms/v1/mt"

// inviaSms invia un sms al destinatario.
func inviaSms(ctx context.Context, token, shortnumber, cell, message string) error {
	// Modificato il contesto impostando un timout.
	ctxInvio, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	errori := make(chan error, 1)
	var err error

	select {
	case err = <-errori:
		return err

	// Se la funzione impiega più del tempo previsto dal timeout.
	case <-ctxInvio.Done():
		errori <- ctxInvio.Err()  // "context deadline exceeded"

	default:
		// modifica cell aggiungendo quanto richiesto da easyapi e il formato internazionale +39
		address := "tel:+39" + cell

		// effetuttua verifiche sui formati dati passati.
		err = verificheFormali(address, message, token)
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

		// Effettua il marshalling dei campi sms in []byte.
		bodyreq, err := xml.Marshal(nuovoSMS)
		if err != nil {
			errori <- fmt.Errorf("Impossibile parsare dati in xml: %s", err.Error())
		}

		_, err = httpRequest(ctxInvio, easyapiMTURL, "POST", token, bodyreq)
		if err != nil {
			errori <- fmt.Errorf("Errore Richiesta http fallita: %v", err.Error())
		}

	}

	return err
}


// isCell è il formato internazionale italiano dei cellulari.
var isCell = regexp.MustCompile(`(?m)^tel:\+39\d{9,12}$`)

// isToken è il formato che deve avere un token easyapi ben formattato.
var isToken = regexp.MustCompile(`(?m)[0-9a-z]{8,8}-[0-9a-z]{4,4}-[0-9a-z]{4,4}-[0-9a-z]{4,4}-[0-9a-z]{12,12}`)

func verificheFormali(address, message, token string) error {

	switch {
	// Verifica il formato del cellulare.
	case !isCell.MatchString(address):
		return fmt.Errorf("Formato del numero di cellulare non corretto: %s", cell)

		// Verifica che il messsaggio non super 160 caratteri.
	case len(message) > 160:
		return fmt.Errorf("Messaggio troppo lungo: %d caratteri. Max 160 caratteri ammessi", len(message))
		// Verifica che il token sia nel formato corretto.
	case !isToken.MatchString(token):
		return fmt.Errorf("Token non nel formato standard")
	}

	return nil
}
