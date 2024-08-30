// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hcldec

import (
	"encoding/gob"
)

func init() {
	// Every Spec implementation should be registered with gob, so that
	// specs can be sent over gob channels, such as using
	// github.com/hashicorp/go-plugin with plugins that need to describe
	// what shape of configuration they are expecting.
	// NOTE(i4k): use RegisterName() instead of Register() because the latter panics
	// if the registered types are used in alias types.
	gob.RegisterName("hcldec.ObjectSpec", ObjectSpec(nil))
	gob.RegisterName("hcldec.TupleSpec", TupleSpec(nil))
	gob.RegisterName("*hcldec.AttrSpec", (*AttrSpec)(nil))
	gob.RegisterName("*hcldec.LiteralSpec", (*LiteralSpec)(nil))
	gob.RegisterName("*hcldec.ExprSpec", (*ExprSpec)(nil))
	gob.RegisterName("*hcldec.BlockSpec", (*BlockSpec)(nil))
	gob.RegisterName("*hcldec.BlockListSpec", (*BlockListSpec)(nil))
	gob.RegisterName("*hcldec.BlockSetSpec", (*BlockSetSpec)(nil))
	gob.RegisterName("*hcldec.BlockMapSpec", (*BlockMapSpec)(nil))
	gob.RegisterName("*hcldec.BlockLabelSpec", (*BlockLabelSpec)(nil))
	gob.RegisterName("*hcldec.DefaultSpec", (*DefaultSpec)(nil))
}
