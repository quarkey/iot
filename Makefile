build:
	go build ./cmd/api
	go build ./cmd/init

testdata:
	./init --conf exampleconfig.json --testdata

clean:
	rm -f api
	rm -f init


