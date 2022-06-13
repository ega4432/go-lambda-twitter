TEMPLATE = template.yaml
PACKAGED_TEMPLATE = packaged.yaml

.PHONY: build
build: clean
	make -C lambda build
	sam build

.PHONY: clean
clean:
	rm -f $(PACKAGED_TEMPLATE)

.PHONY: test
test:
	make -C lambda test

.PHONY: install
install:
	make -C lambda install

.PHONY: local
local: build
	sam local invoke --env-vars env.json

.PHONY: api
api: build
	sam local start-api --env-vars env.json

.PHONY: validate
validate:
	sam validate --config-file samconfig.toml --template-file $(TEMPLATE)

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
