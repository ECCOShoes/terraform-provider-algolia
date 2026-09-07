package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// expandStringList converts a Terraform list of strings into a Go slice. A null
// or unknown list yields a nil slice.
func expandStringList(ctx context.Context, list types.List) ([]string, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}
	out := make([]string, 0, len(list.Elements()))
	diags := list.ElementsAs(ctx, &out, false)
	return out, diags
}

// flattenStringList converts a Go slice into a Terraform list of strings. A nil
// slice yields a null list so that unset values do not show as empty lists.
func flattenStringList(values []string) types.List {
	if values == nil {
		return types.ListNull(types.StringType)
	}
	elems := make([]types.String, 0, len(values))
	for _, v := range values {
		elems = append(elems, types.StringValue(v))
	}
	list, diags := types.ListValueFrom(context.Background(), types.StringType, elems)
	if diags.HasError() {
		return types.ListNull(types.StringType)
	}
	return list
}

// int32Ptr converts a Terraform int64 into an *int32, returning nil for null or
// unknown values.
func int32Ptr(v types.Int64) *int32 {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	n := int32(v.ValueInt64())
	return &n
}

// int64Value converts an *int32 into a Terraform int64, returning null for nil.
func int64Value(v *int32) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*v))
}

// stringPtr converts a Terraform string into a *string, returning nil for null,
// unknown, or empty values.
func stringPtr(v types.String) *string {
	if v.IsNull() || v.IsUnknown() || v.ValueString() == "" {
		return nil
	}
	s := v.ValueString()
	return &s
}

// stringValue converts a *string into a Terraform string, returning null for nil
// or empty values.
func stringValue(v *string) types.String {
	if v == nil || *v == "" {
		return types.StringNull()
	}
	return types.StringValue(*v)
}

// boolPtr converts a Terraform bool into a *bool, returning nil for null or
// unknown values.
func boolPtr(v types.Bool) *bool {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	b := v.ValueBool()
	return &b
}

// boolValue converts a *bool into a Terraform bool, returning null for nil.
func boolValue(v *bool) types.Bool {
	if v == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*v)
}

// ptr returns a pointer to v. It is used to build SDK request payloads whose
// optional fields are pointers.
func ptr[T any](v T) *T {
	return &v
}
