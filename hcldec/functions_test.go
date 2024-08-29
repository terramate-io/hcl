// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hcldec

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/terramate-io/hcl/v2"
	"github.com/terramate-io/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

// This is inspired by hcldec/variables_test.go

func TestFunctions(t *testing.T) {
	tests := []struct {
		config string
		spec   Spec
		want   []hcl.Traversal
	}{
		{
			``,
			&ObjectSpec{},
			nil,
		},
		{
			"a = foo()\n",
			&ObjectSpec{},
			nil, // "a" is not actually used, so "foo" is not required
		},
		{
			"a = foo()\n",
			&AttrSpec{
				Name: "a",
			},
			[]hcl.Traversal{
				{
					hcl.TraverseRoot{
						Name: "foo",
						SrcRange: hcl.Range{
							Start: hcl.Pos{Line: 1, Column: 5, Byte: 4},
							End:   hcl.Pos{Line: 1, Column: 8, Byte: 7},
						},
					},
				},
			},
		},
		{
			"a = foo()\nb = bar()\n",
			&DefaultSpec{
				Primary: &AttrSpec{
					Name: "a",
				},
				Default: &AttrSpec{
					Name: "b",
				},
			},
			[]hcl.Traversal{
				{
					hcl.TraverseRoot{
						Name: "foo",
						SrcRange: hcl.Range{
							Start: hcl.Pos{Line: 1, Column: 5, Byte: 4},
							End:   hcl.Pos{Line: 1, Column: 8, Byte: 7},
						},
					},
				},
				{
					hcl.TraverseRoot{
						Name: "bar",
						SrcRange: hcl.Range{
							Start: hcl.Pos{Line: 2, Column: 5, Byte: 14},
							End:   hcl.Pos{Line: 2, Column: 8, Byte: 17},
						},
					},
				},
			},
		},
		{
			"a = foo()\n",
			&ObjectSpec{
				"a": &AttrSpec{
					Name: "a",
				},
			},
			[]hcl.Traversal{
				{
					hcl.TraverseRoot{
						Name: "foo",
						SrcRange: hcl.Range{
							Start: hcl.Pos{Line: 1, Column: 5, Byte: 4},
							End:   hcl.Pos{Line: 1, Column: 8, Byte: 7},
						},
					},
				},
			},
		},
		{
			`
b {
  a = foo()
}
`,
			&BlockSpec{
				TypeName: "b",
				Nested: &AttrSpec{
					Name: "a",
				},
			},
			[]hcl.Traversal{
				{
					hcl.TraverseRoot{
						Name: "foo",
						SrcRange: hcl.Range{
							Start: hcl.Pos{Line: 3, Column: 7, Byte: 11},
							End:   hcl.Pos{Line: 3, Column: 10, Byte: 14},
						},
					},
				},
			},
		},
		{
			`
b {
  a = foo()
  b = bar()
}
					`,
			&BlockAttrsSpec{
				TypeName:    "b",
				ElementType: cty.String,
			},
			[]hcl.Traversal{
				{
					hcl.TraverseRoot{
						Name: "foo",
						SrcRange: hcl.Range{
							Start: hcl.Pos{Line: 3, Column: 7, Byte: 11},
							End:   hcl.Pos{Line: 3, Column: 10, Byte: 14},
						},
					},
				},
				{
					hcl.TraverseRoot{
						Name: "bar",
						SrcRange: hcl.Range{
							Start: hcl.Pos{Line: 4, Column: 7, Byte: 23},
							End:   hcl.Pos{Line: 4, Column: 10, Byte: 26},
						},
					},
				},
			},
		},
		{
			`
b {
  a = foo()
}
b {
  a = bar()
}
c {
  a = baz()
}
`,
			&BlockListSpec{
				TypeName: "b",
				Nested: &AttrSpec{
					Name: "a",
				},
			},
			[]hcl.Traversal{
				{
					hcl.TraverseRoot{
						Name: "foo",
						SrcRange: hcl.Range{
							Start: hcl.Pos{Line: 3, Column: 7, Byte: 11},
							End:   hcl.Pos{Line: 3, Column: 10, Byte: 14},
						},
					},
				},
				{
					hcl.TraverseRoot{
						Name: "bar",
						SrcRange: hcl.Range{
							Start: hcl.Pos{Line: 6, Column: 7, Byte: 29},
							End:   hcl.Pos{Line: 6, Column: 10, Byte: 32},
						},
					},
				},
			},
		}, /**/
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d-%s", i, test.config), func(t *testing.T) {
			file, diags := hclsyntax.ParseConfig([]byte(test.config), "", hcl.Pos{Line: 1, Column: 1, Byte: 0})
			if len(diags) != 0 {
				t.Errorf("wrong number of diagnostics from ParseConfig %d; want %d", len(diags), 0)
				for _, diag := range diags {
					t.Logf(" - %s", diag.Error())
				}
			}
			body := file.Body

			got := Functions(body, test.spec)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.want)
			}
		})
	}

}
