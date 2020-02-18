.PHONY: all clean

.DEFAULT_GOAL := all

GOSOURCES = \
	httphandle_solve.go \
	httphandle_static.go \
	main.go \
	temp_web.go \
	resources.go

all: weqalign

temp_web.go: genResources.sh web/index.html web/main.js
	go generate

clean:
	rm -f temp_web.go weqalign

distclean: clean
	rm -rf --one-file-system builder/BUILD/

weqalign: $(GOSOURCES)
	go build
	strip main
	mv main $@

