build:
	go build ./cmd/api
	go build ./cmd/init

testdata:
	go build ./cmd/init
	./init --conf ./config/exampleconfig.json --automigrate

clean:
	rm -f api
	rm -f init


