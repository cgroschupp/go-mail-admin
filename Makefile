client-build: clean
	npm --prefix frontend run build

clean:
	rm -rf frontend/dist

interface-build: client-build
	go build ./cmd/go-mail-admin/...

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

swag:
	swag init --generalInfo cmd/go-mail-admin/main.go -o internal/docs

generate:
	tsp compile .
	go generate ./...
