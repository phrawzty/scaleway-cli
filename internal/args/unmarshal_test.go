package args

import (
	"reflect"
	"testing"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalStruct(t *testing.T) {

	type TestCase struct {
		args     []string
		error    string
		expected interface{}
		data     interface{}
	}

	stringPtr := "test"

	run := func(testCase TestCase) func(t *testing.T) {
		return func(t *testing.T) {

			if testCase.data == nil {
				testCase.data = reflect.New(reflect.TypeOf(testCase.expected).Elem()).Interface()
			}
			err := UnmarshalStruct(testCase.args, testCase.data)

			if testCase.error == "" {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, testCase.data)
			} else {
				assert.Equal(t, testCase.error, err.Error())
			}

		}
	}

	RegisterUnmarshalFunc((*height)(nil), unmarshalHeight)

	t.Run("basic", run(TestCase{
		args: []string{
			"string=test",
			"int=42",
			"int16=16",
			"int32=32",
			"int64=64",
			"u-int16=16",
			"u-int32=32",
			"u-int64=64",
			"float32=3.2",
			"float64=6.4",
			"string-ptr=test",
			"bool",
		},
		expected: &Basic{
			String:    "test",
			Int:       42,
			Int16:     16,
			Int32:     32,
			Int64:     64,
			UInt16:    16,
			UInt32:    32,
			UInt64:    64,
			Float32:   3.2,
			Float64:   6.4,
			StringPtr: &stringPtr,
			Bool:      true,
		},
	}))

	t.Run("false-bool", run(TestCase{
		args: []string{
			"bool=false",
		},
		expected: &Basic{
			Bool: false,
		},
	}))

	t.Run("data-must-be-a-pointer", run(TestCase{
		args:  []string{},
		data:  Basic{},
		error: "data must be a pointer to a struct",
	}))

	t.Run("invalid-arg-name", run(TestCase{
		args: []string{
			"testCase=12",
		},
		expected: &Basic{},
		error:    "cannot unmarshal arg 'testCase=12': arg name must only contain lowercase letters, numbers or dashes",
	}))

	t.Run("field-do-not-exist", run(TestCase{
		args: []string{
			"unknown-field=12",
		},
		expected: &Basic{},
		error:    "cannot unmarshal arg 'unknown-field=12': unknown argument",
	}))

	t.Run("invalid-bool", run(TestCase{
		args: []string{
			"bool=invalid",
		},
		expected: &Basic{},
		error:    "cannot unmarshal arg 'bool=invalid': *bool is not unmarshalable: invalid boolean value",
	}))

	t.Run("missing-slice-index", run(TestCase{
		args: []string{
			"strings.1=2",
		},
		expected: &Slice{},
		error:    "cannot unmarshal arg 'strings.1=2': missing index 0, all indices prior to 1 must be set as well",
	}))

	t.Run("missing-slice-indices", run(TestCase{
		args: []string{
			"strings.5=2",
		},
		expected: &Slice{},
		error:    "cannot unmarshal arg 'strings.5=2': missing indices, 0,1,2,3,4 all indices prior to 5 must be set as well",
	}))

	t.Run("missing-slice-indices-overflow", run(TestCase{
		args: []string{
			"strings.99999=2",
		},
		expected: &Slice{},
		error:    "cannot unmarshal arg 'strings.99999=2': missing indices, 0,1,2,3,4,5,6,7,8,9,... all indices prior to 99999 must be set as well",
	}))

	t.Run("duplicate-slice-index", run(TestCase{
		args: []string{
			"basics.0.string=2",
			"basics.0.string=2",
		},
		expected: &Slice{},
		error:    "cannot unmarshal arg 'basics.0.string=2': duplicate argument",
	}))

	t.Run("slice-with-negative-index", run(TestCase{
		args: []string{
			"strings.0=2",
			"strings.-1=2",
		},
		expected: &Slice{},
		error:    "cannot unmarshal arg 'strings.-1=2': invalid index '-1' is not a positive integer",
	}))

	t.Run("nested-slice-with-invalid-index", run(TestCase{
		args: []string{
			"basics.string=test",
		},
		expected: &Slice{},
		error:    "cannot unmarshal arg 'basics.string=test': invalid index 'string' is not a positive integer",
	}))

	t.Run("basic-slice", run(TestCase{
		args: []string{
			"strings.0=1",
			"strings.1=2",
			"strings.2=3",
			"strings.3=test",
			"strings-ptr.0=test",
			"strings-ptr.1=test",
			"basics.0.string=test",
			"basics.0.int=42",
			"basics.1.string=test",
			"basics.1.int=42",
		},
		expected: &Slice{
			Strings:    []string{"1", "2", "3", "test"},
			StringsPtr: []*string{&stringPtr, &stringPtr},
			Basics: []Basic{
				{
					String: "test",
					Int:    42,
				},
				{
					String: "test",
					Int:    42,
				},
			},
		},
	}))

	t.Run("well-known-types", run(TestCase{
		args: []string{
			"size=20gb",
		},
		expected: &WellKnownTypes{
			Size: 20 * scw.GB,
		},
	}))

	t.Run("nested-basic", run(TestCase{
		args: []string{
			"basic.string=test",
		},
		expected: &Nested{
			Basic: Basic{
				String: "test",
			},
		},
	}))

	t.Run("map-basic", run(TestCase{
		args: []string{
			"map.key1=value1",
			"map.key2=value2",
		},
		expected: &Map{
			Map: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}))

	t.Run("custom", run(TestCase{
		args: []string{
			"custom-struct=test",
			"custom-string=test",
		},
		expected: &CustomWrapper{
			CustomStruct: &CustomStruct{
				value: "TEST",
			},
			CustomString: CustomString("TEST"),
		},
	}))

	t.Run("insane", run(TestCase{
		args: []string{
			"map.key1.key2.basic.string=test",
		},
		expected: func() interface{} {
			n1 := &Nested{Basic: Basic{String: "test"}}
			m1 := &map[string]**Nested{"key2": &n1}
			m2 := map[string]**map[string]**Nested{"key1": &m1}
			return &Insane{Map: &m2}
		}(),
	}))

	t.Run("data-is-a-map", run(TestCase{
		args: []string{
			"key1=v1",
			"key2=v2",
		},
		expected: &map[string]string{
			"key1": "v1",
			"key2": "v2",
		},
		data: &map[string]string{},
	}))

	t.Run("data-is-an-enum", run(TestCase{
		args: []string{
			"color=blue",
			"size=1",
		},
		expected: &Enum{Color: ColorBlue, Size: Size1},
	}))

	t.Run("data-is-raw-args", run(TestCase{
		args: []string{
			"pro.access_key",
			"access_key",
		},
		expected: &RawArgs{
			"pro.access_key",
			"access_key",
		},
	}))

	h := height(14)
	t.Run("height-set", run(TestCase{
		args: []string{
			"height=14cm",
		},
		data: &CustomArgs{},
		expected: &CustomArgs{
			Height: &h,
		},
	}))

	t.Run("height-not-set", run(TestCase{
		args:     []string{},
		data:     &CustomArgs{},
		expected: &CustomArgs{},
	}))

	t.Run("duplicate-keys-simple", run(TestCase{
		args: []string{
			"custom-struct=test",
			"custom-struct=test2",
			"custom-string=test",
		},
		data:  &CustomWrapper{},
		error: "cannot unmarshal arg 'custom-struct=test2': duplicate argument",
	}))

	t.Run("duplicate-keys-insane", run(TestCase{
		args: []string{
			"map.key1.key2.basic.string=test",
			"map.key1.key2.basic.string=test2",
		},
		data:  &Insane{},
		error: "cannot unmarshal arg 'map.key1.key2.basic.string=test2': duplicate argument",
	}))

	t.Run("anonymous-nested-field", run(TestCase{
		args: []string{
			"all=all",
			"merge1=1",
			"merge2=2",
			"merge-only=2",
		},
		expected: &Merge{
			Merge1: Merge1{
				All:       "",
				Merge1:    "1",
				MergeOnly: "",
			},
			Merge2: &Merge2{
				All:       "",
				Merge2:    "2",
				MergeOnly: "2",
			},
			All: "all",
		},
	}))
}

func TestIsUmarshalableValue(t *testing.T) {

	type TestCase struct {
		expected bool
		data     interface{}
	}

	run := func(testCase TestCase) func(t *testing.T) {
		return func(t *testing.T) {

			value := IsUmarshalableValue(testCase.data)
			assert.Equal(t, testCase.expected, value)
		}
	}

	RegisterUnmarshalFunc((*height)(nil), unmarshalHeight)

	strPtr := "This is a pointer"
	heightPtr := height(42)
	customStringPtr := CustomString("test")

	t.Run("string", run(TestCase{
		data:     "a simple string",
		expected: true,
	}))

	t.Run("int", run(TestCase{
		data:     42,
		expected: true,
	}))

	t.Run("custom", run(TestCase{
		data:     CustomString("CUSTOM-STRING"),
		expected: true,
	}))

	t.Run("nil", run(TestCase{
		data:     nil,
		expected: false,
	}))

	t.Run("custom-func", run(TestCase{
		data:     height(42),
		expected: true,
	}))

	t.Run("a-struct", run(TestCase{
		data:     &Basic{},
		expected: false,
	}))

	t.Run("str-pointer", run(TestCase{
		data:     &strPtr,
		expected: true,
	}))
	t.Run("custom-func-pointer", run(TestCase{
		data:     &heightPtr,
		expected: true,
	}))
	t.Run("custom-pointer", run(TestCase{
		data:     &customStringPtr,
		expected: true,
	}))
	t.Run("custom-pointer", run(TestCase{
		data:     map[string]string{},
		expected: false,
	}))

	t.Run("custom-pointer", run(TestCase{
		data:     []string{},
		expected: false,
	}))
}
