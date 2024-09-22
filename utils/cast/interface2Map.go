package cast

import (
	"fmt"
)

func Interface2StringMap(in interface{}) (map[string]interface{}, error) {
	inMap, ok := in.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot parse generic map to string map")
	}
	out := make(map[string]interface{})
	for k, v := range inMap {
		key, ok := k.(string)
		if ok {
			out[key] = v
		} else {
			return nil, fmt.Errorf("cannot parse generic map to string map")
		}
	}
	return out, nil
}
