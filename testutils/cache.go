package testutils

import (
	"fmt"

	"gitlab.com/codebox4073715/codebox/cache"
)

/*
Remove all items from cache, this function removes:
- ratelimit keys
*/
func ClearCache() error {
	// list ratelimit keys and remove them
	ratelimitKeys, err := cache.GetKeysByPatternFromCache("ratelimit-*")
	if err != nil {
		return fmt.Errorf(
			"cannot list ratelimit keys, %s\n", err,
		)
	}
	for _, key := range ratelimitKeys {
		if err := cache.DeleteKeyFromCache(key); err != nil {
			return fmt.Errorf("failed to delete an entry, %s\n", err)
		}
	}
	// list violations keys and remove them
	violationKeys, err := cache.GetKeysByPatternFromCache("violation-*")
	if err != nil {
		return fmt.Errorf(
			"cannot list violation keys, %s", err,
		)
	}

	for _, key := range violationKeys {
		if err := cache.DeleteKeyFromCache(key); err != nil {
			return fmt.Errorf("failed to delete an entry, %s\n", err)
		}
	}

	return nil
}
