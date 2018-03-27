package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
	"reflect"
	"strings"
)

func ParseCommandLine(v *viper.Viper) {
	var flags = pflag.CommandLine

	// bind all the command line arguments to the flags
	bindCommandLine(reflect.ValueOf(getDefault()).Elem().Type(), flags, "")

	// push the flags to viper
	v.BindPFlags(flags)

	// parse our command line arguments
	flags.Parse(os.Args[1:])
}

func bindCommandLine(defaultType reflect.Type, flags *pflag.FlagSet, path string) {
	for i := 0; i < defaultType.NumField(); i++ {
		var field = defaultType.Field(i)
		var key = strings.ToLower(field.Name)
		var fieldType = field.Type

		// build our current path
		var currentPath = appendPath(path, key)

		bindCommandLineField(fieldType, flags, currentPath)
	}
}

func bindCommandLineField(defaultType reflect.Type, flags *pflag.FlagSet, path string) {
	if defaultType.Kind() == reflect.Invalid {
		return
	}

	switch defaultType.Kind() {
	case reflect.Ptr:
		bindCommandLineField(defaultType.Elem(), flags, path)
	case reflect.Struct, reflect.Interface:
		bindCommandLine(defaultType, flags, path)

	case reflect.Bool:
		flags.Bool(path, false, "")

	case reflect.Int, reflect.Int16:
		flags.Int(path, 0, "")
	case reflect.Int8:
		flags.Int8(path, 0, "")
	case reflect.Int32:
		flags.Int32(path, 0, "")
	case reflect.Int64:
		flags.Int64(path, 0, "")

	case reflect.Uint:
		flags.Uint(path, 0, "")
	case reflect.Uint8:
		flags.Uint8(path, 0, "")
	case reflect.Uint16:
		flags.Uint16(path, 0, "")
	case reflect.Uint32:
		flags.Uint32(path, 0, "")
	case reflect.Uint64:
		flags.Uint64(path, 0, "")

	case reflect.Float32:
		flags.Float32(path, 0, "")
	case reflect.Float64:
		flags.Float64(path, 0, "")

	case reflect.String:
		flags.String(path, "", "")

	default:
		log.Printf("Unexpected type %s at %s", defaultType.String(), path)
	}
}
