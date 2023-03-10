package validators

import (
	"encoding/json"
	"testing"

	"github.com/asaskevich/govalidator"
)

func TestReadme(t *testing.T) {
	type message struct {
		Statement string `json:"statement" valid:"saysHelloWorld"`
	}

	govalidator.CustomTypeTagMap.Set("saysHelloWorld", NewEquals("Hello World"))

	unmarshalMessage := func(b []byte) (*message, error) {
		var m = &message{}
		_ = json.Unmarshal(b, m)
		_, err := govalidator.ValidateStruct(m)
		if err != nil {
			return nil, err
		}
		return m, nil
	}

	type tc struct {
		m  string
		ok bool
	}

	tests := []tc{
		{
			`
{
	"statement":"Hello World"
}`,
			true,
		},
		{
			`
{
	"statement":"Hello World!!"
}`,
			false,
		},
	}

	for _, test := range tests {
		_, err := unmarshalMessage([]byte(test.m))
		if test.ok && err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !test.ok && err == nil {
			t.Fatal("Expected error, but passed")
		}
	}
}
