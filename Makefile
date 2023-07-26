build:
	go build ./cmd/api
	go build ./cmd/init
	go build ./cmd/simulator

sim: 
	go build ./cmd/simulator
	./simulator --conf ./config/exampleconfig.json
	
testdata:
	go build ./cmd/init
	./init --conf ./config/exampleconfig.json --automigrate

clean:
	rm -f api
	rm -f init
	rm -f simulator

downup:
	migrate -database ${POSTGRESQL_URL} -path database/migrations down
	migrate -database ${POSTGRESQL_URL} -path database/migrations up

down:
	migrate -database ${POSTGRESQL_URL} -path database/migrations down

up:
	migrate -database ${POSTGRESQL_URL} -path database/migrations up

run:
	go build ./cmd/api
	./api --conf config/exampleconfig.json --automigrate

race:
	go run -race ./cmd/api --conf config/exampleconfig.json --automigrate
rundebug:
	go build ./cmd/api
	./api --conf config/exampleconfig.json --debug --automigrate

deploy:
	GOOS=linux GOARCH=arm go build ./cmd/api
	scp -r api config/rpi_prod.json resources database client/dist slundin@192.168.10.128:/iot

golint:
	golangci-lint run --out-format code-climate | tee gl-code-quality-report.json | jq -r '.[] | "\(.location.path):\(.location.lines.begin) \(.description)"'