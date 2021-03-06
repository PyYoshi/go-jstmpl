package types

import (
	"sort"

	"github.com/go-jstmpl/go-jstmpl/validations"
	"github.com/lestrrat/go-jsschema"
)

type Context struct {
	Key         string
	Raw         *schema.Schema
	Validations map[string]bool
}

func (ctx *Context) AddValidation(v validations.Validation) {
	ctx.Validations[v.Func()] = true
}

func (ctx *Context) RequiredValidations() []string {
	vs := []string{}
	for v, b := range ctx.Validations {
		if !b {
			continue
		}
		vs = append(vs, v)
	}
	sort.Strings(vs)
	return vs
}
