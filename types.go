// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Configuration tiene gli elementi di configurazione
type Configuration struct {
	UsernameEasyaPi string `json:"username"`
	Password        string `json:"password"`
}

// ShortNum è il numero breve da usare.
type ShortNum struct {
	Number string `xml:"shortNumber"`
}

type sms struct {
	Address  string `xml:"address"`
	Msgid    string `xml:"msgid"`
	Notify   string `xml:"notify"`
	Validity string `xml:"validity"`
	Oadc     string `xml:"oadc"`
	Message  string `xml:"message"`
}
