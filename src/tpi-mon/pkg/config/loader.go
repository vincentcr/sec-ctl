package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/imdario/mergo"
)

const envKeyPrefix = "Tpimon."

// Load loads the config from default sources
func Load() (Config, error) {
	config := Config{}

	loaders := []func() (Config, error){
		func() (Config, error) { return defaultConfig, nil },
		loadFromEnv,
	}

	for _, loader := range loaders {
		src, err := loader()
		if err != nil {
			return Config{}, err
		}
		if err = mergo.MergeWithOverwrite(&config, src); err != nil {
			return Config{}, err
		}
	}

	log.Printf("Loaded config: %#v\n", config)

	return config, nil
}

func loadFromEnv() (Config, error) {

	config := Config{}

	for _, e := range os.Environ() {
		kv := strings.Split(e, "=")
		k := kv[0]
		v := kv[1]
		if strings.HasPrefix(k, envKeyPrefix) {
			keyPath := strings.Split(k, ".")[1:]
			if err := loadEnvVar(&config, keyPath, v); err != nil {
				return Config{}, err
			}
		}
	}

	return config, nil
}

func loadEnvVar(dstPtr interface{}, keyPath []string, val string) error {

	var fldVal reflect.Value

	// navigate through dstPtr with keyPath
	for _, k := range keyPath {

		// get reflect.Value of dereferenced dstPtr
		dstVal := reflect.ValueOf(dstPtr).Elem()

		// verify that k is a valid field of dstVal
		if _, ok := dstVal.Type().FieldByName(k); !ok {
			return fmt.Errorf("Invalid sub key %v in key path %v", k, keyPath)
		}

		// get reflect.Value of dstVal.k
		fldVal = dstVal.FieldByName(k)

		// update dstPtr for next loop
		dstPtr = fldVal.Addr().Interface()
	}

	// last step is to convert val from string to destination type
	converted, err := parseVal(val, fldVal.Type().Kind())
	if err != nil {
		return err
	}
	fldVal.Set(reflect.ValueOf(converted))

	return nil
}

func parseVal(val string, kind reflect.Kind) (interface{}, error) {
	switch kind {
	case reflect.String:
		return val, nil
	case reflect.Bool:
		return strconv.ParseBool(val)
	case reflect.Int8:
		i, e := strconv.ParseInt(val, 10, 8)
		return int8(i), e
	case reflect.Int16:
		i, e := strconv.ParseInt(val, 10, 16)
		return int16(i), e
	case reflect.Int32:
		i, e := strconv.ParseInt(val, 10, 32)
		return int32(i), e
	case reflect.Int64:
		i, e := strconv.ParseInt(val, 10, 64)
		return int64(i), e
	case reflect.Uint8:
		i, e := strconv.ParseUint(val, 10, 8)
		return uint8(i), e
	case reflect.Uint16:
		i, e := strconv.ParseUint(val, 10, 16)
		return uint16(i), e
	case reflect.Uint32:
		i, e := strconv.ParseUint(val, 10, 32)
		return uint32(i), e
	case reflect.Uint64:
		i, e := strconv.ParseUint(val, 10, 64)
		return uint64(i), e
	case reflect.Int:
		i, e := strconv.ParseUint(val, 10, 64)
		return int(i), e
	case reflect.Float64:
		f, e := strconv.ParseFloat(val, 64)
		return float64(f), e
	case reflect.Float32:
		f, e := strconv.ParseFloat(val, 32)
		return float32(f), e
	default:
		panic(fmt.Errorf("Unexpected kind %v for val %v", kind, val))
	}
}
