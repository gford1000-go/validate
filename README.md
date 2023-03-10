[![Go Reference](https://pkg.go.dev/badge/github.com/gford1000-go/validate.svg)](https://pkg.go.dev/github.com/gford1000-go/validate)

validate
========

validate provides convenience functions to create `CustomTypeValidator` instances for use within the [govalidator](https://github.com/asaskevich/govalidator) `CustomTypeTagMap`.

These can be used together with `govalidator` to ensure JSON objects are valid.

```json
{
	"statement": "Hello World"
}
```

## Use

The "valid" tag on a struct attribute defines the list of custom validators (or govalidator defined validators) which are then
applied sequentially to the attribute value.

```go
type Message struct {
	Statement string `json:"statement" valid:"saysHelloWorld"`
}

func init() {
	govalidator.CustomTypeTagMap.Set("saysHelloWorld", NewEquals("Hello World"))
}

func UnmarshalMessage(b []byte)  (*Message, error) {
	var m = &Message{}
	_ = json.Unmarshal(b, m)
	_, err := govalidator.ValidateStruct(m)
	if err != nil {
		return nil, err
	}
	return m
}
```

## How?

The command line is all you need.

```
go get github.com/gford1000-go/validate
```
