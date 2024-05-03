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
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	EmailStruct struct {
		Email string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	}

	App struct {
		Version string `validate:"len:5"`
	}

	AppArr struct {
		Version []string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	RoleStruct struct {
		Role UserRole `validate:"in:admin,stuff"`
	}

	MinAgeStruct struct {
		Age int `validate:"min:18"`
	}

	MaxAgeStruct struct {
		Age int `validate:"max:20"`
	}

	AgeStruct struct {
		Age int `validate:"min:18|max:20"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "219bbdc3-78f5-4cd8-86a7-e1a8c373a571",
				Name:   "John",
				Age:    20,
				Email:  "john@example.com",
				Role:   "stuff",
				Phones: []string{"12345678900"},
			},
			expectedErr: nil,
		},
		// failed validation
		{
			in: App{Version: "aaaaaa"},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrValidateLen,
				},
			},
		},
		{
			in: User{
				ID:     "219bbdc3-78f5-4cd8-86a7",
				Name:   "John",
				Age:    17,
				Email:  "johnexample.com",
				Role:   "stufff",
				Phones: []string{"123"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrValidateLen,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrValidateMin,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrValidateRegexp,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrValidateIn,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrValidateLen,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedErr.Error())
			}
			_ = tt
		})
	}
}

func TestValidateLen(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          AppArr{Version: []string{"aaaaa", "bbbbb"}},
			expectedErr: nil,
		},
		{
			in:          App{Version: "aaaaa"},
			expectedErr: nil,
		},
		// failed validation
		{
			in: App{Version: "aaaaaa"},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrValidateLen,
				},
			},
		},
		{
			in: AppArr{Version: []string{"aaaaa", "bbbbbb"}},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrValidateLen,
				},
			},
		},
		{
			in: struct {
				Version []int `validate:"len:5"`
			}{Version: []int{123, 123}},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrValidate,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedErr.Error())
			}
			_ = tt
		})
	}
}

func TestValidateRegexp(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          EmailStruct{Email: "asd@asd.com"},
			expectedErr: nil,
		},
		{
			in: struct {
				Emails []string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
			}{Emails: []string{"asd@asd.com", "bsd@bsd.com"}},
			expectedErr: nil,
		},

		// failed validation
		{
			in: EmailStruct{Email: "asd@asd"},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Email",
					Err:   ErrValidateRegexp,
				},
			},
		},
		{
			in: EmailStruct{Email: ""},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Email",
					Err:   ErrValidateRegexp,
				},
			},
		},
		{
			in: struct {
				Emails []string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
			}{Emails: []string{"asd@asd.com", "bsd@bsd"}},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Emails",
					Err:   ErrValidateRegexp,
				},
			},
		},
		// wrong type
		{
			in: 123,
			expectedErr: ValidationErrors{
				ValidationError{
					Err: ErrValidate,
				},
			},
		},
		{
			in: struct {
				Email int `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
			}{
				Email: 123,
			},
			expectedErr: ValidationErrors{ValidationError{
				Field: "Email",
				Err:   ErrValidate,
			}},
		},
		{
			in: struct {
				Emails []int `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
			}{
				Emails: []int{123},
			},
			expectedErr: ValidationErrors{ValidationError{
				Field: "Emails",
				Err:   ErrValidate,
			}},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedErr.Error())
			}
			_ = tt
		})
	}
}

func TestValidateIn(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: Response{
				Code: 200,
			},
			expectedErr: nil,
		},
		{
			in: RoleStruct{
				Role: "admin",
			},
			expectedErr: nil,
		},
		{
			in: struct {
				Codes []int `validate:"in:200,404,500"`
			}{
				Codes: []int{200, 404},
			},
			expectedErr: nil,
		},
		{
			in: struct {
				Roles []string `validate:"in:admin,stuff"`
			}{
				Roles: []string{"admin", "stuff"},
			},
			expectedErr: nil,
		},
		// failed validation
		{
			in: RoleStruct{
				Role: "super admin",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Role",
					Err:   ErrValidateIn,
				},
			},
		},
		{
			in: Response{
				Code: 201,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrValidateIn,
				},
			},
		},
		{
			in: struct {
				Codes []int `validate:"in:200,404,500"`
			}{
				Codes: []int{201, 404},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Codes",
					Err:   ErrValidateIn,
				},
			},
		},
		// wrong type
		{
			in:          123,
			expectedErr: ValidationErrors{ValidationError{Err: ErrValidate}},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedErr.Error())
			}
			_ = tt
		})
	}
}

func TestValidateMin(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: MinAgeStruct{
				Age: 20,
			},
			expectedErr: nil,
		},
		{
			in: MinAgeStruct{
				Age: 18,
			},
			expectedErr: nil,
		},
		{
			in: struct {
				Ages []int `validate:"min:18"`
			}{
				Ages: []int{18, 20},
			},
			expectedErr: nil,
		},
		// failed validation
		{
			in: MinAgeStruct{
				Age: 17,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   ErrValidateMin,
				},
			},
		},
		{
			in: struct {
				Ages []int `validate:"min:18"`
			}{
				Ages: []int{17, 20},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Ages",
					Err:   ErrValidateMin,
				},
			},
		},
		// wrong type
		{
			in:          123,
			expectedErr: ValidationErrors{ValidationError{Err: ErrValidate}},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedErr.Error())
			}
			_ = tt
		})
	}
}

func TestValidateMax(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: MaxAgeStruct{
				Age: 20,
			},
			expectedErr: nil,
		},
		{
			in: MaxAgeStruct{
				Age: 18,
			},
			expectedErr: nil,
		},
		{
			in: struct {
				Ages []int `validate:"max:21"`
			}{
				Ages: []int{18, 20},
			},
			expectedErr: nil,
		},
		// failed validation
		{
			in: MaxAgeStruct{
				Age: 21,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   ErrValidateMax,
				},
			},
		},
		{
			in: struct {
				Ages []int `validate:"max:17"`
			}{
				Ages: []int{18, 17},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Ages",
					Err:   ErrValidateMax,
				},
			},
		},
		// wrong type
		{
			in:          123,
			expectedErr: ValidationErrors{ValidationError{Err: ErrValidate}},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedErr.Error())
			}
			_ = tt
		})
	}
}

func TestValidateMinMax(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: AgeStruct{
				Age: 20,
			},
			expectedErr: nil,
		},
		{
			in: AgeStruct{
				Age: 18,
			},
			expectedErr: nil,
		},
		{
			in: struct {
				Ages []int `validate:"min:18|max:20"`
			}{
				Ages: []int{18, 20},
			},
			expectedErr: nil,
		},

		// failed validation
		{
			in: AgeStruct{
				Age: 17,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   ErrValidateMin,
				},
			},
		},
		{
			in: AgeStruct{
				Age: 21,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   ErrValidateMax,
				},
			},
		},
		{
			in: struct {
				Ages []int `validate:"min:18|max:20"`
			}{
				Ages: []int{17, 20},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Ages",
					Err:   ErrValidateMin,
				},
			},
		},
		// // wrong type
		{
			in:          123,
			expectedErr: ValidationErrors{ValidationError{Err: ErrValidate}},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedErr.Error())
			}
			_ = tt
		})
	}
}

func TestValidateRegexpLen(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: struct {
				Email string `validate:"regexp:^\\w+@\\w+\\.\\w+$|len:11"`
			}{
				Email: "asd@asd.com",
			},
			expectedErr: nil,
		},
		// failed validation
		{
			in: struct {
				Email string `validate:"regexp:^\\w+@\\w+\\.\\w+$|len:11"`
			}{
				Email: "asdasd.com",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Email",
					Err:   ErrValidateRegexp,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrValidateLen,
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedErr.Error())
			}
			_ = tt
		})
	}
}
