# go-lambda-twitter

[![ci](https://github.com/ega4432/go-lambda-twitter/actions/workflows/ci.yaml/badge.svg)](https://github.com/ega4432/go-lambda-twitter/actions/workflows/ci.yaml)

This is a repository for **go-lambda-twitter** - Twitter client built by API Gateway + Lambda in AWS.

```shell
.
├── Makefile                    <-- Make to automate build
├── README.md                   <-- This instructions file
├── lambda                      <-- Source code for a lambda function
│   ├── main.go                 <-- Lambda function code
│   └── main_test.go            <-- Unit tests
└── template.yaml
```

## Requirements

* AWS CLI already configured with Administrator permission
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

## Setup process

### Installing dependencies & building the target

```shell
make
```

### Local development

**Invoking function locally through local API Gateway**

```shell
make api
```

**Calling health check endpoint with cURL**

```shell
# GET
curl --request GET \
    --url http://127.0.0.1:3000/health | jq .
{
  "message": "OK"
}

# POST
curl --request POST \
  --url http://127.0.0.1:3000/tweet \
  -d '{ "text": "test" }' | jq .
{
  "message": "Tweeted successfully",
  "tweet_text": "test",
  "tweet_url":"https://twitter.com/ega4432/status/1555603271665913856"
}
```

## Packaging and deployment

From local

```shell
make deploy
```
