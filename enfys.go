package enfys

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

func splitTag(tag string) (string, bool, error) {
	parts := strings.Split(tag, ",")
	if len(parts) == 0 {
		return "", false, fmt.Errorf("empty tag")
	}

	if len(parts) > 2 {
		return "", false, fmt.Errorf("invalid tag format")
	}

	if len(parts) == 2 {
		if parts[1] == "required" {
			return parts[0], true, nil
		} else {
			return "", false, fmt.Errorf("unknown tag option: %s", parts[1])
		}
	}

	return parts[0], false, nil
}

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

	missing := []string{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		envTag := field.Tag.Get("enfys")
		if envTag == "" {
			continue
		}

		key, required, err := splitTag(envTag)
		if err != nil {
			return fmt.Errorf("invalid tag for field %s: %v", field.Name, err)
		}

		envValue, exists := os.LookupEnv(key)
		if !exists {
			if required {
				missing = append(missing, key)
			}

			defaultTag := field.Tag.Get("default")
			if defaultTag != "" {
				envValue = defaultTag
			} else {
				continue
			}
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

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return nil
}
