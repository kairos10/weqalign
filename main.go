package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	SESSION_COOKIE                = "t2_session"
	DEFAULT_ALIGN_POINTS_DISTANCE = 2
)

var cmdParams struct {
	tcpPort   string
	imgFolder string
	imgFolderWeb string
}

func main() {
	//////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////
	// test1
	go func() {
		<-time.After(time.Second)
		_ = resources.addResource("api/am/testres/img1.jpg", "")
		_ = resources.addResource("api/am/testres/img2.jpg", "")
		buf := make([]byte, 100)
		for i := 0; i < 4; i++ {
			fmt.Println("req", i)
			k := fmt.Sprintf("%d", 1+i%2)
			res, _ := http.Get("http://localhost:8080/SolveField?imgid=" + k)
			res.Body.Read(buf)
			fmt.Println(string(buf))
		}

		//res, _ := http.Get("http://localhost:8080/")
		//res.Body.Read(buf)
		//fmt.Println("index.html: ", string(buf))

		//res, _ = http.Get("http://localhost:8080/web/main.js")
		//res.Body.Read(buf)
		//fmt.Println("main.js: ", string(buf))

		res, _ := http.Get("http://localhost:8080/resources?lastImgId=1")
		res.Body.Read(buf)
		fmt.Println("resources?lastImgId=1: ", string(buf))

		fmt.Println("TEST DONE!")
		//os.Exit(0)
	}()
	//////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////

	cmdParams.tcpPort = "8080"
	cmdParams.imgFolder = "/tmp/_images"
	cmdParams.imgFolderWeb = "/img"
	
	err := setImageFolder(cmdParams.imgFolder, cmdParams.imgFolderWeb)
	if err != nil {
		log.Fatalf("ERROR setting image folder[%s]->[%s]: %v\n", cmdParams.imgFolder, cmdParams.imgFolderWeb, err)
	}

	// setup static files
	http.HandleFunc("/web/", HttpStaticFilesHandler())
	http.Handle("/", http.RedirectHandler("/web/index.html", http.StatusPermanentRedirect))

	// setup plate solver
	http.HandleFunc("/SolveField", GetHttpPlateSolver())

	// publish resources
	http.HandleFunc("/resources", PublishResourcesHandler())

	fmt.Println("Listening...")
	log.Fatal(http.ListenAndServe(":"+cmdParams.tcpPort, nil))
}

func setImageFolder(imgFolder, imgFolderWeb string) (err0 error) {
	// setup [/img]
	os.MkdirAll(imgFolder, os.ModeDir|os.ModePerm|0755)
	d, err := os.Stat(imgFolder)
	if err0=err; err0 != nil { return }
	if !d.IsDir() {
		err0 = fmt.Errorf("Something is wrong with [%v] [%v]\n", cmdParams.imgFolder, err)
		return
	}

	err0 = StartResGeneratorFromFolder(imgFolder, imgFolderWeb)
	if err0 != nil { return }

	fs := http.FileServer(http.Dir(imgFolder))
	http.Handle(imgFolderWeb + "/", http.StripPrefix(imgFolderWeb, fs))

	return
}
