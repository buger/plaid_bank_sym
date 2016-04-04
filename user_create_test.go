package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func clearDB() {
	OpenDB(".test").RemoveAll()
	Init(1, ".test")
}

func TestUserCreate(t *testing.T) {
	clearDB()

	name := "John Dou"
	birthDate := time.Date(1987, time.May, 26, 0, 0, 0, 0, time.UTC)

	user, pass, err := CurrentBank.CreateUser(name, birthDate)

	if err != nil {
		t.Errorf("Should create user %v", err)
	} else {
		if user.Name != name || pass == "" || len(pass) != 10 {
			t.Errorf("Should fill user record with data: %v, %s", user, pass)
		}
	}

	user, _, err = CurrentBank.CreateUser(name, birthDate)
	if err == nil {
		t.Errorf("Should not allow users with same credentials %v", err)
	} else {
		if _, ok := err.(*UserExistsErr); !ok {
			t.Errorf("Should raise user exist error")
		}
	}

	user, _, err = CurrentBank.CreateUser(name, birthDate.Add(time.Hour))
	if err != nil {
		t.Errorf("Should create user with different birth date: %v", err)
	}

	user, _, err = CurrentBank.CreateUser(name, time.Now())
	if err == nil {
		t.Errorf("Should not allow users yonger then 18 years %v", err)
	} else {
		if _, ok := err.(*UserAgeErr); !ok {
			t.Errorf("Should raise user age error")
		}
	}

	user, _, err = CurrentBank.CreateUser(name, time.Now().Add(-101*durYear))
	if err == nil {
		t.Errorf("Should not allow users older then 100 years")
	} else {
		if _, ok := err.(*UserAgeErr); !ok {
			t.Errorf("Should raise user age error %v", err)
		}
	}
}

func testUserCreateAPI(t *testing.T, name string, birthDate string, expectedStatus int) (data map[string]string) {
	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()

	v := url.Values{}
	v.Set("name", name)
	v.Set("dateOfBirth", birthDate)

	resp, _ := http.PostForm(ts.URL+"/user/create", v)

	if resp.StatusCode != expectedStatus {
		t.Errorf("Expected status %d, was: %d", expectedStatus, resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	json.Unmarshal(body, &data)

	return data
}

func TestUserCreateAPISuccess(t *testing.T) {
	clearDB()

	data := testUserCreateAPI(
		t,
		"John Dou",
		time.Date(1987, time.May, 26, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		http.StatusCreated,
	)

	if data["error"] != "" {
		t.Errorf("Error field should be nil %s", data["error"])
	}

	if data["id"] == "" || len(data["id"]) != 10 {
		t.Errorf("Error return id field %s", data["id"])
	}

	if data["password"] == "" || len(data["password"]) != 10 {
		t.Errorf("Error return password field %s", data["password"])
	}
}

func TestUserCreateAPIDuplicate(t *testing.T) {
	clearDB()

	testUserCreateAPI(
		t,
		"John Dou",
		time.Date(1987, time.May, 26, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		http.StatusCreated,
	)
	data := testUserCreateAPI(
		t,
		"John Dou",
		time.Date(1987, time.May, 26, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		http.StatusBadRequest,
	)

	if data["error"] == "" {
		t.Errorf("Error field should not be blank %v", data)
	}

	if data["id"] != "" || data["password"] != "" {
		t.Errorf("Only error field should be presented")
	}
}

func TestUserCreateAPIWrongDate(t *testing.T) {
	clearDB()

	data := testUserCreateAPI(
		t,
		"John Dou",
		"asdasd", // Not date at all
		http.StatusBadRequest,
	)

	if data["error"] == "" {
		t.Errorf("Error field should not be nil %v", data)
	}
}
