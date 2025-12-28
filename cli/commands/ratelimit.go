package commands

import (
	"fmt"

	"gitlab.com/codebox4073715/codebox/cache"
)

/*
This function handles the command to reset ratelimits
*/
func HandleResetRatelimits() uint {
	// list ratelimit keys and remove them
	fmt.Println("Listing ratelimit entries...")
	ratelimitKeys, err := cache.GetKeysByPatternFromCache("ratelimit-*")
	if err != nil {
		fmt.Printf(
			"Cannot list ratelimit keys, %s\n", err,
		)
		return 1
	}

	fmt.Printf("%d ratelimit entries have been found, removing them...\n", len(ratelimitKeys))

	for _, key := range ratelimitKeys {
		if err := cache.DeleteKeyFromCache(key); err != nil {
			fmt.Printf("Failed to delete an entry, %s\n", err)
			return 1
		}
	}

	fmt.Printf("%d ratelimit entries have been removed\n", len(ratelimitKeys))

	// list violations keys and remove them
	fmt.Println("Listing violation entries...")

	violationKeys, err := cache.GetKeysByPatternFromCache("violation-*")
	if err != nil {
		fmt.Printf(
			"Cannot list violation keys, %s\n", err,
		)
		return 1
	}

	fmt.Printf("%d violation entries have been found, removing them...\n", len(violationKeys))

	for _, key := range violationKeys {
		if err := cache.DeleteKeyFromCache(key); err != nil {
			fmt.Printf("Failed to delete an entry, %s\n", err)
			return 1
		}
	}

	fmt.Printf("%d violation entries have been removed\n", len(violationKeys))

	return 0
}
