package provider

import (
	"testing"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExpandFlattenTypoTolerance(t *testing.T) {
	cases := []struct {
		name  string
		value string
	}{
		{"true", "true"},
		{"false", "false"},
		{"min", "min"},
		{"strict", "strict"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expanded := expandTypoTolerance(types.StringValue(tc.value))
			if expanded == nil {
				t.Fatalf("expandTypoTolerance(%q) returned nil", tc.value)
			}
			got := flattenTypoTolerance(expanded)
			if got.ValueString() != tc.value {
				t.Fatalf("roundtrip mismatch: got %q, want %q", got.ValueString(), tc.value)
			}
		})
	}
}

func TestExpandTypoToleranceNull(t *testing.T) {
	if got := expandTypoTolerance(types.StringNull()); got != nil {
		t.Fatalf("expected nil for null typo_tolerance, got %#v", got)
	}
	if got := flattenTypoTolerance(nil); !got.IsNull() {
		t.Fatalf("expected null for nil typo tolerance, got %q", got.ValueString())
	}
}

func TestFlattenStringListNilIsNull(t *testing.T) {
	if got := flattenStringList(nil); !got.IsNull() {
		t.Fatalf("expected null list for nil slice, got %#v", got)
	}
	got := flattenStringList([]string{"a", "b"})
	if got.IsNull() {
		t.Fatal("expected non-null list for populated slice")
	}
	if len(got.Elements()) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(got.Elements()))
	}
}

func TestACLRoundtrip(t *testing.T) {
	in := []string{"search", "browse", "addObject"}
	acls := stringsToACL(in)
	out := aclToStrings(acls)
	if len(out) != len(in) {
		t.Fatalf("length mismatch: got %d, want %d", len(out), len(in))
	}
	for i := range in {
		if out[i] != in[i] {
			t.Fatalf("element %d mismatch: got %q, want %q", i, out[i], in[i])
		}
	}
	if _, ok := any(acls[0]).(search.Acl); !ok {
		t.Fatal("expected search.Acl elements")
	}
}

func TestInt32PtrConversions(t *testing.T) {
	if got := int32Ptr(types.Int64Null()); got != nil {
		t.Fatalf("expected nil for null int64, got %d", *got)
	}
	p := int32Ptr(types.Int64Value(42))
	if p == nil || *p != 42 {
		t.Fatalf("expected 42, got %#v", p)
	}
	if got := int64Value(nil); !got.IsNull() {
		t.Fatal("expected null int64 for nil pointer")
	}
	if got := int64Value(p); got.ValueInt64() != 42 {
		t.Fatalf("expected 42, got %d", got.ValueInt64())
	}
}
