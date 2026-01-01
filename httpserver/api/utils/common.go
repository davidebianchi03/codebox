package utils

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
Get an unsigned integer parameter from the Gin context.
Returns an error if the parameter is not found or is not a valid unsigned integer.
*/
func GetUIntParamFromContext(ctx *gin.Context, paramName string) (uint, error) {
	paramValue, found := ctx.Params.Get(paramName)
	if !found {
		return 0, fmt.Errorf("parameter %s not found", paramName)
	}

	paramUint, err := strconv.ParseUint(paramValue, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid parameter %s: %w", paramName, err)
	}

	return uint(paramUint), nil
}
