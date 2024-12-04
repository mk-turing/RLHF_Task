package _85848

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestUnmarshalPerson2A(t *testing.T) {
	testCases := []struct {
		name      string
		inputData []byte
		expected  interface{}
		err       error
	}{
		{
			name:      "Successfully unmarshal valid JSON",
			inputData: []byte(`{"name": "John Doe", "age": 30, "city": "New York"}`),
			expected: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
				City string `json:"city"`
			}{
				Name: "John Doe",
				Age:  30,
				City: "New York",
			},
			err: nil,
		},
		{
			name:      "Successfully unmarshal with missing fields",
			inputData: []byte(`{"name": "Jane Smith"}`),
			expected: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
				City string `json:"city"`
			}{
				Name: "Jane Smith",
				Age:  0,
				City: "",
			},
			err: nil,
		},
		{
			name:      "Successfully unmarshal with extra fields",
			inputData: []byte(`{"name": "Alice Brown", "age": 25, "city": "Paris", "favoriteFood": "Pizza"}`),
			expected: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
				City string `json:"city"`
			}{
				Name: "Alice Brown",
				Age:  25,
				City: "Paris",
			},
			err: nil,
		},
		{
			name:      "Error unmarshalling invalid JSON",
			inputData: []byte(`{"invalid": "json"}`),
			expected: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
				City string `json:"city"`
			}{},
			err: &json.UnmarshalTypeError{},
		},
		{
			name:      "Error with incompatible field type",
			inputData: []byte(`{"name": "Emma", "age": "thirty"}`),
			expected: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{},
			err: &json.UnmarshalTypeError{},
		},
		{
			name:      "Missing required field",
			inputData: []byte(`{"age": 35}`),
			expected: struct {
				Name string `json:"name" required:"true"`
				Age  int    `json:"age"`
			}{},
			err: &json.UnmarshalTypeError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result interface{}
			err := UnmarshalPerson(tc.inputData, &result)
			if !errors.Is(err, tc.err) {
				t.Errorf("expected error %v, got %v", tc.err, err)
				return
			}

			if err == nil {
				// Assert equality based on the concrete type of the expected result
				assertStructEqual(t, tc.expected, result)
			}
		})
	}
}

func assertStructEqual(t *testing.T, expected interface{}, actual interface{}) {
	switch exp := expected.(type) {
	case struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
		City string `json:"city"`
	}:
		act, ok := actual.(struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
			City string `json:"city"`
		})
		if !ok {
			t.Fatalf("expected is not of the correct type")
		}
		if exp.Name != act.Name || exp.Age != act.Age || exp.City != act.City {
			t.Errorf("expected %+v, got %+v", exp, act)
		}
	case struct {
		Name string `json:"name" required:"true"`
		Age  int    `json:"age"`
	}:
		act, ok := actual.(struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		})
		if !ok {
			t.Fatalf("expected is not of the correct type")
		}
		if exp.Name != act.Name || exp.Age != act.Age {
			t.Errorf("expected %+v, got %+v", exp, act)
		}
	default:
		t.Fatalf("unhandled expected type: %T", expected)
	}
}
