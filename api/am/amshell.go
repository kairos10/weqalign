package am

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var AmShellCommands map[string]string

var LoggerFunc func(string)

func logger(s string) {
	if LoggerFunc != nil {
		LoggerFunc(s)
	}
}

func init() {
	if AmShellCommands == nil {
		AmShellCommands = make(map[string]string)
	}
	AmShellCommands["solver"] = "solve-field"
	AmShellCommands["wcsinfo"] = "wcsinfo"
	AmShellCommands["wcs-rd2xy"] = "wcs-rd2xy"
	AmShellCommands["tablist"] = "tablist"

	for _, p := range AmShellCommands {
		_, err := exec.LookPath(p)
		if err != nil {
			fmt.Printf("Checking path: [%s] NOT FOUND!\n\n", p)
		}
	}
}

const (
	DEFAULT_SOLVER_MIN_FIELD = 1
	DEFAULT_SOLVER_MAX_FIELD = 80
	MAX_PARALLEL_SOLVERS     = 1 // number of solvers allowed to run in parallel
)

var AmMinField, AmMaxField int // plate solver limits [min° - max°]; if not initialized, the values are derived from the first solved image

type ImageResource struct {
	FilePath    string // full path to the image file
	isProcessed bool   // file is already processed
	isInProcess bool   // the file is being worked on
	isSolved    bool   // the file has a solution

	FieldW, FieldH float64 // field width, height
	RACenter, DECCenter float64 // image center [RA;DEC]
	PixScale       float64 // pixel scale
	NegParity       bool   // the image has negative parity

	Stars       []XYStarPos // position of the stars identifies within the image
}
func (r ImageResource) String() string {
	status := "na"
	if r.isSolved { status="solved" } else if r.isProcessed { status="failed" } else if r.isInProcess { status="running" }
	return fmt.Sprintf("%s [%s]", r.FilePath, status)
}

func GetPlateSolver() func(*ImageResource) bool {
	procPool := struct {
		sync.Mutex
		availableProcesses int
	}{availableProcesses: MAX_PARALLEL_SOLVERS}

	if AmMinField == 0 {
		AmMinField = DEFAULT_SOLVER_MIN_FIELD
	}
	if AmMaxField == 0 {
		AmMaxField = DEFAULT_SOLVER_MAX_FIELD
	}

	return func(r *ImageResource) bool {
		var ok bool
		fileBase := r.FilePath[0 : len(r.FilePath)-len(filepath.Ext(r.FilePath))]

		if r.isProcessed {
			return true
		} else if r.isInProcess {
			return false
		} else if _, err := os.Stat(fileBase + ".working"); err == nil {
			// maybe some external solver is working on the file
			return false
		} else if _, err := os.Stat(fileBase + ".solved"); err == nil {
			// file is already solved
			r.isProcessed = true
			r.isSolved = true
			logger("SOLVED: " + fileBase)
			return true
		} else {
			r.isInProcess = true
		}

		procPool.Lock()
		if procPool.availableProcesses > 0 {
			procPool.availableProcesses--
			ok = true
		} else {
			ok = false
		}
		procPool.Unlock()

		if ok {
			logger(fmt.Sprintf("Solving [%v]\n", fileBase))
			cmd := exec.Command(AmShellCommands["solver"], r.FilePath, "-L", fmt.Sprintf("%v", AmMinField), "-H", fmt.Sprintf("%v", AmMaxField))
			logger("SOLVER+")
			startTime := time.Now()
			err := cmd.Run()
			if err != nil {
				logger(fmt.Sprintln("CMD Error: ", err))
			}
			logger(fmt.Sprintln("SOLVER-[", time.Now().Sub(startTime), "]"))

			procPool.Lock()
			procPool.availableProcesses++
			procPool.Unlock()

			r.isProcessed = true

			if _, err := os.Stat(fileBase + ".solved"); err == nil {
				r.isSolved = true
				PlateGetWcsInfo(r)
				min := int(math.Min(r.FieldW, r.FieldH)) - 1
				max := int(math.Max(r.FieldW, r.FieldH)) + 1
				if AmMinField == DEFAULT_SOLVER_MIN_FIELD && AmMaxField == DEFAULT_SOLVER_MAX_FIELD && min >= AmMinField && max <= AmMaxField {
					AmMinField, AmMaxField = min, max
					logger(fmt.Sprintln("Changing min/max solver field: [", AmMinField, ", ", AmMaxField, "]"))
				}
				logger("SOLVED!: " + fileBase)
			}
		}

		r.isInProcess = false
		return ok
	}
}

func PlateGetWcsInfo(res *ImageResource) {
	if res.PixScale != 0 || !res.isSolved {
		return
	}
	fileBase := res.FilePath[0 : len(res.FilePath)-len(filepath.Ext(res.FilePath))]
	cmd := exec.Command(AmShellCommands["wcsinfo"], fileBase+".wcs")
	stdout, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(stdout)
	err := cmd.Start()
	if err != nil {
		fmt.Println("WCS ERROR: ", err)
	}
	re := regexp.MustCompile("[^ \t]+")
	for scanner.Scan() {
		b := scanner.Bytes()
		m := re.FindAll(b, -1)
		if len(m) >= 2 {
			switch string(m[0]) {
			case "pixscale":
				pixScale, err := strconv.ParseFloat(string(m[1]), 64)
				if err == nil {
					res.PixScale = pixScale
				}
				//fmt.Println("PIXSCALE: ", pixScale)
			case "fieldw":
				fieldW, err := strconv.ParseFloat(string(m[1]), 64)
				if err == nil {
					res.FieldW = fieldW
				}
				//fmt.Println("fieldw: ", int(fieldW))
			case "fieldh":
				fieldH, err := strconv.ParseFloat(string(m[1]), 64)
				if err == nil {
					res.FieldH = fieldH
				}
				//fmt.Println("fieldh: ", int(fieldH))
			case "ra_center":
				raCenter, err := strconv.ParseFloat(string(m[1]), 64)
				if err == nil {
					res.RACenter = raCenter
				}
			case "dec_center":
				decCenter, err := strconv.ParseFloat(string(m[1]), 64)
				if err == nil {
					res.DECCenter = decCenter
				}
			case "parity":
				parity, err := strconv.ParseFloat(string(m[1]), 64)
				if err == nil {
					res.NegParity = (parity == 1)
				}
			}
		}
	}
	stdout.Close()
	cmd.Wait()
}


// calculate virtual X,Y for the given set of celestial coordinates (xy can be outside of the actual image)
// 	wcs-rd2xy -w img1.wcs -r 0 -d 90
//
//	RA,Dec (0.0000000000, 90.0000000000) -> pixel (141.8137606751, 191.8948501378)
func PlateRD2XY(res *ImageResource, ra int, dec int) (x float64, y float64) {
	fileBase := res.FilePath[0 : len(res.FilePath)-len(filepath.Ext(res.FilePath))]
	r := fmt.Sprintf("%v", ra)
	d := fmt.Sprintf("%v", dec)
	cmd := exec.Command(AmShellCommands["wcs-rd2xy"], "-w", fileBase+".wcs", "-r", r, "-d", d)
	stdout, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(stdout)
	cmd.Start()
	re := regexp.MustCompile("[+-]?[0-9.]+")
	for lineOk := false; scanner.Scan() && !lineOk; {
		b := scanner.Bytes()

		m := re.FindAll(b, -1)
		//fmt.Printf("[%v] %v %v %v %v\n", string(b), string(m[0]), string(m[1]), string(m[2]), string(m[3]))
		if len(m) >= 4 {
			lineOk = true
			x, _ = strconv.ParseFloat(string(m[2]), 64)
			y, _ = strconv.ParseFloat(string(m[3]), 64)
		}
	}
	stdout.Close()
	cmd.Wait()

	return
}


// star position in image
type XYStarPos struct {
	X          int
	Y          int
	Flux       float64
	Background float64
}

const MAX_NUM_STARS = 15
// tablist testres/img1.axy
//
//          X              Y           FLUX     BACKGROUND
// 1        147.220        124.876        249.000        3.00000
// 2        878.206        480.664        101.148        3.85217
func PlateGetStars(res *ImageResource) []XYStarPos {
	if res.Stars != nil {
		return res.Stars
	}
	stars := make([]XYStarPos, MAX_NUM_STARS) // return up to MAX_NUM_STARS star positions
	numStars := 0

	fileBase := res.FilePath[0 : len(res.FilePath)-len(filepath.Ext(res.FilePath))]
	cmd := exec.Command(AmShellCommands["tablist"], fileBase+".axy")
	stdout, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(stdout)
	cmd.Start()
	re := regexp.MustCompile("[+-]?[0-9.]+")
	for scanner.Scan() && numStars < len(stars) {
		b := scanner.Bytes()
		m := re.FindAll(b, -1)
		if len(m) < 5 {
			continue
		}

		x, _ := strconv.ParseFloat(string(m[1]), 64)
		y, _ := strconv.ParseFloat(string(m[2]), 64)
		stars[numStars].X = int(x)
		stars[numStars].Y = int(y)
		stars[numStars].Flux, _ = strconv.ParseFloat(string(m[3]), 64)
		stars[numStars].Background, _ = strconv.ParseFloat(string(m[4]), 64)

		numStars++
	}
	stdout.Close()
	cmd.Wait()
	if numStars > 0 {
		res.Stars = stars[0:numStars]
	}
	return res.Stars
}
