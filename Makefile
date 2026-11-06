.PHONY: proto

proto:
	go tool buf format -w
	go tool buf lint
	go tool buf generate
	cd gen/go && go mod tidy
