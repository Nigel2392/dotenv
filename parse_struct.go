package dotenv

import (
	"reflect"
	"strconv"
	"strings"
)

func (e *Env) unmarshal(s ...any) error {
	for _, structie := range s {
		// Parse the form data into the struct
		structValue := reflect.ValueOf(structie)
		if structValue.Kind() == reflect.Ptr {
			structValue = structValue.Elem()
		}
		if structValue.Kind() != reflect.Struct {
			panic("not a struct")
		}
		structType := structValue.Type()
		for key, value := range e.Variables {
			var keys = strings.Split(key, ".")
			structName := structType.Name()
			if strings.EqualFold(structName, keys[0]) {
				for i := 0; i < structValue.NumField(); i++ {
					field := structValue.Field(i)
					//	if field.Kind() == reflect.Pointer {
					//		// Create a new value for the pointer and set the field to it
					//		field.Set(reflect.New(field.Type().Elem()))
					//		field = field.Elem()
					//	}

					if field.Kind() == reflect.Struct || field.Kind() == reflect.Ptr {
						//	if err := e.unmarshal(field.Addr().Interface()); err != nil {
						//		return err
						//	}
						continue
					} else {
						if strings.EqualFold(structType.Field(i).Tag.Get("env"), keys[1]) {
							var casted = castType(field, value)
							field.Set(reflect.ValueOf(casted))
						}
					}
				}
			}
		}
	}
	return nil
}

func castType(f reflect.Value, val []string) any {
	switch f.Kind() {
	case reflect.String:
		return val[0]
	case reflect.Bool:
		var b, err = strconv.ParseBool(val[0])
		if err != nil {
			panic(err)
		}
		return b
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i, err = strconv.Atoi(val[0])
		if err != nil {
			panic(err)
		}
		return i
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var i, err = strconv.Atoi(val[0])
		if err != nil {
			panic(err)
		}
		return uint(i)
	case reflect.Slice:
		var slice = reflect.MakeSlice(f.Type(), len(val), len(val))
		for i, v := range val {
			slice.Index(i).SetString(v)
		}
		return slice.Interface()
	default:
		panic("unsupported type")
	}
}
