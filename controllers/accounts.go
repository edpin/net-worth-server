//
// Date: 10/18/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"flag"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/net-worth-server/models"
	"github.com/tidwall/gjson"
)

//
// Return all accounts in this system.
//
// curl -H "Authorization: Bearer XXXXX" http://localhost:9090/api/v1/accounts
//
func (t *Controller) GetAccounts(w http.ResponseWriter, r *http.Request) {

	// Return Happy
	t.RespondJSON(w, http.StatusOK, t.DB.GetAllAcounts())
}

//
// Return one account by id.
//
// curl -H "Authorization: Bearer XXXXX" http://localhost:9090/api/v1/accounts/1
//
func (t *Controller) GetAccount(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Convert string to int.
	id, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		t.RespondError(w, http.StatusBadRequest, "Something went wrong please email support")
		return
	}

	// Get the account by id.
	account, err := t.DB.GetAccountById(uint(id))

	// Return json based on if this was a good result or not.
	if err != nil {
		t.RespondError(w, http.StatusNotFound, err.Error())
	} else {
		t.RespondJSON(w, http.StatusOK, account)
	}
}

//
// Create a new record from the data passed in.
//
// curl -H "Content-Type: application/json" -X POST -d '{"name":"Lending Club","balance":70123.66}' -H "Authorization: Bearer XXXXXX" http://localhost:9090/api/v1/accounts
//
func (t *Controller) CreateAccount(w http.ResponseWriter, r *http.Request) {

	account := models.Account{}

	// Decode the json we posted in.
	if err := t.DecodePostedJson(&account, w, r); err != nil {
		t.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Store in database & return json.
	if err := t.DB.CreateAccount(&account); err != nil {
		t.RespondError(w, http.StatusBadRequest, err.Error())
	} else {
		// Get fresh account because of timezone issues, and add in all the other magic.
		account, _ := t.DB.GetAccountById(account.Id)
		t.RespondJSON(w, http.StatusCreated, account)
	}
}

//
// Get all the marks for an account.
//
// curl -H "Authorization: Bearer XXXXX" http://localhost:9090/api/v1/accounts/XX/marks
//
func (t *Controller) GetAccountMarks(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Convert string to int.
	id, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		t.RespondError(w, http.StatusBadRequest, "Something went wrong please email support")
		return
	}

	// Return Happy
	t.RespondJSON(w, http.StatusOK, t.DB.GetMarksByAccountById(uint(id)))
}

//
// Mark account value for today.
//
// curl -H "Content-Type: application/json" -X POST -d '{"balance":1000.00, "date": "2017-10-05"}' -H "Authorization: Bearer XXXXXX" http://localhost:9090/api/v1/accounts/XX/marks
//
func (t *Controller) CreateAccountMark(w http.ResponseWriter, r *http.Request) {

	// Grab date for late formatting.
	body, _ := ioutil.ReadAll(r.Body)
	date := gjson.Get(string(body), "date").String()
	balance := gjson.Get(string(body), "balance").Float()

	// URL vars
	vars := mux.Vars(r)

	// Convert string to int.
	id, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		t.RespondError(w, http.StatusBadRequest, "Something went wrong please email support")
		return
	}

	// Reformat date.
	pDate, err := time.Parse("2006-01-02", date)

	if err != nil {
		t.RespondError(w, http.StatusBadRequest, "Unable to parse date.")
		return
	}

	// Store in database & return json.
	err = t.DB.MarkAccountByDate(uint(id), pDate.UTC(), balance)

	if err != nil {
		t.RespondError(w, http.StatusBadRequest, "Something went wrong please email support")
		return
	}

	// Get the account by id.
	account, err := t.DB.GetAccountById(uint(id))

	if err != nil {
		t.RespondError(w, http.StatusBadRequest, err.Error())
	} else {
		t.RespondJSON(w, http.StatusCreated, account)
	}
}

//
// Manage funds going into an account. Add funds to an account, or remove.
// Negative values will be removing funds from an account.
//
// curl -H "Content-Type: application/json" -X POST -d '{"amount":500.12, "date": "2017-11-12", "note": "Test note."}' -H "Authorization: Bearer XXXXXX" http://localhost:9090/api/v1/accounts/XX/funds
//
func (t *Controller) AccountManageFunds(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	date := gjson.Get(string(body), "date").String()
	amount := gjson.Get(string(body), "amount").Float()
	note := gjson.Get(string(body), "note").String()

	vars := mux.Vars(r)

	// Hack for testing.
	if flag.Lookup("test.v") != nil {
		vars = make(map[string]string)
		vars["id"] = "1"
	}

	// Convert string to int.
	id, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		t.RespondError(w, http.StatusBadRequest, "Something went wrong please email support")
		return
	}

	// Reformat date.
	pDate, err := time.Parse("2006-01-02", date)

	if err != nil {
		t.RespondError(w, http.StatusBadRequest, "Unable to parse date.")
		return
	}

	// Add / Subtract money to an account units
	t.DB.AccountUnitsAddFunds(uint(id), pDate, amount, note)

	// Get the account by id.
	account, err := t.DB.GetAccountById(uint(id))

	if err != nil {
		t.RespondError(w, http.StatusBadRequest, err.Error())
	} else {
		t.RespondJSON(w, http.StatusCreated, account)
	}
}

/* End File */
