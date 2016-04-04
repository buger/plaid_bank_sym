package main

import (
    "testing"
    // "encoding/json"
    "net/url"
    "net/http/httptest"
    // "io/ioutil"
    "net/http"
    "strings"
)

func TestTransactionDeposit(t *testing.T) {
    clearDB()
    user, pass := createUser()

    ts := httptest.NewServer(http.DefaultServeMux)
    defer ts.Close()

    s, _ := NewUserSession(user.Id, pass, "127.0.0.1")

    accounts := user.Accounts()

    v := url.Values{}
    v.Set("to", accounts[0].Id)
    v.Set("amount", "1000.1")

    client := &http.Client{}
    req, _ := http.NewRequest("POST", ts.URL + "/deposit", strings.NewReader(v.Encode()))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Add("X-Auth-Token", s.Token)
    req.Header.Add("X-Auth-User", user.Id)
    resp, _ := client.Do(req)

    // body, _ := ioutil.ReadAll(resp.Body)
    // resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("Should process transaction")
    }

    accounts = user.Accounts()
    if accounts[0].Ballance != 1000.1 {
        t.Errorf("Should update ballance")
    }
}

func TestTransactionDepositWrongAcc(t *testing.T) {
    clearDB()
    user, pass := createUser()

    ts := httptest.NewServer(http.DefaultServeMux)
    defer ts.Close()

    s, _ := NewUserSession(user.Id, pass, "127.0.0.1")

    accounts := user.Accounts()

    v := url.Values{}
    v.Set("to", "123")
    v.Set("amount", "1000.1")

    client := &http.Client{}
    req, _ := http.NewRequest("POST", ts.URL + "/deposit", strings.NewReader(v.Encode()))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Add("X-Auth-Token", s.Token)
    req.Header.Add("X-Auth-User", user.Id)
    resp, _ := client.Do(req)

    if resp.StatusCode != http.StatusNotFound {
        t.Errorf("Should not find account")
    }

    accounts = user.Accounts()
    if accounts[0].Ballance != 0 {
        t.Errorf("Should not update ballance")
    }
}


func TestTransactionDepositWithdraw(t *testing.T) {
    clearDB()
    user, pass := createUser()

    ts := httptest.NewServer(http.DefaultServeMux)
    defer ts.Close()

    s, _ := NewUserSession(user.Id, pass, "127.0.0.1")

    accounts := user.Accounts()

    // Deposit
    v := url.Values{}
    v.Set("to", accounts[0].Id)
    v.Set("amount", "1000.1")

    client := &http.Client{}
    req, _ := http.NewRequest("POST", ts.URL + "/deposit", strings.NewReader(v.Encode()))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Add("X-Auth-Token", s.Token)
    req.Header.Add("X-Auth-User", user.Id)
    resp, _ := client.Do(req)

    // Withdraw
    v = url.Values{}
    v.Set("from", accounts[0].Id)
    v.Set("amount", "500.1")

    client = &http.Client{}
    req, _ = http.NewRequest("POST", ts.URL + "/withdraw", strings.NewReader(v.Encode()))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Add("X-Auth-Token", s.Token)
    req.Header.Add("X-Auth-User", user.Id)
    resp, _ = client.Do(req)

    if resp.StatusCode != http.StatusOK {
        t.Errorf("Should process transaction %d", resp.StatusCode)
    }

    accounts = user.Accounts()
    if accounts[0].Ballance != 500 {
        t.Errorf("Should update ballance")
    }
}

func TestTransactionDepositWithdrawNoMoney(t *testing.T) {
    clearDB()
    user, pass := createUser()

    ts := httptest.NewServer(http.DefaultServeMux)
    defer ts.Close()

    s, _ := NewUserSession(user.Id, pass, "127.0.0.1")

    accounts := user.Accounts()

    // Deposit
    v := url.Values{}
    v.Set("to", accounts[0].Id)
    v.Set("amount", "1000.1")

    client := &http.Client{}
    req, _ := http.NewRequest("POST", ts.URL + "/deposit", strings.NewReader(v.Encode()))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Add("X-Auth-Token", s.Token)
    req.Header.Add("X-Auth-User", user.Id)
    resp, _ := client.Do(req)

    // Withdraw
    v = url.Values{}
    v.Set("from", accounts[0].Id)
    v.Set("amount", "1500.1") // More then we have

    client = &http.Client{}
    req, _ = http.NewRequest("POST", ts.URL + "/withdraw", strings.NewReader(v.Encode()))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Add("X-Auth-Token", s.Token)
    req.Header.Add("X-Auth-User", user.Id)
    resp, _ = client.Do(req)

    if resp.StatusCode != http.StatusBadRequest {
        t.Errorf("Should not process transaction %d", resp.StatusCode)
    }

    accounts = user.Accounts()
    if accounts[0].Ballance != 1000.1 {
        t.Errorf("Should not update ballance")
    }
}

func TestTransactionTransfer(t *testing.T) {
    clearDB()
    user, pass := createUser()
    user2, _ := createUser2()

    ts := httptest.NewServer(http.DefaultServeMux)
    defer ts.Close()

    s, _ := NewUserSession(user.Id, pass, "127.0.0.1")

    accounts := user.Accounts()
    accounts2 := user2.Accounts()

    // Deposit money to first account
    v := url.Values{}
    v.Set("to", accounts[0].Id)
    v.Set("amount", "1000.1")

    client := &http.Client{}
    req, _ := http.NewRequest("POST", ts.URL + "/deposit", strings.NewReader(v.Encode()))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Add("X-Auth-Token", s.Token)
    req.Header.Add("X-Auth-User", user.Id)
    resp, _ := client.Do(req)

    // Transfer money to second account
    v = url.Values{}
    v.Set("from", accounts[0].Id)
    v.Set("to", accounts2[0].Id)
    v.Set("amount", "500.1")

    client = &http.Client{}
    req, _ = http.NewRequest("POST", ts.URL + "/transfer", strings.NewReader(v.Encode()))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Add("X-Auth-Token", s.Token)
    req.Header.Add("X-Auth-User", user.Id)
    resp, _ = client.Do(req)

    if resp.StatusCode != http.StatusOK {
        t.Errorf("Should not process transaction %d", resp.StatusCode)
    }

    accounts = user.Accounts()
    if accounts[0].Ballance != 500 {
        t.Errorf("Should update ballance of first user")
    }

    accounts2 = user2.Accounts()
    if accounts2[0].Ballance != 500.1 {
        t.Errorf("Should update ballance of second user")
    }
}