client-build:
	rm -f -r  ./mailserver-configurator-client/dist/
	cd ./mailserver-configurator-client; npm install
	cd ./mailserver-configurator-client; npm run build

interface-copy-client:
	rm -f -r ./internal/public/*
	cp -r ./mailserver-configurator-client/dist/* ./internal/public/

interface-install-deps:
	go get github.com/rakyll/statik

interface-build:
	cd ./internal; ~/go/bin/statik -f -src=./public
	cd ./internal; go build -o ../go-mail-admin-with-gui ./

statik:
	cd ./internal; ~/go/bin/statik -f -src=./public

run:
	GOMAILADMIN_PASSWORD_SCHEME=BLF-CRYPT GOMAILADMIN_DB="root:develop@tcp(127.0.0.1:3306)/vmail" GOMAILADMIN_V3="on" GOMAILADMIN_AUTH_Username="test" GOMAILADMIN_AUTH_Password="test"  go run ./internal

gorelease-vue:
	go get github.com/rakyll/statik
	rm -f -r  ./mailserver-configurator-client/dist/
	cd ./mailserver-configurator-client; npm install
	cd ./mailserver-configurator-client; npm run build
	mkdir ./mailserver-configurator-interface/public/
	cp -r ./mailserver-configurator-client/dist/* ./mailserver-configurator-interface/public/
	cd ./mailserver-configurator-interface; ~/go/bin/statik -f -src=./public

init-test:
	docker-compose down
	docker-compose rm
	docker-compose up -d
	sleep 10

test:
	GOMAILADMIN_DB="root:develop@tcp(127.0.0.1:3306)/vmail" GOMAILADMIN_V3="off" go test ./internal


build: client-build interface-copy-client interface-install-deps interface-build
