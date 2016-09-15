package types

import (
	"github.com/go-jstmpl/go-jstmpl/helpers"
	"github.com/go-jstmpl/go-jstmpl/validations"
	"github.com/lestrrat/go-jsschema"
)

type Object struct {
	Schema      *schema.Schema `json:"-"`
	NativeType  string         `json:"-"`
	TableName   string
	ColumnName  string
	ColumnType  string
	Type        string
	Name        string
	key         string
	IsPrivate   bool `json:"-"`
	Required    []string
	Validations []validations.Validation
	Properties  []Schema
}

func NewObject(ctx *Context, s *schema.Schema) *Object {
	var tn string
	if s.Extras["table"] != nil {
		tn, _ = helpers.GetTableData(s)
	} else {
		tn, _ = helpers.GetTableData(ctx.Raw)
	}
	return &Object{
		Schema:     s,
		NativeType: "object",
		TableName:  tn,
		Type:       helpers.SpaceToUpperCamelCase(s.Title),
		Name:       helpers.SpaceToUpperCamelCase(s.Title),
		key:        ctx.Key,
		Required:   ctx.Raw.Required,
		IsPrivate:  false,
		Properties: []Schema{},
	}
}

func (o Object) Raw() *schema.Schema {
	return o.Schema
}

func (o Object) Title() string {
	return o.Schema.Title
}

func (o Object) Format() string {
	return string(o.Schema.Format)
}

func (o Object) Key() string {
	return o.key
}

func (o Object) Example() interface{} {
	e := make(map[string]interface{})
	for _, s := range o.Properties {
		e[s.Key()] = s.Example()
	}
	return e
}
