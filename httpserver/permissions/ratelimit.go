package permissions

import (
	"fmt"
	"math"
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
  - Each request key is created with a base TTL equal
    to `periodSeconds` multiplied for the power of two the number of
    violations in last 24hrs
  - When the number of requests exceeds `callsPerPeriod`, a violation record is created
*/
func IPRateLimitedRoute(
	handler gin.HandlerFunc,
	callsPerPeriod int,
	defaultPeriodSeconds int,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ipAddress := c.ClientIP()
		requestPath := c.FullPath()

		// TODO: ignore whitelisted ip addresses

		baseKey := fmt.Sprintf(
			"ratelimit-%s-%s",
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

			// count how many times the ratelimit has been previously hitted
			violationsBaseKey := fmt.Sprintf(
				"violation-%s-%s",
				ipAddress,
				requestPath,
			)
			violations, err := cache.GetKeysByPatternFromCache(fmt.Sprintf("%s*", violationsBaseKey))
			if err != nil {
				// TODO: log error
				utils.ErrorResponse(
					c,
					http.StatusInternalServerError,
					"unknown error",
				)
				return
			}
			violationsCount := len(violations)

			if excess == 1 {
				// record violation only once per burst,
				// violations are recorded for 24 hours
				err = cache.SetKeyToCache(
					fmt.Sprintf(
						"%s-%d",
						violationsBaseKey,
						int(time.Now().UnixMilli()),
					),
					[]byte(""),
					24*60*60,
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

				if violationsCount == 2 {
					// this is the third violation (violations are enumerated before recording this one),
					// notify it to the admins
					// TODO: send an email to administrators if this setting is enabled
				}
			}

			// calculate the key ttl multiplier using
			// considering violations in last 24hrs
			multiplier := int(math.Min(
				math.Pow(2, float64(violationsCount)),
				1024, // max value for mulitplier 2^1024 = ~17h
			))
			keyTTL *= multiplier
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
				"too many requests, try again later",
			)
		} else {
			handler(c)
		}
	}
}
