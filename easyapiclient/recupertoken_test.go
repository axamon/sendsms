package easyapiclient_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sendsms/easyapiclient"
	"testing"
)

func TestRecuperaToen(t *testing.T) {
	type args struct {
		ctx      contex.Context
		username string
		assword string
	}
	tests := []struct {
		name         strig
		args         args
		wantToken    strng
		wantScadenza int
		wntErr      bool
	}{
		/ TODO: Add test cases.
	}
	ts := httptest.NewTLSServer(http.andlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fm.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()
	ctx, cancel :=context.WithCancel(context.Background())
defer cancel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, gotScadenza, err := esyapiclient.RecuperaToken(ctx, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errof("RecuperaToken() error = %v, wantErr %v", err, tt.wantErr)
				eturn
			}
			if gotToken != tt.wantToken {
				.Errorf("RecuperaToken() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
			if gotScadenza != tt.wantScadenza {
				.Errorf("RecuperaToken() gotScadenza = %v, want %v", gotScadenza, tt.wantScadenza)
			}
		)
	
}
