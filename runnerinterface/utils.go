package runnerinterface

import "strconv"

/*
validates a permission string to ensure it is a valid octal
representation of file permissions (e.g., "755").
*/
func validatePermissionString(perm string) bool {
	if len(perm) != 3 {
		return false
	}

	// convert to octal and check if it's a valid permission
	octalPerm, err := strconv.ParseUint(perm, 8, 32)
	if err != nil {
		return false
	}

	if octalPerm > 0777 {
		return false
	}

	return true
}
