package enfys

import (
	"fmt"
	"os"
	"reflect"
)

func isTrue(value string) bool {
	accept := []string{
		"true",
		"t",
		"yes",
		"y",
		"on",
		"1",
	}

	for _, a := range accept {
		if value == a {
			return true
		}
	}

	return false
}

func Fill(s any) error {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected pointer to struct, got %s", t.Kind())
	}

	t = t.Elem()
	v = v.Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		envTag := field.Tag.Get("enfys")
		if envTag == "" {
			continue
		}

		envValue, exists := os.LookupEnv(envTag)
		if !exists {
			continue
		}

		if !value.CanSet() {
			return fmt.Errorf("cannot set field %s", field.Name)
		}

		switch value.Kind() {
		case reflect.Int:
			var intValue int
			_, err := fmt.Sscanf(envValue, "%d", &intValue)
			if err != nil {
				return fmt.Errorf("failed to parse int for field %s: %v", field.Name, err)
			}
			value.SetInt(int64(intValue))
		case reflect.Bool:
			value.SetBool(isTrue(envValue))
		case reflect.String:
			value.SetString(envValue)
		case reflect.Float32, reflect.Float64:
			var floatValue float64
			_, err := fmt.Sscanf(envValue, "%f", &floatValue)
			if err != nil {
				return fmt.Errorf("failed to parse float for field %s: %v", field.Name, err)
			}
			value.SetFloat(floatValue)
		default:
			return fmt.Errorf("unsupported field type %s for field %s", value.Kind(), field.Name)
		}
	}

	return nil
}
