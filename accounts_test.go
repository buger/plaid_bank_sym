package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestAccountsList(t *testing.T) {
	clearDB()
	user, pass := createUser()

	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()

	s, _ := NewUserSession(user.Id, pass, "127.0.0.1")

	client := &http.Client{}
	req, _ := http.NewRequest("GET", ts.URL+"/accounts", nil)
	req.Header.Add("X-Auth-Token", s.Token)
	req.Header.Add("X-Auth-User", user.Id)
	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var data UserAccountsResp
	json.Unmarshal(body, &data)

	// Should first create account with user
	if len(data.Accounts) != 1 {
		t.Errorf("Should create first account when user created %v", data)
	}

	if data.Accounts[0].UserId != user.Id {
		t.Errorf("Should create account for proper user")
	}

	if data.Accounts[0].Ballance != 0 {
		t.Errorf("Should create account 0 ballance")
	}
}

func TestAccountsListAuthError(t *testing.T) {
	clearDB()
	user, pass := createUser()

	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()

	NewUserSession(user.Id, pass, "127.0.0.1")

	client := &http.Client{}
	req, _ := http.NewRequest("GET", ts.URL+"/accounts", nil)
	req.Header.Add("X-Auth-Token", "wrong token")
	req.Header.Add("X-Auth-User", user.Id)
	resp, _ := client.Do(req)

	if resp.StatusCode == 200 {
		t.Errorf("Should not process request with wrong token")
	}
}

func TestAccountsListAuthError2(t *testing.T) {
	clearDB()
	user, pass := createUser()

	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()

	s, _ := NewUserSession(user.Id, pass, "127.0.0.1")

	client := &http.Client{}
	req, _ := http.NewRequest("GET", ts.URL+"/accounts", nil)
	req.Header.Add("X-Auth-Token", s.Token)
	req.Header.Add("X-Auth-User", "wrong")
	resp, _ := client.Do(req)

	if resp.StatusCode == 200 {
		t.Errorf("Should not process request with wrong user id")
	}
}

func createAccount(t *testing.T) (*User, string, *UserSession, string) {
	user, pass := createUser()

	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()

	s, _ := NewUserSession(user.Id, pass, "127.0.0.1")

	client := &http.Client{}
	req, _ := http.NewRequest("POST", ts.URL+"/account/create", nil)
	req.Header.Add("X-Auth-Token", s.Token)
	req.Header.Add("X-Auth-User", user.Id)
	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Should create account %v", resp, body)
	}

	var data map[string]string
	json.Unmarshal(body, &data)

	return user, pass, s, data["id"]
}

func TestAccountCreate(t *testing.T) {
	clearDB()
	createAccount(t)
}

func TestAccountInfo(t *testing.T) {
	clearDB()

	user, pass, s, accId := createAccount(t)

	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()

	v := url.Values{}
	v.Set("user", user.Id)
	v.Set("password", pass)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", ts.URL+"/account?id="+accId, nil)
	req.Header.Add("X-Auth-Token", s.Token)
	req.Header.Add("X-Auth-User", user.Id)
	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Should create account %v", resp, string(body))
	}

	var acc Account
	json.Unmarshal(body, &acc)

	if acc.UserId != user.Id {
		t.Errorf("Should found user account")
	}
}
