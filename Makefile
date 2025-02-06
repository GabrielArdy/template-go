.PHONY: generate build init

API_SPEC := apis/api.yml

generate:
	@echo "Generating API client code..."
	@mkdir -p generated
	oapi-codegen \
		--package generated \
		-generate types,server,spec \
		$(API_SPEC) > generated/api.gen.go
init:
	go mod tidy

build:
	go build 
