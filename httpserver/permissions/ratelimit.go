package permissions

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/cache"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
)

/*
IPRateLimitedRoute wraps a Gin route handler with IP-based rate limiting.

The rate limit is enforced per client IP address and request path. Each incoming
request creates a cache key with a time-based suffix and a TTL representing the
rate-limit window.

TTL Calculation Algorithm:
  - Under normal conditions, each request key is created with a base TTL equal
    to `periodSeconds`.
  - When the number of requests exceeds `callsPerPeriod`, the TTL of newly
    created keys is increased linearly based on the number of excess requests:
    TTL = periodSeconds + (excessRequests * periodSeconds)
  - If the request count exceeds five times the allowed limit, an additional
    penalty is applied by multiplying the calculated TTL by 10.
  - This escalating TTL acts as a progressive backoff mechanism, increasing
    the cooldown period as abuse intensity increases.
*/
func IPRateLimitedRoute(
	handler gin.HandlerFunc,
	callsPerPeriod int,
	defaultPeriodSeconds int,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ipAddress := c.ClientIP()
		requestPath := c.FullPath()

		baseKey := fmt.Sprintf(
			"%s-%s",
			ipAddress,
			requestPath,
		)

		ratelimitExceeded := false

		// count items matching that key
		keys, err := cache.GetKeysByPatternFromCache(fmt.Sprintf("%s*", baseKey))
		if err != nil {
			// TODO: log error
			utils.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"unknown error",
			)
			return
		}

		keyTTL := defaultPeriodSeconds

		if len(keys) > callsPerPeriod {
			ratelimitExceeded = true

			// increase TTL based on how much the limit is exceeded
			excess := len(keys) - callsPerPeriod
			keyTTL = defaultPeriodSeconds + (excess * defaultPeriodSeconds)

			if len(keys) > 5*callsPerPeriod {
				keyTTL *= 10
			}
		}

		// set key
		err = cache.SetKeyToCache(
			fmt.Sprintf(
				"%s-%d",
				baseKey,
				int(time.Now().UnixMilli()),
			),
			[]byte(""),
			keyTTL,
		)

		if err != nil {
			// TODO: log error
			utils.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"unknown error",
			)
			return
		}

		if ratelimitExceeded {
			utils.ErrorResponse(
				c,
				http.StatusTooManyRequests,
				"ratelimit exceeded, try again later",
			)
		} else {
			handler(c)
		}
	}
}
