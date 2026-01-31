package validator

//import (
//	"context"
//	"testing"
//
//	"github.com/hashicorp/terraform-plugin-framework/diag"
//	"github.com/hashicorp/terraform-plugin-framework/path"
//	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
//	"github.com/hashicorp/terraform-plugin-framework/types"
//	"github.com/hashicorp/terraform-plugin-go/tftypes"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestRequiredConfiguredCredentialsValidator(t *testing.T) {
//	t.Parallel()
//
//	type testCase struct {
//		req      tfsdk.ValidateAttributeRequest
//		expected *tfsdk.ValidateAttributeResponse
//	}
//	tests := map[string]testCase{
//		"non-type_attribute": {
//			req: tfsdk.ValidateAttributeRequest{
//				AttributeConfig: types.String{Value: "test-value"},
//				Config: tfsdk.Config{
//					Schema: tfsdk.Schema{
//						Attributes: map[string]tfsdk.Attribute{
//							"test": {Optional: true, Type: types.StringType},
//						},
//					},
//					Raw: tftypes.NewValue(
//						tftypes.Object{
//							AttributeTypes: map[string]tftypes.Type{
//								"test": tftypes.String,
//							},
//						},
//						map[string]tftypes.Value{
//							"test": tftypes.NewValue(tftypes.String, "test-value"),
//						},
//					),
//				},
//			},
//			expected: &tfsdk.ValidateAttributeResponse{
//				Diagnostics: diag.Diagnostics{
//					diag.NewAttributeErrorDiagnostic(
//						path.Path{},
//						"Ensure that integration is correctly configured",
//						"Validator must be applied for integration 'type' only.",
//					),
//				},
//			},
//		},
//		"type_attribute": {
//			req: tfsdk.ValidateAttributeRequest{
//				AttributeConfig:         types.String{Value: "type-value"},
//				AttributePath:           path.Root("type"),
//				AttributePathExpression: path.MatchRoot("type"),
//				Config: tfsdk.Config{
//					Schema: tfsdk.Schema{
//						Attributes: map[string]tfsdk.Attribute{
//							"type": {Optional: true, Type: types.StringType},
//						},
//					},
//					Raw: tftypes.NewValue(
//						tftypes.Object{
//							AttributeTypes: map[string]tftypes.Type{
//								"type": tftypes.String,
//							},
//						},
//						map[string]tftypes.Value{
//							"type": tftypes.NewValue(tftypes.String, "type-value"),
//						},
//					),
//				},
//			},
//			expected: &tfsdk.ValidateAttributeResponse{},
//		},
//	}
//
//	for name, test := range tests {
//		t.Run(name, func(t *testing.T) {
//			actual := tfsdk.ValidateAttributeResponse{}
//
//			RequiresConfiguredCredentials().Validate(context.TODO(), test.req, &actual)
//
//			assert.Equal(t, test.expected.Diagnostics, actual.Diagnostics)
//		})
//	}
//}
//
//func TestRequiredConfiguredCredentialsValidator_github(t *testing.T) {
//	t.Parallel()
//
//	type testCase struct {
//		req      tfsdk.ValidateAttributeRequest
//		expected *tfsdk.ValidateAttributeResponse
//	}
//	tests := map[string]testCase{
//		"with-token": {
//			req: tfsdk.ValidateAttributeRequest{
//				AttributeConfig:         types.String{Value: "github"},
//				AttributePath:           path.Root("type"),
//				AttributePathExpression: path.MatchRoot("type"),
//				Config: tfsdk.Config{
//					Schema: tfsdk.Schema{
//						Attributes: map[string]tfsdk.Attribute{
//							"type":  {Optional: true, Type: types.StringType},
//							"token": {Optional: true, Type: types.StringType},
//						},
//					},
//					Raw: tftypes.NewValue(
//						tftypes.Object{
//							AttributeTypes: map[string]tftypes.Type{
//								"type":  tftypes.String,
//								"token": tftypes.String,
//							},
//						},
//						map[string]tftypes.Value{
//							"type":  tftypes.NewValue(tftypes.String, "github"),
//							"token": tftypes.NewValue(tftypes.String, "github-token"),
//						},
//					),
//				},
//			},
//			expected: &tfsdk.ValidateAttributeResponse{},
//		},
//		"without-token": {
//			req: tfsdk.ValidateAttributeRequest{
//				AttributeConfig:         types.String{Value: "github"},
//				AttributePath:           path.Root("type"),
//				AttributePathExpression: path.MatchRoot("type"),
//				Config: tfsdk.Config{
//					Schema: tfsdk.Schema{
//						Attributes: map[string]tfsdk.Attribute{
//							"type":  {Optional: true, Type: types.StringType},
//							"token": {Optional: true, Type: types.StringType},
//						},
//					},
//					Raw: tftypes.NewValue(
//						tftypes.Object{
//							AttributeTypes: map[string]tftypes.Type{
//								"type":  tftypes.String,
//								"token": tftypes.String,
//							},
//						},
//						map[string]tftypes.Value{
//							"type":  tftypes.NewValue(tftypes.String, "github"),
//							"token": tftypes.NewValue(tftypes.String, ""),
//						},
//					),
//				},
//			},
//			expected: &tfsdk.ValidateAttributeResponse{
//				Diagnostics: diag.Diagnostics{
//					diag.NewAttributeErrorDiagnostic(
//						path.Root("type"),
//						"Ensure that integration is correctly configured",
//						"token must be defined and not empty for 'github' integration",
//					),
//				},
//			},
//		},
//	}
//
//	for name, test := range tests {
//		t.Run(name, func(t *testing.T) {
//			actual := tfsdk.ValidateAttributeResponse{}
//
//			RequiresConfiguredCredentials().Validate(context.TODO(), test.req, &actual)
//
//			assert.Equal(t, test.expected.Diagnostics, actual.Diagnostics)
//		})
//	}
//}
