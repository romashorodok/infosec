package ent

import (
	"entgo.io/ent/schema"
	"github.com/masseelch/elk"
	"github.com/masseelch/elk/spec"
)

const SecurityBearerAuth string = "BearerAuth"

var securitySpec spec.Security = spec.Security{{SecurityBearerAuth: {}}}

var ElkSecurity schema.Annotation = elk.SchemaSecurity(securitySpec)
