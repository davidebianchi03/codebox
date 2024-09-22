package cast

import "fmt"

func Interface2StringArray(i interface{}) ([]string, error) {
	var stringArray []string

	intArray, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot convert interface to []string")
	}

	for _, value := range intArray {
		strValue, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("interface is not a string array")
		}
		stringArray = append(stringArray, strValue)
	}

	return stringArray, nil

}
