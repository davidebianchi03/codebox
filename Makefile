SERVER_VERSION ?= dev
RECOMMENDED_RUNNER_VERSION ?= latest
RECOMMENDED_CLI_VERSION ?= latest
ANALYTICS_ENDPOINT ?=
ANALYTICS_API_KEY ?=

LDFLAGS = -X gitlab.com/codebox4073715/codebox/config.ServerVersion=$(SERVER_VERSION) \
          -X gitlab.com/codebox4073715/codebox/config.RecommendedRunnerVersion=$(RECOMMENDED_RUNNER_VERSION) \
          -X gitlab.com/codebox4073715/codebox/config.RecommendedCLIVersion=$(RECOMMENDED_CLI_VERSION)
          -X gitlab.com/codebox4073715/codebox/config.AnalyticsEndpoint=$(ANALYTICS_ENDPOINT)
          -X gitlab.com/codebox4073715/codebox/config.AnalyticsApiKey=$(ANALYTICS_API_KEY)
.ONESHELL:
build:
	go build -o bin/codebox -ldflags "$(LDFLAGS)" main.go
