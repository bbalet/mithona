package main

/*  mithona allows to simply share information on a server with its users
    Copyright (C) 2013 Benjamin BALET

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.*/

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"net/http"
)

// loginActionHandler is the HTTP handler for the login action
// If wrong credentials, then redirect to loginFormHandler with Flash message
func loginActionHandler(w http.ResponseWriter, r *http.Request) {

	// Read the private key
	pemData, err := ioutil.ReadFile(PRIVATE_KEY_FILE)
	checkHttpError(err, w)

	// Extract the PEM-encoded data block
	block, _ := pem.Decode(pemData)
	if block == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Decode the RSA private key
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	checkHttpError(err, w)

	// Decode the Base64 into binary
	cipheredValue, err := base64.StdEncoding.DecodeString(r.FormValue("CipheredValue"))
	checkHttpError(err, w)

	// Decrypt the data
	var out []byte
	out, err = rsa.DecryptPKCS1v15(rand.Reader, priv, cipheredValue)
	checkHttpError(err, w)

	//home or login
	session, _ := store.Get(r, "goServerView")
	h := sha512.New()
	h.Write(out)
	str := base64.StdEncoding.EncodeToString(h.Sum([]byte{}))

	if str == config.Password {
		session.Values["isConnected"] = true
		session.Values["user"] = r.FormValue("User")
		session.Save(r, w)
		homeHandler(w, r)
	} else {
		session.Values["isConnected"] = false
		session.AddFlash("Either the username or the password provided is wrong")
		session.Save(r, w)
		loginFormHandler(w, r)
	}
}

// loginFormHandler is the HTTP handler for the login form
func loginFormHandler(w http.ResponseWriter, r *http.Request) {
	// Read the public key
	pemData, err := ioutil.ReadFile(PUBLIC_KEY_FILE)
	checkHttpError(err, w)
	var p = pageContent(r, "Login", pemData)
	err = templates["login"].ExecuteTemplate(w, "base", p)
	checkHttpError(err, w)
}

// logoutHandler closes the admin session
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "goServerView")
	session.Values["isConnected"] = false
	session.Save(r, w)
	homeHandler(w, r)
}
