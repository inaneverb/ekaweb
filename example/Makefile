GO_BUILD_TAGS:=

full:
	go mod tidy
	go run -tags=${GO_BUILD_TAGS} server_full/full.go

minimal:
	go mod tidy
	go run -tags=${GO_BUILD_TAGS} server_minimal/minimal.go

query:
	go mod tidy
	go run -tags=${GO_BUILD_TAGS} server_query/query.go

jrpc:
	go mod tidy
	go run -tags=${GO_BUILD_TAGS} server_jrpc/jrpc.go