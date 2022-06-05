TEMPLATE = template.yaml
PACKAGED_TEMPLATE = packaged.yaml

.PHONY: build
build: clean
	sam build

.PHONY: clean
clean:
	rm -f $(PACKAGED_TEMPLATE)

.PHONY: test
test:
	go test -v ./go-lambda-twitter

.PHONY: install
install:
	go get ./go-lambda-twitter

.PHONY: local
local: build
	sam local invoke --env-vars env.json

.PHONY: api
api: build
	sam local start-api --env-vars env.json

.PHONY: validate
validate:
	sam validate --config-file samconfig.toml --template-file template.yaml

.PHONY: package
package: build
	. .env
	sam package --s3-bucket $(S3_BUCKET) \
		--template-file $(TEMPLATE) \
		--output-template-file packaged.yaml

.PHONY: deploy
deploy: package
	. .env
	sam deploy --stack-name $(STACK_NAME) \
		--template-file $(PACKAGED_TEMPLATE) \
		--no-confirm-changeset
