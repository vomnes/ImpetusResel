package utils

import "reflect"

// StringInArray take a string or an array of strings and a array of string as parameter
// If the first argument is an array of string, the function return true if
// at list one of the elements array in the array of the second argument
// Return true if the string is in the array of string else false
func StringInArray(a interface{}, list []string) bool {
	var elements []string
	typeA := reflect.TypeOf(a)
	if typeA.String() == "string" {
		elements = []string{a.(string)}
	}
	if typeA.String() == "[]string" {
		elements = a.([]string)
	}
	if len(elements) == 0 {
		return false
	}
	for _, b := range list {
		for _, elem := range elements {
			if b == elem {
				return true
			}
		}
	}
	return false
}
