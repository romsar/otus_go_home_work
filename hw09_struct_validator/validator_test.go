package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	//nolint:unused,structcheck
	User struct {
		ID        string `json:"id" validate:"len:36"`
		Name      string
		Age       int      `validate:"min:18|max:50"`
		BloodType int      `validate:"in:1,2,3,4"` // группа крови
		Email     string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role      UserRole `validate:"in:admin,stuff"`
		Phones    []string `validate:"len:11"`
		Company   string   `json:"omitempty"`
		meta      json.RawMessage
	}

	Dog struct {
		Name string `validate:"len:foo|foo:bar"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name           string
		in             interface{}
		expectedErr    error
		expectedErrMsg string
	}{
		{
			name: "no validation errors",
			in: User{
				ID:        "05588232-205c-11ed-861d-0242ac120002",
				Age:       20,
				BloodType: 1,
				Email:     "email@mail.ru",
				Role:      "admin",
				Phones:    []string{"12345678901"},
			},
			expectedErr:    nil,
			expectedErrMsg: "",
		},
		//nolint:lll
		{
			name: "invalid values #1",
			in: User{
				ID:        "foo bar",
				Age:       55,
				BloodType: 6,
				Email:     "email",
				Role:      "foo",
				Phones:    []string{"1234567890111"},
			},
			expectedErr:    ValidationErrors{},
			expectedErrMsg: "invalid value: value length (7) is not match required length (36), invalid value: value (55) is more than (50), invalid value: value (6) is not in (1,2,3,4), invalid value: value (email) is not matched regexp (^\\w+@\\w+\\.\\w+$), invalid value: value (foo) is not in (admin,stuff), invalid value: value length (13) is not match required length (11)",
		},
		{
			name: "invalid values #2",
			in: User{
				ID:        "05588232-205c-11ed-861d-0242ac120002",
				Age:       1,
				BloodType: 0,
				Email:     "email@mail.ru",
				Role:      "admin",
				Phones:    []string{},
			},
			expectedErr:    ValidationErrors{},
			expectedErrMsg: "invalid value: value (1) is less than (18), invalid value: value (0) is not in (1,2,3,4)",
		},
		{
			name: "invalid rules",
			in: Dog{
				Name: "baron",
			},
			expectedErr:    ErrInvalidValidationRule,
			expectedErrMsg: "invalid validation rule",
		},
		{
			name:           "invalid entity",
			in:             "foobar",
			expectedErr:    ErrNotValidatable,
			expectedErrMsg: "that entity cannot be validated",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
				return
			}

			require.Error(t, err, tt.expectedErr)
			if err != nil {
				require.EqualError(t, err, tt.expectedErrMsg)
			}
		})
	}
}
