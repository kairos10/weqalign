package main

import (
	"fmt"
	"golang.org/x/net/webdav"
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
}

func main() {
	//////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////
	// test1
	go func() {
		<-time.After(time.Second)
		_ = resources.addResource("api/am/testres/img1.jpg")
		_ = resources.addResource("api/am/testres/img2.jpg")
		buf := make([]byte, 100)
		for i := 0; i < 4; i++ {
			fmt.Println("req", i)
			k := fmt.Sprintf("%d", 1+i%2)
			res, _ := http.Get("http://localhost:8080/SolveField?imgid=" + k)
			res.Body.Read(buf)
			fmt.Println(string(buf))
		}
		res, _ := http.Get("http://localhost:8080/index.html")
		res.Body.Read(buf)
		fmt.Println("index.html: ", string(buf))
		res, _ = http.Get("http://localhost:8080/main.js")
		res.Body.Read(buf)
		fmt.Println("main.js: ", string(buf))

		fmt.Println("TEST DONE!")
		//os.Exit(0)
	}()
	//////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////

	cmdParams.tcpPort = "8080"
	cmdParams.imgFolder = "/tmp/_images"

	// setup [/img]
	os.MkdirAll(cmdParams.imgFolder, os.ModeDir|os.ModePerm|0755)
	d, err := os.Stat(cmdParams.imgFolder)
	if err != nil || !d.IsDir() {
		log.Fatalf("Something is wrong with [%v] [%v]\n", cmdParams.imgFolder, err)
	}
	fs := http.FileServer(http.Dir(cmdParams.imgFolder))
	http.Handle("/img/", http.StripPrefix("/img", fs))

	// setup webDav
	wDavH := &webdav.Handler{
		Prefix:     "/dav/",
		FileSystem: webdav.Dir(cmdParams.imgFolder),
		LockSystem: webdav.NewMemLS(),
		//Logger: func(r *http.Request, err error) { fmt.Println("DAV: " + r.Method + " " + r.URL.String()) },
	}
	http.Handle("/dav/", wDavH)
	hostname, _ := os.Hostname()
	fmt.Println("webdav server set up at: \\\\" + hostname + "@" + cmdParams.tcpPort + "\\dav" + "\\DavWWWRoot")
	fmt.Println("http server set up at: http://" + hostname + ":" + cmdParams.tcpPort + "/")

	// setup static files
	http.HandleFunc("/", HttpFileHandler("web/index.html"))
	http.HandleFunc("/main.js", HttpFileHandler("web/main.js"))

	// setup plate solver
	http.HandleFunc("/SolveField", GetHttpPlateSolver())

	fmt.Println("Listening...")
	log.Fatal(http.ListenAndServe(":"+cmdParams.tcpPort, nil))
}
