package utils

func RecursiveCatchField(field string, data map[string]interface{}) interface{} {
	if data[field] != nil {
		return data[field]
	}

	for _, value := range data {
		if value != nil {
			switch value := value.(type) {
			case []interface{}:
				for _, v := range value {
					if nestedMap, ok := v.(map[string]interface{}); ok {
						result := RecursiveCatchField(field, nestedMap)
						if result != nil {
							return result
						}
					}
				}
			}
		}
	}

	return nil
}
