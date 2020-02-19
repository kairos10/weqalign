package main

//go:generate ./genResources.sh

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
	"strings"
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

func HttpStaticFilesHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fName := r.URL.Path
		prefix := "/"
		if strings.HasPrefix(fName, prefix) {
			fName = fName[len(prefix):]
		}
		switch {
		case true:
			_ = GetSessionId(w, r)

			fContentStr, ok := _genFiles[fName]
			if ok {
				fContent, err := base64.StdEncoding.DecodeString(fContentStr)
				if err == nil {
					http.ServeContent(w, r, fName, time.Time{}, bytes.NewReader(fContent))
					break
				}
				fmt.Println("BASE64 ERR: ", err)
			}

			fallthrough
		default:
			fmt.Println("NOT found: ", fName)
			http.Error(w, "not found", http.StatusNotFound)
		}

	}
}
