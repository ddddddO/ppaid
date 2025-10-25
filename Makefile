build:
	go build -o puco cmd/puco/*.go

deploy: build
	mv puco $(HOME)/go/bin/

credit:
	gocredits . > CREDITS
# gocredits . > CREDITS
# could not find the license for "github.com/mattn/go-localereader"
# make: *** [Makefile:8: credit] Error 1