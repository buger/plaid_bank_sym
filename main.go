package main

import (
	"fmt"
	"net/http"
)

var db *FileDB
var CurrentBank Bank

func Init(bankId int, dataPath string) {
	// Limit context to current bank
	db = OpenDB(fmt.Sprintf(dataPath+"/bank/%d", bankId))

	var banks = []Bank{Bank{1, "Bank of America"}, Bank{2, "Wells Fargo"}}
	for _, bank := range banks {
		db.WriteJSON(fmt.Sprintf("banks/%d", bank.Id), bank)
	}

	if err := db.ReadJSON(fmt.Sprintf("banks/%d", bankId), &CurrentBank); err != nil {
		panic(fmt.Sprintf("Can't initialize bank: %d %v", bankId, err))
	}
}

func main() {
	Init(1, ".data")
	http.ListenAndServe(":80", nil)
}
