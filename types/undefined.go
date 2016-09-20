package types

import (
	"github.com/go-jstmpl/go-jstmpl/helpers"
	"github.com/go-jstmpl/go-jstmpl/validations"
	"github.com/lestrrat/go-jsschema"
)

type Undefined struct {
	Schema      *schema.Schema `json:"-"`
	NativeType  string         `json:"-"`
	ColumnName  string
	ColumnType  string
	Type        string
	Name        string
	key         string
	IsPrivate   bool `json:"-"`
	Validations []validations.Validation
}

func NewUndefined(ctx *Context, s *schema.Schema) *Undefined {
	vs := []validations.Validation{}
	var cn, ct string

	if s.Extras["column"] != nil {
		cn, ct, _ = helpers.GetColumnData(s)
	} else {
		cn, ct, _ = helpers.GetColumnData(ctx.Raw)
	}

	return &Undefined{
		Schema:      s,
		NativeType:  "undefined",
		Type:        "undefined",
		ColumnName:  cn,
		ColumnType:  ct,
		Name:        helpers.SpaceToUpperCamelCase(s.Title),
		key:         ctx.Key,
		IsPrivate:   true,
		Validations: vs,
	}
}

func (o Undefined) Raw() *schema.Schema {
	return o.Schema
}

func (o Undefined) Title() string {
	return o.Schema.Title
}

func (o Undefined) Description() string {
	return o.Schema.Description
}

func (o Undefined) Key() string {
	return o.key
}

func (o Undefined) ReadOnly() bool {
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

func (o Undefined) Example(withoutReadOnly bool) interface{} {
	v := o.Schema.Extras["example"]
	if v == nil {
		return ""
	}
	return v
}

func (o Undefined) ExampleString() string {
	return helpers.Serialize(o.Example(false))
}
