package util

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
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

	assertConfTestSimpleEqual(t, defaults, cfg, "defaults only")
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

	assertConfTestSimpleEqual(t, defaults, cfg, "defaults with env")
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

	assertconfTestComplexEqual(t, defaults, cfg, "nested env vars")
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

	assertConfTestSimpleEqual(t, defaults, cfg, "loadConfigFromDefaults")
}

func TestDeepCopy(t *testing.T) {

	src := mkComplex()
	dst := confTestComplex{}

	if err := deepCopy(&dst, &src); err != nil {
		t.Fatalf("deepCopy errored out: %v", err)
	}

	t.Log("src =", src)
	t.Log("dst =", dst)

	assertconfTestComplexEqual(t, mkComplex(), dst, "complex struct")
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

		assertEqual(t, tc.expected, a, "parseVal(%v, %v)", tc.input, tc.kind)
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

func assertConfTestSimpleEqual(t *testing.T, e, a confTestSimple, msgFmt string, msgArgs ...interface{}) {
	t.Helper()

	msg := fmt.Sprintf(msgFmt, msgArgs...)

	assertEqual(t, e.I1, a.I1, "%v: I1", msg)
	assertEqual(t, e.S1, a.S1, "%v: S1", msg)
	assertEqual(t, e.S2, a.S2, "%v: S2", msg)
	assertEqual(t, e.F1, a.F1, "%v: F1", msg)
}

func assertconfTestComplexEqual(t *testing.T, e, a confTestComplex, msgFmt string, msgArgs ...interface{}) {
	t.Helper()

	msg := fmt.Sprintf(msgFmt, msgArgs...)

	assertEqual(t, e.I1, a.I1, "%v: I1", msg)
	assertEqual(t, e.S1, a.S1, "%v: S1", msg)

	assertConfTestSimpleEqual(t, e.ST1, a.ST1, "%v: ST1", msg)

	assertNotEqual(t, &e.A1, &a.A1, "%v: A1 should not be same ref", msg)

	assertEqual(t, len(e.A1), len(a.A1), "%v: A1 len", msg)
	for i, v := range e.A1 {
		assertEqual(t, v, a.A1[i], "%v: A1[%v]", msg, i)
	}

	assertNotEqual(t, &e.M1, &a.M1, "%v: M1 should not be same ref", msg)
	assertEqual(t, len(e.M1), len(a.M1), "%v: M1 len", msg)
	for k, v := range e.M1 {
		assertEqual(t, v, a.M1[k], "%v: M1[%v]", msg, k)
	}

	assertNotEqual(t, &e.A2, &a.A2, "%v: A2 should not be same ref", msg)
	assertEqual(t, len(e.A2), len(a.A2), "%v: A2 len", msg)
	for i, v := range e.A2 {
		assertConfTestSimpleEqual(t, v, a.A2[i], "%v: A2[%v]", msg, i)
	}

	assertNotEqual(t, &e.M2, &a.M2, "%v: M2 should not be same ref", msg)
	assertEqual(t, len(e.M2), len(a.M2), "%v: M2 len", msg)
	for k, v := range e.M2 {
		assertConfTestSimpleEqual(t, v, a.M2[k], "%v: M2[%v]", msg, k)
	}

}

func assertNotEqual(t *testing.T, e interface{}, a interface{}, msgFmt string, msgArgs ...interface{}) {
	t.Helper()

	if e == a {
		assertInfo := fmt.Sprintf("%#v == %#v", e, a)
		msg := fmt.Sprintf(msgFmt, msgArgs...)
		t.Fatalf("Assertion failed: %v\n\t%v", msg, assertInfo)
	}

}

func assertEqual(t *testing.T, e interface{}, a interface{}, msgFmt string, msgArgs ...interface{}) {
	t.Helper()

	eV := reflect.ValueOf(e)
	aV := reflect.ValueOf(a)

	if eV.Kind() != aV.Kind() {
		assertInfo := fmt.Sprintf("%v != %v", eV.Kind(), aV.Kind())
		msg := fmt.Sprintf(msgFmt, msgArgs...)
		t.Fatalf("Assertion failed: not same kind: %v\n\t%v", msg, assertInfo)
	}

	if eV.Type() != aV.Type() {
		assertInfo := fmt.Sprintf("%v != %v", eV.Type(), aV.Type())
		msg := fmt.Sprintf(msgFmt, msgArgs...)
		t.Fatalf("Assertion failed: not same type: %v\n\t%v", msg, assertInfo)
	}

	if e != a {
		assertInfo := fmt.Sprintf("%#v != %#v", e, a)
		msg := fmt.Sprintf(msgFmt, msgArgs...)
		t.Fatalf("Assertion failed: %v\n\t%v", msg, assertInfo)
	}
}

func assertSliceEqual(t *testing.T, e interface{}, a interface{}, msgFmt string, msgArgs ...interface{}) {
	t.Helper()

	msg := fmt.Sprintf(msgFmt, msgArgs...)

	eV := reflect.ValueOf(e)
	aV := reflect.ValueOf(a)

	if eV.Kind() != reflect.Slice {
		t.Fatalf("%v: should be a slice", msg)
	}

	assertEqual(t, eV.Type(), aV.Type(), "%v: type not equal", msg)
	assertEqual(t, eV.Len(), aV.Len(), "%v: len not equal", msg)

	for i := 0; i < eV.Len(); i++ {
		ei := eV.Index(i).Interface()
		ai := aV.Index(i).Interface()
		assertEqual(t, ei, ai, "%v: not equal at %v", msg, i)
	}
}

func assertMapEqual(t *testing.T, e interface{}, a interface{}, msgFmt string, msgArgs ...interface{}) {
	t.Helper()

	msg := fmt.Sprintf(msgFmt, msgArgs...)

	eV := reflect.ValueOf(e)
	aV := reflect.ValueOf(a)

	if eV.Kind() != reflect.Map {
		t.Fatalf("%v: should be a map", msg)
		t.Fail()
		return
	}

	assertEqual(t, eV.Type(), aV.Type(), "%v: type not equal", msg)
	assertEqual(t, eV.Len(), aV.Len(), "%v: len not equal", msg)

	for _, k := range eV.MapKeys() {
		ek := eV.MapIndex(k).Interface()
		ak := aV.MapIndex(k).Interface()
		assertEqual(t, ek, ak, "%v: not equal at %v", msg, k)
	}
}
