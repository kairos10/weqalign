package main

import (
	"net/http"
	"time"
	"encoding/json"
	"strconv"
)

func PublishResourcesHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type jsonImgRes struct {
			Path string
			SolverResponse interface{}
		}
		jsonResponse := make(map[string]jsonImgRes)
			
		_ = GetSessionId(w, r)
		time.Sleep(100 * time.Millisecond) // slow down crazy clients

		lastResourceId := resources.getLastResourceId()

		// serve all resources with ID > client's last ID
		clientLastImgId, err := strconv.ParseUint(r.FormValue("lastImgId"), 10, 64)
		if err != nil { clientLastImgId=0 }

		// serve only the requested resource
		clientImgId, err := strconv.ParseUint(r.FormValue("imgid"), 10, 64)
		if err != nil { clientImgId=0 }

		var fromId, toId uint64
		if clientLastImgId > 0 {
			fromId, toId = clientLastImgId+1, lastResourceId
		} else if clientImgId > 0 {
			fromId, toId = clientImgId, clientImgId
		} else {
			fromId, toId = 1, lastResourceId
		}

		for id:=fromId; id<=toId; id++ {
			resource, ok := resources.getResource(id)
			if !ok { continue }
			var jsonRes jsonImgRes
			jsonRes.Path = resource.webPath
			jsonRes.SolverResponse = resource.solverResponse
			jsonResponse[resource.imgId] = jsonRes
		}

		if len(jsonResponse) > 0 {
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Cache-Control", "no-cache, must-revalidate")
			w.Header().Set("Expires", "0")
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
			json, _ := json.MarshalIndent(jsonResponse, "", "   ")
			w.Write(json)
		} else {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "", http.StatusServiceUnavailable)
		}
	}
}

