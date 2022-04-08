build:
	go build ./cmd/api
	go build ./cmd/init
	go build ./cmd/simulator

testdata:
	go build ./cmd/init
	./init --conf ./config/exampleconfig.json --automigrate

clean:
	rm -f api
	rm -f init

downup:
	migrate -database ${POSTGRESQL_URL} -path database/migrations down
	migrate -database ${POSTGRESQL_URL} -path database/migrations up

run:
	go build ./cmd/api
	./api --conf config/exampleconfig.json
