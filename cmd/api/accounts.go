package main

import (
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/zsw10/BoulderAPI/internal/data"
)

func (app *application) createAccountHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	privateKey, err := certcrypto.GeneratePrivateKey(certcrypto.RSA2048)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user := data.User{
		Email: input.Email,
		Key:   privateKey,
	}

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

	user.CreatedAt = time.Now()
	user.Status = reg.Body.Status
	id, err := strconv.Atoi(path.Base(reg.URI))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user.ID = id

	err = app.models.User.Insert(&user)

	res := envelope{"AccountID": user.ID, "Status": user.Status, "CreatedAt": user.CreatedAt.Format(time.RFC3339)}

	err = app.writeJSON(w, http.StatusCreated, res, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
