package types

import (
	"encoding/json"

	"github.com/go-jstmpl/go-jstmpl/helpers"
	"github.com/go-jstmpl/go-jstmpl/validations"
	"github.com/lestrrat/go-jsschema"
)

type Integer struct {
	*schema.Schema
	NativeType  string `json:"-"`
	ColumnName  string
	ColumnType  string
	Type        string
	Name        string
	key         string
	IsPrivate   bool `json:"-"`
	Validations []validations.Validation
}

func NewInteger(ctx *Context, s *schema.Schema) *Integer {
	vs := []validations.Validation{}
	if v, err := validations.NewMaximumValidation(s); err == nil {
		ctx.AddValidation(v)
		vs = append(vs, v)
	}
	if v, err := validations.NewMinimumValidation(s); err == nil {
		ctx.AddValidation(v)
		vs = append(vs, v)
	}
	var cn, ct string

	if s.Extras["column"] != nil {
		cn, ct, _ = helpers.GetColumnData(s)
	} else {
		cn, ct, _ = helpers.GetColumnData(ctx.Raw)
	}

	return &Integer{
		Schema:      s,
		NativeType:  "number",
		Type:        "number",
		ColumnName:  cn,
		ColumnType:  ct,
		Name:        helpers.SpaceToUpperCamelCase(s.Title),
		key:         ctx.Key,
		IsPrivate:   true,
		Validations: vs,
	}
}

func (o Integer) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Title":       o.Title(),
		"Description": o.Description,
		"Type":        o.Type,
		"Name":        o.Name,
		"Required":    o.Required,
		"Validations": o.Validations,
		"Properties":  o.Properties,
		"ColumnName":  o.ColumnName,
		"ColumnType":  o.ColumnType,
	})
}

func (o Integer) Raw() *schema.Schema {
	return o.Schema
}

func (o Integer) Title() string {
	return o.Schema.Title
}

func (o Integer) Format() string {
	return string(o.Schema.Format)
}

func (o Integer) Key() string {
	return o.key
}

func (o Integer) ReadOnly() bool {
	v := o.Schema.Extras["readOnly"]
	if v == nil {
		return false
	}
	r, ok := v.(bool)
	if !ok {
		return false
	}
	return r
}

func (o Integer) Example(withoutReadOnly bool) interface{} {
	v := o.Schema.Extras["example"]
	if v == nil {
		return 0
	}
	return v
}

func (o Integer) ExampleString() string {
	return helpers.Serialize(o.Example(false))
}
