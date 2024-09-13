package main

import (
	"encoding/json"
	"errors"
	"strings"
)

func ToJsonMulti(tenders any) (string, error) {
	jsonData, err := json.Marshal(tenders)
	if err != nil {
		return "", errors.New("error while formating query result to JSON")
	}

	return string(jsonData), nil
}

func ToJson(tenders any) (string, error) {
	jsonData, err := json.Marshal(tenders)
	if err != nil {
		return "", errors.New("error while formating query result to JSON")
	}

	return string(jsonData), nil
}

func UpdateMapFromAnother(map1, map2 map[string]string) map[string]string {
	for key1 := range map1 {
		for key2, value2 := range map2 {
			if strings.EqualFold(key1, key2) {
				map1[key1] = value2
				break
			}
		}
	}
	return map1
}
