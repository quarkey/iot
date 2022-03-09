build-all:
	go build
	go build ./cmd/init

initdb:
	./init --conf exampleconfig.json --testdata --drop

