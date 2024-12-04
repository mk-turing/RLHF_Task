package _85848

import (
	"encoding/json"
	"errors"
	"testing"
)

func unmarshalPerson(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func TestUnmarshalPersonA(t *testing.T) {
	var person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	testCases := []struct {
		name     string
		jsonData []byte
		expected struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		err error
	}{
		{
			name:     "Successfully unmarshals valid JSON",
			jsonData: []byte(`{"name": "Jenny Doe", "age": 30}`),
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
			name:     "Error unmarshals invalid JSON",
			jsonData: []byte(`{"invalid": "json"}`),
			expected: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{},
			err: errors.New("invalid format: \"invalid\": invalid string literal"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := unmarshalPerson(tc.jsonData, &person); err != nil {
				if tc.err == nil {
					t.Errorf("unexpected error: %v", err)
				} else if err.Error() != tc.err.Error() {
					t.Errorf("error message mismatch: expected %v, got %v", tc.err, err)
				}
			} else if tc.err != nil {
				t.Error("expected error, got none")
			}

			if person != tc.expected {
				t.Errorf("unmarshaled result mismatch: expected %+v, got %+v", tc.expected, person)
			}

			// Reset person for the next test case
			person = struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{}
		})
	}
}
