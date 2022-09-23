package validator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ tfsdk.AttributeValidator = notEmptyStringValidator{}

// notEmptyStringValidator validates that a string Attribute's content is not empty.
type notEmptyStringValidator struct{}

func (validator notEmptyStringValidator) Description(_ context.Context) string {
	return "string must not be empty"
}

func (validator notEmptyStringValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (validator notEmptyStringValidator) Validate(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) {
	s := request.AttributeConfig.(types.String)

	if s.Unknown {
		return
	}

	if s.Value == "" {
		response.Diagnostics.AddAttributeError(
			request.AttributePath,
			validator.Description(ctx),
			"",
		)
	}
}

// NotEmptyString checks that the string is not empty.
func NotEmptyString() tfsdk.AttributeValidator {
	return &notEmptyStringValidator{}
}
