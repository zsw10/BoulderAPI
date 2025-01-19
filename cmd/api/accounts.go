package main

import (
	"crypto"
	"net/http"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

type User struct {
	Email         string
	Resgistration *registration.Resource
	key           crypto.PrivateKey
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u User) GetRegistration() *registration.Resource {
	return u.Resgistration
}

func (u *User) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func (app *application) createAccountHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	err := app.readJSON(w, r, user.Email)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	privateKey, err := certcrypto.GeneratePrivateKey(certcrypto.RSA2048)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user.key = privateKey

	config := lego.NewConfig(&user)
	config.CADirURL = app.config.boulder.url
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user.Resgistration = reg

	res := envelope{"AccountID": reg.URI, "Status": "registered", "CreatedAt": time.Now().Format(time.RFC3339)}

	err = app.writeJSON(w, http.StatusCreated, res, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
