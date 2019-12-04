// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/xml"
	"fmt"
)

const easyapiGetInfo = "https://easyapi.telecomitalia.it:8248/sms/v1/info"

// Info recupera lo shortnumber da usare per inviare sms.
func Info(ctx context.Context, token string) (shortnumber string, err error) {


	bodyresp, err := httpRequest(ctx, easyapiGetInfo, "GET", token, nil)
	if err != nil {
		return "", fmt.Errorf("Errore Richiesta http fallita: %v", err.Error())
	}

	sNum := new(ShortNum)

	err = xml.Unmarshal(bodyresp, &sNum)
	if err != nil {
		return "", fmt.Errorf("Error Impossibile effettuare caricamento shortnumber: %v", err.Error())
	}

	return sNum.Number, err

}
