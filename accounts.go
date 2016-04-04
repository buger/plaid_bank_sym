package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Account struct {
	Id       string
	UserId   string
	Ballance float64
}

func (a *Account) Transactions() (transactions []Transaction) {
	transactionIds, _ := db.AllKeys(fmt.Sprintf("accounts/transactions/%s", a.Id))

	for _, id := range transactionIds {
		var t Transaction
		db.ReadJSON(fmt.Sprintf("accounts/transactions/%s/%s", a.Id, id), &t)
		transactions = append(transactions, t)
	}

	return
}

type UserAccountsResp struct {
	Accounts []Account `json:"accounts"`
}

var userAccountsHandler = requireAuthHandler(func(user *User, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&UserAccountsResp{user.Accounts()})
})

type AccountCreatedResp struct {
	Id string `json:"id"`
}

var createAccountHandler = requireAuthHandler(func(user *User, w http.ResponseWriter, r *http.Request) {
	acc := CurrentBank.CreateAccount(user)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&AccountCreatedResp{acc.Id})
})

type ErrorResp struct {
	Error string `json:"error,omniempty"`
}

var accountInfoHandler = requireAuthHandler(func(user *User, w http.ResponseWriter, r *http.Request) {
	accId := r.URL.Query().Get("id")

	var acc Account
	if err := db.ReadJSON(fmt.Sprintf("users/%s/accounts/%s", user.Id, accId), &acc); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&ErrorResp{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&acc)
})

var transactionsHandler = requireAuthHandler(func(user *User, w http.ResponseWriter, r *http.Request) {
	accId := r.URL.Query().Get("id")

	var acc Account
	if err := db.ReadJSON(fmt.Sprintf("users/%s/accounts/%s", user.Id, accId), &acc); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(&ErrorResp{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(acc.Transactions())
})

func init() {
	http.HandleFunc("/account/create", createAccountHandler)
	http.HandleFunc("/account/transactions", transactionsHandler)
	http.HandleFunc("/accounts", userAccountsHandler)
	http.HandleFunc("/account", accountInfoHandler)
}
