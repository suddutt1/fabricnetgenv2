package util

import (
	"fmt"

	"strconv"
)

func GetMap(element interface{}) map[string]interface{} {
	retMap, ok := element.(map[string]interface{})
	if ok == true {
		return retMap
	}

	return nil
}
func GetString(element interface{}) string {
	retString, ok := element.(string)
	if ok == true {
		return retString
	}
	return ""
}
func GetNumber(element interface{}) int {

	s := fmt.Sprintf("%v", element)
	retString, err := strconv.Atoi(s)
	if err == nil {
		return retString
	}
	return 0
}
func GetBoolean(element interface{}) bool {
	retString, ok := element.(string)
	if ok == true {
		retBool, inValid := strconv.ParseBool(retString)
		if inValid == nil {
			return retBool
		}
	}
	return false
}
func IfEntryExistsInMap(mapToCheck map[string]interface{}, attribute string) bool {
	_, ok := mapToCheck[attribute]
	return ok
}
