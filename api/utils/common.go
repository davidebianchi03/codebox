package utils

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
