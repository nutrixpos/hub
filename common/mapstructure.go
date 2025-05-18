package common

import (
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
)

func StringToTimeHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}
		if s, ok := data.(string); ok {
			return time.Parse(time.RFC3339, s)
		}
		return data, nil
	}
}
