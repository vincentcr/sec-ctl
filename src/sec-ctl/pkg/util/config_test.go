package util

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/vincentcr/testify/assert"
)

const testKeyPfx = "ConfTest"

type confTestSimple struct {
	I1 int
	S1 string
	S2 string
	F1 float32
	F2 float32
}

func (cfg *confTestSimple) AppName() string {
	return testKeyPfx + "Simple"
}

type confTestComplex struct {
	I1  int
	S1  string
	ST1 confTestSimple
	A1  []int
	M1  map[int]string
	A2  []confTestSimple
	M2  map[string]confTestSimple
	A3  []string
}

func (cfg *confTestComplex) AppName() string {
	return testKeyPfx + "Complex"
}

func mkSimple() confTestSimple {
	return confTestSimple{I1: 1, S1: "abc", F1: 123.0}
}

func mkComplex() confTestComplex {

	return confTestComplex{
		I1:  123,
		S1:  "abcdefgh",
		ST1: mkSimple(),
		A1:  []int{-1, 3, -5, -7},
		M1:  map[int]string{1: "abcdef", 2: "klmnop", 4: "defghi"},
		A2:  []confTestSimple{mkSimple(), mkSimple(), mkSimple()},
		M2: map[string]confTestSimple{
			"a": mkSimple(),
			"b": mkSimple(),
			"Z": mkSimple(),
		},
	}
}

func clearEnv() {
	for _, env := range os.Environ() {
		kv := strings.Split(env, "=")
		k := kv[0]
		if strings.HasPrefix(k, envKeyPrefix+testKeyPfx) {
			os.Unsetenv(k)
		}
	}
}

func TestLoadConfigDefaultsOnly(t *testing.T) {
	clearEnv()
	defaults := mkSimple()
	cfg := confTestSimple{}

	if err := LoadConfig(&cfg, &defaults); err != nil {
		t.Fatalf("LoadConfig errored out: %v", err)
	}

	assert.Equal(t, defaults, cfg)
}

func TestLoadConfigDefaultsWithEnv(t *testing.T) {
	clearEnv()
	defaults := mkSimple()
	cfg := confTestSimple{}

	os.Setenv(envKeyPrefix+"ConfTestSimple.S2", "abc")
	os.Setenv(envKeyPrefix+"ConfTestSimple.F1", "3.141592")

	if err := LoadConfig(&cfg, &defaults); err != nil {
		t.Fatalf("LoadConfig errored out: %v", err)
	}

	defaults.S2 = "abc"
	defaults.F1 = 3.141592

	assert.Equal(t, defaults, cfg)
}

func TestLoadConfigDefaultsWithEnvNestedObject(t *testing.T) {
	clearEnv()
	defaults := mkComplex()
	cfg := confTestComplex{}

	os.Setenv(envKeyPrefix+"ConfTestComplex.I1", "-246")
	os.Setenv(envKeyPrefix+"ConfTestComplex.ST1.F1", "3.141592")
	os.Setenv(envKeyPrefix+"ConfTestComplex.ST1.F2", "2.718281")

	if err := LoadConfig(&cfg, &defaults); err != nil {
		t.Fatalf("LoadConfig errored out: %v", err)
	}

	defaults.I1 = -246
	defaults.ST1.F1 = 3.141592
	defaults.ST1.F2 = 2.718281

	assert.Equal(t, defaults, cfg)
}

func TestLoadConfigErrorsOnInvalidEnvKeyPath(t *testing.T) {
	clearEnv()
	defaults := mkSimple()
	cfg := confTestSimple{}

	os.Setenv(envKeyPrefix+"ConfTestSimple.FOO", "-246")

	if err := LoadConfig(&cfg, &defaults); err == nil {
		t.Fatalf("Expected LoadConfig to fail on invalid key")
	} else if !strings.Contains(err.Error(), "FOO") {
		t.Fatalf("Expected LoadConfig error to show the key: instead got: %v", err)
	}
}

func TestCreateConfigFromDefaults(t *testing.T) {
	defaults := mkSimple()
	cfg := confTestSimple{}

	if err := loadConfigFromDefaults(&cfg, &defaults); err != nil {
		t.Fatalf("loadConfigFromDefaults errored out: %v", err)
	}

	assert.Equal(t, defaults, cfg)
}

func TestDeepCopy(t *testing.T) {

	src := mkComplex()
	dst := confTestComplex{}

	if err := deepCopy(&dst, &src); err != nil {
		t.Fatalf("deepCopy errored out: %v", err)
	}

	t.Log("src =", src)
	t.Log("dst =", dst)

	assert.Equal(t, mkComplex(), dst)
}

func TestParseVal(t *testing.T) {

	type testCase struct {
		expected interface{}
		input    string
		kind     reflect.Kind
	}

	testCases := []testCase{
		testCase{"abc", "abc", reflect.String},
		testCase{"", "", reflect.String},
		testCase{true, "true", reflect.Bool},
		testCase{false, "false", reflect.Bool},

		testCase{int(-9223372036854775808), "-9223372036854775808", reflect.Int},
		testCase{int(-12345678), "-12345678", reflect.Int},
		testCase{int(-1), "-1", reflect.Int},
		testCase{int(0), "0", reflect.Int},
		testCase{int(1), "1", reflect.Int},
		testCase{int(12345678), "12345678", reflect.Int},
		testCase{int(9223372036854775807), "9223372036854775807", reflect.Int},

		testCase{int64(-9223372036854775808), "-9223372036854775808", reflect.Int64},
		testCase{int64(-12345678), "-12345678", reflect.Int64},
		testCase{int64(-1), "-1", reflect.Int64},
		testCase{int64(0), "0", reflect.Int64},
		testCase{int64(1), "1", reflect.Int64},
		testCase{int64(12345678), "12345678", reflect.Int64},
		testCase{int64(9223372036854775807), "9223372036854775807", reflect.Int64},

		testCase{int32(-2147483648), "-2147483648", reflect.Int32},
		testCase{int32(-12345678), "-12345678", reflect.Int32},
		testCase{int32(-1), "-1", reflect.Int32},
		testCase{int32(0), "0", reflect.Int32},
		testCase{int32(1), "1", reflect.Int32},
		testCase{int32(12345678), "12345678", reflect.Int32},
		testCase{int32(2147483647), "2147483647", reflect.Int32},

		testCase{int16(-32768), "-32768", reflect.Int16},
		testCase{int16(-1234), "-1234", reflect.Int16},
		testCase{int16(-1), "-1", reflect.Int16},
		testCase{int16(0), "0", reflect.Int16},
		testCase{int16(1), "1", reflect.Int16},
		testCase{int16(1234), "1234", reflect.Int16},
		testCase{int16(32767), "32767", reflect.Int16},

		testCase{int8(-128), "-128", reflect.Int8},
		testCase{int8(-12), "-12", reflect.Int8},
		testCase{int8(-1), "-1", reflect.Int8},
		testCase{int8(0), "0", reflect.Int8},
		testCase{int8(1), "1", reflect.Int8},
		testCase{int8(12), "12", reflect.Int8},
		testCase{int8(127), "127", reflect.Int8},

		testCase{uint(0), "0", reflect.Uint},
		testCase{uint(1), "1", reflect.Uint},
		testCase{uint(12345678), "12345678", reflect.Uint},
		testCase{uint(18446744073709551615), "18446744073709551615", reflect.Uint},

		testCase{uint64(0), "0", reflect.Uint64},
		testCase{uint64(1), "1", reflect.Uint64},
		testCase{uint64(12345678), "12345678", reflect.Uint64},
		testCase{uint64(18446744073709551615), "18446744073709551615", reflect.Uint64},

		testCase{uint32(0), "0", reflect.Uint32},
		testCase{uint32(1), "1", reflect.Uint32},
		testCase{uint32(12345678), "12345678", reflect.Uint32},
		testCase{uint32(4294967295), "4294967295", reflect.Uint32},

		testCase{uint16(0), "0", reflect.Uint16},
		testCase{uint16(1), "1", reflect.Uint16},
		testCase{uint16(1234), "1234", reflect.Uint16},
		testCase{uint16(65535), "65535", reflect.Uint16},

		testCase{uint8(0), "0", reflect.Uint8},
		testCase{uint8(1), "1", reflect.Uint8},
		testCase{uint8(127), "127", reflect.Uint8},
		testCase{uint8(255), "255", reflect.Uint8},

		testCase{float32(-3.40282346638528859811704183484516925440e+38), "-3.40282346638528859811704183484516925440e+38", reflect.Float32},
		testCase{float32(-12), "-12", reflect.Float32},
		testCase{float32(-1), "-1", reflect.Float32},
		testCase{float32(-1.401298464324817070923729583289916131280e-45), "-1.401298464324817070923729583289916131280e-45", reflect.Float32},
		testCase{float32(0), "0", reflect.Float32},
		testCase{float32(0), "0.0", reflect.Float32},
		testCase{float32(1.401298464324817070923729583289916131280e-45), "1.401298464324817070923729583289916131280e-45", reflect.Float32},
		testCase{float32(0.000001), "0.000001", reflect.Float32},
		testCase{float32(1), "1", reflect.Float32},
		testCase{float32(12), "12", reflect.Float32},
		testCase{float32(3.141592), "3.141592", reflect.Float32},
		testCase{float32(3.40282346638528859811704183484516925440e+38), "3.40282346638528859811704183484516925440e+38", reflect.Float32},

		testCase{float64(-1.797693134862315708145274237317043567981e+308), "-1.797693134862315708145274237317043567981e+308", reflect.Float64},
		testCase{float64(-12), "-12", reflect.Float64},
		testCase{float64(-1), "-1", reflect.Float64},
		testCase{float64(-4.940656458412465441765687928682213723651e-324), "-4.940656458412465441765687928682213723651e-324", reflect.Float64},
		testCase{float64(0), "0", reflect.Float64},
		testCase{float64(0), "0.0", reflect.Float64},
		testCase{float64(4.940656458412465441765687928682213723651e-324), "4.940656458412465441765687928682213723651e-324", reflect.Float64},
		testCase{float64(0.000001), "0.000001", reflect.Float64},
		testCase{float64(1), "1", reflect.Float64},
		testCase{float64(12), "12", reflect.Float64},
		testCase{float64(3.141592), "3.141592", reflect.Float64},
		testCase{float64(1.797693134862315708145274237317043567981e+308), "1.797693134862315708145274237317043567981e+308", reflect.Float64},
	}

	for _, tc := range testCases {
		a, err := parseVal(tc.input, tc.kind)
		if err != nil {
			t.Fatalf("parseVal failed for (%v,%v): %v", tc.input, tc.kind, err)
		}

		assert.Equal(t, tc.expected, a, "parseVal(%v, %v)", tc.input, tc.kind)
	}
}

func TestParseValRejectsInvalid(t *testing.T) {
	type testCase struct {
		input string
		kind  reflect.Kind
	}

	testCases := []testCase{
		testCase{"abc", reflect.Ptr},
		testCase{"abc", reflect.Interface},
		testCase{"-10", reflect.Uint},
	}

	for _, tc := range testCases {
		a, err := parseVal(tc.input, tc.kind)
		if err == nil {
			t.Fatalf("Expected parseVal to fail on (%v,%v): instead, got: %v", tc.input, tc.kind, a)
		}
	}

}
