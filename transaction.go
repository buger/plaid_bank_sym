package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type OperationType int

const (
	Transfer = OperationType(iota)
	Withdraw
	Deposit
)

type TransactionState int

const (
	TStarted = TransactionState(iota)
	TConfirmed
	TCompleted
)

type Transaction struct {
	Id            string
	BankFrom      int64
	AccountFrom   string
	BankTo        int64
	AccountTo     string
	Amount        float64
	State         TransactionState
	OperationType OperationType
}

func LogTransaction(op OperationType, accountFrom, accountTo string, amount float64) *Transaction {
	t := &Transaction{
		Id:            randStringBytes(50),
		BankFrom:      CurrentBank.Id,
		BankTo:        CurrentBank.Id,
		AccountFrom:   accountFrom,
		AccountTo:     accountTo,
		Amount:        amount,
		OperationType: op,
	}

	switch op {
	case Transfer:
		// In case of transfer we are creating 2 transactions
		db.WriteJSON(fmt.Sprintf("accounts/transactions/%s/%s", accountFrom, t.Id), t)
		db.WriteJSON(fmt.Sprintf("accounts/transactions/%s/%s", accountTo, t.Id), t)
	case Deposit:
		db.WriteJSON(fmt.Sprintf("accounts/transactions/%s/%s", accountTo, t.Id), t)
	case Withdraw:
		db.WriteJSON(fmt.Sprintf("accounts/transactions/%s/%s", accountFrom, t.Id), t)
	}

	return t
}

var transferTransactionHandler = requireAuthHandler(func(user *User, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	db.Transaction(func(db *FileDB) {
		from := r.FormValue("from")
		to := r.FormValue("to")
		amountStr := r.FormValue("amount")
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResp{Error: "Amount is not a number"})
			return
		}

		var accFrom, accTo Account
		if err := db.ReadJSON(fmt.Sprintf("accounts/%s", from), &accFrom); err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(&ErrorResp{Error: "From account not found"})
			return
		}

		if err := db.ReadJSON(fmt.Sprintf("accounts/%s", to), &accTo); err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(&ErrorResp{Error: "To account not found"})
			return
		}

		if accFrom.UserId != user.Id {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResp{Error: "From account do not belong to the logged user"})
			return
		}

		if accFrom.Ballance < amount {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResp{Error: "From account does not have enough money"})
			return
		}

		accFrom.Ballance -= amount
		accTo.Ballance += amount

		// Update accounts data
		db.WriteJSON(fmt.Sprintf("accounts/%s", accFrom.Id), &accFrom)
		db.WriteJSON(fmt.Sprintf("accounts/%s", accTo.Id), &accTo)

		LogTransaction(Transfer, accFrom.Id, accTo.Id, amount)

		w.WriteHeader(http.StatusOK)
	})
})

var depositTransactionHandler = requireAuthHandler(func(user *User, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	db.Transaction(func(db *FileDB) {
		to := r.FormValue("to")
		amountStr := r.FormValue("amount")

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResp{Error: "Amount is not a number"})
			return
		}

		var accTo Account
		if err := db.ReadJSON(fmt.Sprintf("accounts/%s", to), &accTo); err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(&ErrorResp{Error: "To account not found"})
			return
		}

		if accTo.UserId != user.Id {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResp{Error: "To account do not belong to the logged user"})
			return
		}

		accTo.Ballance += amount

		db.WriteJSON(fmt.Sprintf("accounts/%s", accTo.Id), &accTo)
		LogTransaction(Deposit, "", accTo.Id, amount)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	})
})

var withdrawTransactionHandler = requireAuthHandler(func(user *User, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	db.Transaction(func(db *FileDB) {
		from := r.FormValue("from")
		amountStr := r.FormValue("amount")

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResp{Error: "Amount is not a number"})
			return
		}

		var accFrom Account
		if err := db.ReadJSON(fmt.Sprintf("accounts/%s", from), &accFrom); err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(&ErrorResp{Error: "From account not found"})
			return
		}

		if accFrom.UserId != user.Id {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResp{Error: "From account do not belong to the logged user"})
			return
		}

		if accFrom.Ballance < amount {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&ErrorResp{Error: "From account does not have enough money"})
			return
		}

		accFrom.Ballance -= amount

		db.WriteJSON(fmt.Sprintf("accounts/%s", accFrom.Id), &accFrom)
		LogTransaction(Withdraw, accFrom.Id, "", amount)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	})
})

func init() {
	http.HandleFunc("/transfer", transferTransactionHandler)
	http.HandleFunc("/deposit", depositTransactionHandler)
	http.HandleFunc("/withdraw", withdrawTransactionHandler)
}
