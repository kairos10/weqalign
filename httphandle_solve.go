package main

import (
	"encoding/json"
	"fmt"
	"github.com/kairos10/weqalign/api/am"
	"math"
	"net/http"
)

const alignPointsDistance = 2 // distance in degrees for the DEC displacement

func GetHttpPlateSolver() func(http.ResponseWriter, *http.Request) {
	type solverStatus struct {
		FileId       string
		SolverStatus string
		//NCP, ...
		RA0DEC90_X         float64
		RA0DEC90_Y         float64
		RA0DEC95_X         float64
		RA0DEC95_Y         float64
		RA90DEC95_X        float64
		RA90DEC95_Y        float64
		PixScale           float64
		NegParity          bool
		NorthernHemisphere bool
		Stars              []am.XYStarPos
	}
	plateSolver := am.GetPlateSolver()
	return func(w http.ResponseWriter, r *http.Request) {
		imgId := r.FormValue("imgid")
		imgRes, ok := mapIdResource[imgId]
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		if imgRes.solverResponse != nil {
			ss := imgRes.solverResponse.(*solverStatus)
			json, _ := json.MarshalIndent(*ss, "", "   ")
			w.Write(json)
		} else {
			w.Header().Set("Retry-After", "10")
			http.Error(w, "", http.StatusServiceUnavailable)
			fmt.Println("not ready: " + imgId)
		}

		if imgRes.solverResponse == nil {
			plate := am.ImagePlate{FilePath: imgRes.path}
			_ = plateSolver(&plate)
			stat := plate.GetStatus()
			if stat == am.ImagePlateStatusSOLVED || stat == am.ImagePlateStatusFAILED {
				ss := solverStatus{FileId: imgId, SolverStatus: string(stat)}
				imgRes.solverResponse = &ss
				if stat == am.ImagePlateStatusSOLVED {
					am.PlateGetWcsInfo(&plate)

					ra0dec90x, ra0dec90y := am.PlateRD2XY(&plate, 0, 90)    // XY for NCP
					ra0dec270x, ra0dec270y := am.PlateRD2XY(&plate, 0, 270) // XY for SCP

					if ra0dec270x != 0 && ra0dec270y != 0 &&
						(ra0dec90x == 0 || math.Abs(ra0dec270x) < math.Abs(ra0dec90x)) &&
						(ra0dec90y == 0 || math.Abs(ra0dec270y) < math.Abs(ra0dec90y)) {
						// SCP closer
						ss.RA0DEC90_X, ss.RA0DEC90_Y = ra0dec270x, ra0dec270y
						ss.RA0DEC95_X, ss.RA0DEC95_Y = am.PlateRD2XY(&plate, 0, 270+alignPointsDistance)
						ss.RA90DEC95_X, ss.RA90DEC95_Y = am.PlateRD2XY(&plate, 90, 270+alignPointsDistance)
						ss.NorthernHemisphere = true
					} else {
						// use NCP
						ss.RA0DEC90_X, ss.RA0DEC90_Y = ra0dec90x, ra0dec90y
						ss.RA0DEC95_X, ss.RA0DEC95_Y = am.PlateRD2XY(&plate, 0, 90+alignPointsDistance)
						ss.RA90DEC95_X, ss.RA90DEC95_Y = am.PlateRD2XY(&plate, 90, 90+alignPointsDistance)
						ss.NorthernHemisphere = false
					}
					ss.PixScale = plate.PixScale
					ss.NegParity = plate.NegParity

					ss.Stars = am.PlateGetStars(&plate)
				}
			}
		}

	}
}
