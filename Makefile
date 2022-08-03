TEMPLATE = template.yaml
PACKAGED_TEMPLATE = packaged.yaml

.PHONY: build
build: clean fmt
	sam build

.PHONY: clean
clean:
	rm -rf .aws-sam
	rm -f $(PACKAGED_TEMPLATE)

.PHONY: fmt
fmt:
	make -C lambda fmt

.PHONY: test
test: fmt
	make -C lambda test

.PHONY: install
install:
	make -C lambda install

.PHONY: local
local: fmt build
	sam local invoke --env-vars env.json

.PHONY: api
api: fmt build
	sam local start-api --env-vars env.json

.PHONY: validate
validate:
	sam validate --config-file samconfig.toml \
		--template-file $(TEMPLATE)

.PHONY: package
package: build
	. .env
	sam package --s3-bucket $(S3_BUCKET) \
		--template-file $(TEMPLATE) \
		--output-template-file $(PACKAGED_TEMPLATE)

.PHONY: deploy
deploy: package
	. .env
	sam deploy --stack-name $(STACK_NAME) \
		--template-file $(PACKAGED_TEMPLATE) \
		--no-confirm-changeset
