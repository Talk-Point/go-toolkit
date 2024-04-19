package cli

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func Input(model interface{}, args []string) error {
	argMap := ParseArgs(args)
	return InputFromModel(model, argMap)
}

func ParseArgs(args []string) map[string]string {
	argMap := make(map[string]string)

	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				argMap[strings.TrimPrefix(args[i], "-")] = args[i+1]
				i++
			} else {
				argMap[strings.TrimPrefix(args[i], "-")] = ""
			}
		}
	}

	return argMap
}

func InputFromModel(model interface{}, args map[string]string) error {
	reader := bufio.NewReader(os.Stdin)
	val := reflect.ValueOf(model).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)

		tag := fieldType.Tag.Get("validate")
		if !strings.Contains(tag, "required") {
			continue
		}

		input, ok := args[strings.ToLower(fieldType.Name)]
		if !ok {
			fmt.Printf("Enter %s: ", fieldType.Name)
			inputValue, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("error reading input: %w", err)
			}
			input = strings.TrimSpace(inputValue)
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(input)
		case reflect.Int:
			i, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("error parsing int: %w", err)
			}
			field.SetInt(int64(i))
		case reflect.Ptr:
			if field.Type().Elem().Kind() == reflect.String {
				str := input
				field.Set(reflect.ValueOf(&str))
			} else if field.Type().Elem().Kind() == reflect.Int {
				i, err := strconv.Atoi(input)
				if err != nil {
					return fmt.Errorf("error parsing int: %w", err)
				}
				field.Set(reflect.ValueOf(&i))
			} else {
				fmt.Printf("Unsupported type: %s\n", field.Kind())
				return fmt.Errorf("unsupported type: %s", field.Kind())
			}
		default:
			fmt.Printf("Unsupported type: %s\n", field.Kind())
			return fmt.Errorf("unsupported type: %s", field.Kind())
		}
	}

	return nil
}
