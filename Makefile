build:
	go build -o ppaid cmd/ppaid/*.go

deploy: build
	mv ppaid $(HOME)/go/bin/

credit:
	gocredits . > CREDITS
# gocredits . > CREDITS
# could not find the license for "github.com/mattn/go-localereader"
# make: *** [Makefile:8: credit] Error 1