include .env
export

.PHONY: test
test:
	go test ./auth0/... -tags=integration -cover

.PHONY: test-auth0
test-auth0:
	go test ./auth0 -tags=integration -cover