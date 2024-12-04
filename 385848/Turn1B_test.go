package _85848

import (
	"encoding/json"
	"errors"
	"testing"
)

func UnmarshalPerson(data []byte, person interface{}) error {
	return json.Unmarshal(data, person)
}

func TestUnmarshalPersonB(t *testing.T) {
	testCases := []struct {
		name      string
		inputData []byte
		expected  interface{}
		err       error
	}{
		{
			name:      "Successfully unmarshal valid JSON",
			inputData: []byte(`{"name": "Jenny Doe", "age": 30}`),
			expected: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{
				Name: "Jenny Doe",
				Age:  30,
			},
			err: nil,
		},
		{
			name:      "Error unmarshalling invalid JSON",
			inputData: []byte(`{"invalid": "json"}`),
			expected: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{},
			err: &json.UnmarshalTypeError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}
			err := UnmarshalPerson(tc.inputData, &result)
			if !errors.Is(err, tc.err) {
				t.Errorf("expected error %v, got %v", tc.err, err)
				return
			}

			if err == nil && tc.expected != nil {
				expected, ok := tc.expected.(struct {
					Name string `json:"name"`
					Age  int    `json:"age"`
				})
				if !ok {
					t.Fatalf("expected is not of the correct type")
				}
				if expected.Name != result.Name || expected.Age != result.Age {
					t.Errorf("expected %+v, got %+v", expected, result)
				}
			}
		})
	}
}
