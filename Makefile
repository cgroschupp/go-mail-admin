client-build: clean
	cd frontend && npm install && npx vue-cli-service build

clean:
	rm -rf frontend/dist

interface-build: client-build
	go build

run:
	GOMAILADMIN_PASSWORD_SCHEME=BLF-CRYPT GOMAILADMIN_DB="root:develop@tcp(127.0.0.1:3306)/vmail" GOMAILADMIN_V3="on" GOMAILADMIN_AUTH_Username="test" GOMAILADMIN_AUTH_Password="test"  go run ./...

init-test:
	docker-compose down
	docker-compose rm
	docker-compose up -d
	sleep 10

test:
	GOMAILADMIN_DB="root:develop@tcp(127.0.0.1:3306)/vmail" GOMAILADMIN_V3="off" go test ./...

build: interface-build
