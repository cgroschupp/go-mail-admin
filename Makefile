build-frontend:
	npm --prefix frontend ci
	npm --prefix frontend run build

build-backend:
	go build ./cmd/go-mail-admin/...

build: build-frontend build-backend

clean:
	rm -rf frontend/dist
	rm -rf dist

run:
	GOMAILADMIN_AUTH_PASSWORD=develop GOMAILADMIN_AUTH_SECRET=develop ./cmd/go-mail-admin/...

test:
	go test ./...

generate:
	tsp compile .
	go generate ./...
