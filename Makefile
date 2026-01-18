SERVER_VERSION ?= dev
RECOMMENDED_RUNNER_VERSION ?= latest
RECOMMENDED_CLI_VERSION ?= latest

LDFLAGS = -X gitlab.com/codebox4073715/codebox/config.ServerVersion=$(SERVER_VERSION) \
          -X gitlab.com/codebox4073715/codebox/config.RecommendedRunnerVersion=$(RECOMMENDED_RUNNER_VERSION) \
          -X gitlab.com/codebox4073715/codebox/config.RecommendedCLIVersion=$(RECOMMENDED_CLI_VERSION)
.ONESHELL:
build:
	go build -o bin/codebox -ldflags "$(LDFLAGS)" main.go
