package main

//go:generate ./genResources.sh

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func GetSessionId(w http.ResponseWriter, r *http.Request) (sessId string) {
	cookie, err := r.Cookie(SESSION_COOKIE)
	if err != nil || len(cookie.Value) < 10 {
		sessId = "sess_" + fmt.Sprintf("%v", time.Now().Unix())
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: SESSION_COOKIE, Value: sessId, Path: "/", Expires: expiration, MaxAge: 86400, Secure: false, HttpOnly: false}
		http.SetCookie(w, &cookie)
	} else {
		sessId = cookie.Value
	}
	return
}

func HttpFileHandler(fName string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = GetSessionId(w, r)

		fContent, ok := _genFiles[fName]
		if ok {
			http.ServeContent(w, r, fName, time.Time{}, bytes.NewReader([]byte(fContent)))
		} else {
			fmt.Println("NOT found: ", fName)
			http.Error(w, "not found", http.StatusNotFound)
		}
	}
}
