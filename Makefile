build:
	go build -o ppaid cmd/ppaid/*.go

deploy: build
	mv ppaid $(HOME)/go/bin/