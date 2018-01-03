include .env
export

.PHONY: test
test:
	go test ./auth0/... -tags=integration -cover
