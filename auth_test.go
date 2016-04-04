package main

import (
    "testing"
    "time"
    "encoding/json"
    "net/url"
    "net/http/httptest"
    "io/ioutil"
    "net/http"
)

func createUser() (*User, string) {
    name := "John Dou"
    birthDate := time.Date(1987, time.May, 26, 0, 0, 0, 0, time.UTC)
    user, pass, _ := CurrentBank.CreateUser(name, birthDate)
    return user, pass
}

func createUser2() (*User, string) {
    name := "Billy Joe"
    birthDate := time.Date(1987, time.May, 26, 0, 0, 0, 0, time.UTC)
    user, pass, _ := CurrentBank.CreateUser(name, birthDate)
    return user, pass
}


func TestAuth(t *testing.T) {
    clearDB()

    _, err := NewUserSession("not_existing", "","127.0.0.1")
    if err == nil {
        t.Errorf("Should not create session for not existing user")
    } else {
        if _, ok := err.(*UserNotExistsErr); !ok {
            t.Errorf("Should throw user not exist error")
        }
    }

    user, pass := createUser()

    _, err = NewUserSession(user.Id, "wrongPass", "127.0.0.1")
    if err == nil {
        t.Errorf("Should not create session because of wrong password")
    } else {
        if _, ok := err.(*UserWrondPasswordErr); !ok {
            t.Errorf("Should throw user wrong password error %v", err)
        }
    }

    s1, err := NewUserSession(user.Id, pass, "127.0.0.1")
    if err != nil {
        t.Errorf("Should create session %v", err)
    }

    s2, err := NewUserSession(user.Id, pass, "127.0.0.1")
    if err != nil {
        t.Errorf("Should create second session %v", err)
    }

    if IsSessionActive(user.Id, s1.Token) {
        t.Errorf("First session should expire %v", s1)
    }

    if !IsSessionActive(user.Id, s2.Token) {
        t.Errorf("Second sessions should be active %v", s2)
    }

    clearDB()
    user, pass = createUser()

    s3, err := newUserSessionWithDate(user.Id, pass, "127.0.0.1", time.Now().Add(-31*time.Minute))
    if err != nil {
        t.Errorf("Should create expired session %v", err)
    }

    if IsSessionActive(user.Id, s3.Token) {
        t.Errorf("Third session should be expired %v", s3)
    }
}


func TestAuthAPI(t *testing.T) {
    clearDB()

    user, pass := createUser()

    ts := httptest.NewServer(http.DefaultServeMux)
    defer ts.Close()

    v := url.Values{}
    v.Set("user", user.Id)
    v.Set("password", pass)

    resp, _ := http.PostForm(ts.URL + "/auth", v)

    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status %d, was: %d", http.StatusOK, resp.StatusCode)
    }

    body, _ := ioutil.ReadAll(resp.Body)
    resp.Body.Close()

    var data map[string]string
    json.Unmarshal(body, &data)

    if data["token"] == "" || len(data["token"]) != 50 {
        t.Errorf("Should auth user %v", data)
    }
}

func TestAuthAPIWrongPassword(t *testing.T) {
    clearDB()

    user, _ := createUser()

    ts := httptest.NewServer(http.DefaultServeMux)
    defer ts.Close()

    v := url.Values{}
    v.Set("user", user.Id)
    v.Set("password", "wrong")

    resp, _ := http.PostForm(ts.URL + "/auth", v)

    if resp.StatusCode != http.StatusUnauthorized {
        t.Errorf("Expected status %d, was: %d", http.StatusOK, resp.StatusCode)
    }

    body, _ := ioutil.ReadAll(resp.Body)
    resp.Body.Close()

    var data map[string]string
    json.Unmarshal(body, &data)

    if data["error"] == "" {
        t.Errorf("Should show error %v", data)
    }
}
