.PHONY: all clean

all: weqalign

temp_web.go: web/index.html web/main.js
	go generate

clean:
	rm -f temp_web.go

weqalign: httphandle_solve.go  httphandle_static.go  main.go  temp_web.go
	go build
	strip main
	mv main $@

