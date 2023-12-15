setup:
	go mod download

gen-proto:
	buf generate proto