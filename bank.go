package main

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Bank struct {
	Id   int64
	Name string
}

type User struct {
	Id           string
	Name         string
	DateOfBirth  time.Time
	PasswordHash [sha512.Size]byte
}

func (u *User) Accounts() (accounts []Account) {
	accountIds, _ := db.AllKeys(fmt.Sprintf("users/%s/accounts", u.Id))

	for _, id := range accountIds {
		var acc Account
		db.ReadJSON(fmt.Sprintf("users/%s/accounts/%s", u.Id, id), &acc)
		accounts = append(accounts, acc)
	}

	return
}

var durYear = time.Hour * 24 * 365

type UserAgeErr struct {
	age time.Duration
}

func (e *UserAgeErr) Error() string {
	return fmt.Sprintf("Person age should be between 18 and 100. Current person age: %d", e.age)
}

type UserExistsErr struct {
	name        string
	dateOfBirth time.Time
}

func (e *UserExistsErr) Error() string {
	return fmt.Sprintf("User with name %s and birth date %s already exists", e.name, e.dateOfBirth)
}

func (b *Bank) CreateUser(name string, dateOfBirth time.Time) (user *User, password string, err error) {
	age := time.Now().Sub(dateOfBirth) / durYear
	if age < 18 || age > 100 {
		err = &UserAgeErr{age}
		return
	}

	uIdx := fmt.Sprintf("users/index/name-bdate/%s-%d", name, dateOfBirth.Unix())

	if db.Exists(uIdx) {
		err = &UserExistsErr{name, dateOfBirth}
		return
	}

	// Choose unique user id
	id := randStringBytes(10)
	for db.Exists(fmt.Sprintf("users/%s.json", id)) {
		id = randStringBytes(10)
	}
	password = randStringBytes(10)

	user = &User{
		Id:           id,
		Name:         name,
		DateOfBirth:  dateOfBirth,
		PasswordHash: sha512.Sum512([]byte(password)),
	}

	uKey := fmt.Sprintf("users/%s.json", user.Id)
	db.WriteJSON(uKey, user)
	db.CreateLink(uKey, uIdx)

	b.CreateAccount(user)

	return
}

func (b *Bank) CreateAccount(user *User) *Account {
	// Choose unique account id
	id := randStringBytes(34)
	for db.Exists(fmt.Sprintf("accounts/%s", id)) {
		id = randStringBytes(34)
	}

	acc := &Account{
		Id:       id,
		UserId:   user.Id,
		Ballance: 0.0,
	}

	accKey := fmt.Sprintf("accounts/%s", acc.Id)
	db.WriteJSON(accKey, acc)
	db.CreateLink(accKey, fmt.Sprintf("users/%s/accounts/%s", user.Id, acc.Id))

	return acc
}

type UserCreatedResp struct {
	Id       string `json:"id,omitempty"`
	Password string `json:"password,omitempty"`
	Error    string `json:"error,omitempty"`
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	name := r.FormValue("name")
	date, err := time.Parse(time.RFC3339, r.FormValue("dateOfBirth"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&UserCreatedResp{Error: err.Error()})
		return
	}

	user, pass, err := CurrentBank.CreateUser(name, date)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&UserCreatedResp{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&UserCreatedResp{Id: user.Id, Password: pass})
}

func init() {
	http.HandleFunc("/user/create", createUserHandler)
}
