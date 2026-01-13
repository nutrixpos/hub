package common

import (
	"fmt"
	"regexp"
	"strings"
)

func MaskString(data string, secrets []string) string {
	for _, secret := range secrets {
		data = strings.ReplaceAll(data, secret, "***********")
	}
	return data
}

func InterpretVars(plain string, vars_basket map[string]interface{}) (interpreted string, err error) {

	re := regexp.MustCompile(`\{\{\s*(.*?)\s*\}\}`)
	matches := re.FindAllStringSubmatch(plain, -1)

	var var_names []string
	var var_values []interface{}

	for _, match := range matches {
		if len(match) > 1 {
			var_names = append(var_names, match[1])
		}
	}

	for _, var_name := range var_names {

		found := false

		for basket_var_name, basket_var := range vars_basket {

			if basket_var_name == var_name {
				var_values = append(var_values, basket_var)
				found = true
				break
			}
		}

		if !found {
			return interpreted, fmt.Errorf("requested basket variable not found: %s", var_name)
		}
	}

	interpreted = plain
	for _, value := range var_values {

		re := regexp.MustCompile(`\{\{\s*(.*?)\s*\}\}`)

		// Find the first match and its submatch
		loc := re.FindStringSubmatchIndex(interpreted)

		if loc == nil {
			break
		}

		// loc[0], loc[1] = start/end of full match
		// loc[2], loc[3] = start/end of first submatch (content inside {{}})

		start, end := loc[0], loc[1]
		interpreted = fmt.Sprintf("%s%s%s", interpreted[:start], value, interpreted[end:])
	}

	return interpreted, nil
}
