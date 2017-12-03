package util

/*

requirements:

- handle defaults
- handle environment variables, for production and docker compose
- handle different configurations for dev, test, prod, staging...
- store config to strongly typed object

*/

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const envKeyPrefix = "SecCtl."

type config interface {
	AppName() string
}

// LoadConfig loads the supplied configuration object from defaults and environment variables
func LoadConfig(cfg, defaults config) error {

	if err := loadConfigFromDefaults(cfg, defaults); err != nil {
		return err
	}

	if err := loadConfigFromEnv(cfg); err != nil {
		return err
	}

	dumpConfig(os.Stdout, cfg)

	return nil
}

func dumpConfig(w io.Writer, cfg config) {

	w.Write([]byte("Loaded config:\n  "))

	if dat, err := json.MarshalIndent(cfg, "  ", "  "); err != nil {
		panic(err)
	} else if _, err := w.Write(dat); err != nil {
		panic(err)
	}

	w.Write([]byte("\n"))
}

func loadConfigFromDefaults(cfg, defaults config) error {

	if err := deepCopy(cfg, defaults); err != nil {
		return err
	}

	return nil
}

func deepCopy(dst, src interface{}) error {
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(w)
	if err := enc.Encode(src); err != nil {
		return err
	}

	dec := gob.NewDecoder(r)
	return dec.Decode(dst)
}

func loadConfigFromEnv(cfg config) error {
	env := env2map()
	pfx := envKeyPrefix + cfg.AppName() + "."

	for k, v := range env {
		if strings.HasPrefix(k, pfx) {
			keyPath := strings.Split(k[len(pfx):], ".")
			if err := loadEnvVar(cfg, keyPath, v); err != nil {
				return err
			}
		}
	}

	return nil
}

func env2map() map[string]string {

	env := map[string]string{}

	for _, e := range os.Environ() {
		kv := strings.Split(e, "=")
		k := kv[0]
		v := kv[1]
		env[k] = v
	}
	return env
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
	case reflect.Uint:
		i, e := strconv.ParseUint(val, 10, 64)
		return uint(i), e
	case reflect.Int:
		i, e := strconv.ParseInt(val, 10, 64)
		return int(i), e
	case reflect.Float64:
		f, e := strconv.ParseFloat(val, 64)
		return float64(f), e
	case reflect.Float32:
		f, e := strconv.ParseFloat(val, 32)
		return float32(f), e
	default:
		return nil, fmt.Errorf("Unexpected kind %v for val %v", kind, val)
	}
}
