package validator

//import (
//	"context"
//	"testing"
//
//	"github.com/hashicorp/terraform-plugin-framework/attr"
//	"github.com/hashicorp/terraform-plugin-framework/path"
//	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
//	"github.com/hashicorp/terraform-plugin-framework/types"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestNotEmptyStringValidatorValidator(t *testing.T) {
//	t.Parallel()
//
//	type testCase struct {
//		val         attr.Value
//		expectError bool
//	}
//	tests := map[string]testCase{
//		"unknown_string": {
//			val:         types.String{Unknown: true},
//			expectError: false,
//		},
//		"null_string": {
//			val:         types.String{Null: true},
//			expectError: true,
//		},
//		"empty_string": {
//			val:         types.String{Value: ""},
//			expectError: true,
//		},
//		"non-empty_string": {
//			val:         types.String{Value: "test string"},
//			expectError: false,
//		},
//	}
//
//	for name, test := range tests {
//		t.Run(name, func(t *testing.T) {
//			request := tfsdk.ValidateAttributeRequest{
//				AttributePath:           path.Root("test"),
//				AttributePathExpression: path.MatchRoot("test"),
//				AttributeConfig:         test.val,
//			}
//			response := tfsdk.ValidateAttributeResponse{}
//			NotEmptyString().Validate(context.TODO(), request, &response)
//
//			assert.Condition(t, func() (success bool) {
//				return true
//			})
//
//			if !response.Diagnostics.HasError() && test.expectError {
//				t.Fatalf("expected error, got no error")
//			}
//			if response.Diagnostics.HasError() && !test.expectError {
//				t.Fatalf("got unexpected error: %s", response.Diagnostics)
//			}
//		})
//	}
//}
