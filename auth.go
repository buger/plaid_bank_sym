package main

import (
    "net/http"
    "crypto/sha512"
    "fmt"
    "time"
    "sort"
    "encoding/json"
)

type UserSession struct {
    Token string
    UserId string
    Ip string
    Active bool
    CreatedAt time.Time
}

type UserNotAuthorizedErr struct {
    msg string
}

func (e *UserNotAuthorizedErr) Error() string {
    return "Auth token is wrong or expired"
}

type UserNotExistsErr struct {
    userId string
}

func (e *UserNotExistsErr) Error() string {
    return fmt.Sprintf("User with id %s not exists", e.userId)
}

type UserWrondPasswordErr struct {
    msg string
}

func (e *UserWrondPasswordErr) Error() string {
    return e.msg
}


func NewUserSession(userId, password, ip string) (*UserSession, error) {
    return newUserSessionWithDate(userId, password, ip, time.Now())
}

func newUserSessionWithDate(userId, password, ip string, createdAt time.Time) (*UserSession, error) {
    var user User
    if err := db.ReadJSON(fmt.Sprintf("users/%s.json", userId), &user); err != nil {
        return nil, &UserNotExistsErr{userId}
    }

    if user.PasswordHash != sha512.Sum512([]byte(password)) {
        return nil, &UserWrondPasswordErr{"Wrong password"}
    }

    s := &UserSession{
        Token: randStringBytes(50),
        UserId: user.Id,
        Ip: ip,
        Active: true,
        CreatedAt: createdAt,
    }

    if err := db.WriteJSON(fmt.Sprintf("users/%s/sessions/%s", userId, s.Token), s); err != nil {
        return nil, err
    }

    return s, nil
}

func (s *UserSession) IsValid() bool {
    if !s.Active || time.Now().Sub(s.CreatedAt) > 30 * time.Minute {
        return false
    }

    return true
}


type ByCreatedAt []UserSession

func (a ByCreatedAt) Len() int           { return len(a) }
func (a ByCreatedAt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCreatedAt) Less(i, j int) bool { return a[j].CreatedAt.Before(a[i].CreatedAt)}

// Checks if session with provide token still active or expired
func IsSessionActive(userId string, token string) bool {
    sessionTokens, _ := db.AllKeys(fmt.Sprintf("users/%s/sessions", userId))

    if len(sessionTokens) == 0 {
        return false
    }

    var sessions []UserSession
    for _, tok := range sessionTokens {
        var s UserSession
        db.ReadJSON(fmt.Sprintf("users/%s/sessions/%s", userId, tok), &s)
        sessions = append(sessions, s)
    }

    sort.Sort(ByCreatedAt(sessions))

    if sessions[0].Token != token {
        return false
    }

    if !sessions[0].IsValid() {
        return false
    }

    return true
}

type UserAuthResp struct {
    Token string `json:"token,omitempty"`
    Error string `json:"error,omitempty"`
}

func userAuthHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")

    userId := r.FormValue("user")
    password := r.FormValue("password")

    s, err := NewUserSession(userId, password, r.RemoteAddr)

    if err != nil {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(&UserAuthResp{Error: err.Error()})
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(&UserAuthResp{Token: s.Token})
}

type AuthErrorResp struct {
    Error string `json:"error"`
}

func requireAuthHandler(cb func(*User, http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        token := r.Header.Get("X-Auth-Token")
        userId := r.Header.Get("X-Auth-User")

        if !IsSessionActive(userId, token) {
            w.WriteHeader(http.StatusUnauthorized)
            err := &UserNotAuthorizedErr{}
            json.NewEncoder(w).Encode(&ErrorResp{Error: err.Error()})
            return
        }

        var user User
        db.ReadJSON(fmt.Sprintf("users/%s.json", userId), &user)

        cb(&user, w, r)
    }
}

func init() {
    http.HandleFunc("/auth", userAuthHandler)
}