package utils

func RecursiveCatchField(field string, data map[string]interface{}) string {
	if data[field] != nil {
		if str, ok := data[field].(string); ok {
			return str
		}
	}

	for _, value := range data {
		if value != nil {
			switch value := value.(type) {
			case []interface{}:
				for _, v := range value {
					if nestedMap, ok := v.(map[string]interface{}); ok {
						result := RecursiveCatchField(field, nestedMap)
						if result != "" {
							return result
						}
					}
				}
			}
		}
	}

	return ""
}
