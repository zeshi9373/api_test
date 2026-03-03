package fn

import (
	"encoding/json"
	"fmt"
)

func SearchListValue(data string, searchKey string, searchValue any, getKey string) any {
	dataMap := make([]map[string]any, 0)
	json.Unmarshal([]byte(data), &dataMap)
	if len(dataMap) > 0 {
		for _, item := range dataMap {
			if fmt.Sprintf("%v", item[searchKey]) == fmt.Sprintf("%v", searchValue) {
				return item[getKey]
			}
		}
	}

	return nil
}
