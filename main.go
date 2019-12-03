// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sendsms/easyapiclient"

	"github.com/tkanos/gonfig"
)

// conf è una istanza del type Configuration con username e password da usare
// per accedere a Easyapi
var conf Configuration

// confFile è il file json in cui sono scritte username e password da usare.
var confFile string

// cell è il numero di cellulare a cui inviare SMS.
var cell string

// cell è il messaggio da inviare.
var messaggio string

var author bool

func main() {
	// Crea il contesto iniziale e la funzione cancel per uscire
	// dal programma in modo pulito.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.BoolVar(&author, "a", false, "Per segnalizioni all'autore")
	flag.StringVar(&confFile, "file", "conf.json", "File di configurazione")
	flag.StringVar(&cell, "c", "", "Cellulare a cui inviare SMS")
	flag.StringVar(&messaggio, "m", "", "Messaggio  da inviare")

	// Parsa i parametri non di default passati all'avvio.
	flag.Parse()

	if author {
		fmt.Println("Apire una segnalazione su: ")
		fmt.Println("https://scm.code.telecomitalia.it/00246506/sendsms/issues")
		os.Exit(0)
	}

	// Recupera valori dal file json di configurazione passato come argomento.
	err := gonfig.GetConf(confFile, &conf)
	if err != nil {
		log.Fatalf("Errore Impossibile recuperare informazioni dal file di configurazione: %s", confFile)
	}

	// Recupera un token per inviare sms valido.
	token, err := easyapiclient.RecuperaToken(ctx, conf.UsernameEasyaPi, conf.Password)
	if err != nil {
		log.Fatalf("Errore nel recupero del token sms: %s\n", err.Error())
	}

	// Recupera lo shortnumber da usare per inviare sms.
	shortnumber, err := Info(ctx, token)
	if err != nil {
		log.Fatalf("Errore, impossibile recuperare shortnumber %s\n", err.Error())
	}

	// InviaSms invia sms usando le informazioni recuperate in precedenza.
	err = InviaSms(ctx, token, shortnumber, cell, messaggio)
	if err != nil {
		log.Fatalf("Errore, sms non inviato: %s\n", err)
	}

	// Termina correttamente.
	os.Exit(0)
}
