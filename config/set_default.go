package config

import (
	"github.com/spf13/viper"
	"log"
	"reflect"
	"strings"
)

func SetDefault(v *viper.Viper) {
	setDefaultObject(reflect.ValueOf(getDefault()).Elem(), v, "")
}

func setDefaultObject(defaultValue reflect.Value, v *viper.Viper, path string) {
	var defaultType = defaultValue.Type()

	for i := 0; i < defaultType.NumField(); i++ {
		var field = defaultType.Field(i)
		var key = strings.ToLower(field.Name)
		var value = defaultValue.Field(i)

		// build our current path
		var currentPath = appendPath(path, key)

		setDefaultField(value, v, currentPath)
	}
}

func setDefaultSlice(defaultConfig reflect.Value, v *viper.Viper, path string) {
	log.Fatalln("not implemented")
}

func setDefaultField(value reflect.Value, v *viper.Viper, path string) {
	if value.Kind() == reflect.Invalid {
		return
	}

	switch value.Type().Kind() {
	case reflect.Interface, reflect.Ptr:
		setDefaultField(value.Elem(), v, path)
	case reflect.Struct:
		setDefaultObject(value, v, path)
	case reflect.Array, reflect.Slice:
		setDefaultSlice(value, v, path)

	case reflect.Bool:
		v.SetDefault(path, value.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetDefault(path, value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetDefault(path, value.Uint())
	case reflect.Float32, reflect.Float64:
		v.SetDefault(path, value.Float())
	case reflect.Complex64, reflect.Complex128:
		v.SetDefault(path, value.Complex())
	case reflect.String:
		v.SetDefault(path, value.String())
	default:
		log.Printf("Unexpected type %s at %s", value.Type().String(), path)
	}
}
