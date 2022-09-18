package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestTypeConversion_fromBool(t *testing.T) {
	t.Parallel()

	valTrue := true
	valFalse := false

	type testCase struct {
		val      *bool
		expected types.Bool
	}
	tests := map[string]testCase{
		"nil": {
			val:      nil,
			expected: types.Bool{Null: true},
		},
		"true": {
			val:      &valTrue,
			expected: types.Bool{Value: valTrue},
		},
		"false": {
			val:      &valFalse,
			expected: types.Bool{Value: valFalse},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := fromBoolPtr(test.val)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestTypeConversion_toBool(t *testing.T) {
	t.Parallel()

	expectedTrue := true
	expectedFalse := false

	type testCase struct {
		val      types.Bool
		expected *bool
	}
	tests := map[string]testCase{
		"null_bool": {
			val:      types.Bool{Null: true},
			expected: nil,
		},
		"unknown_bool": {
			val:      types.Bool{Unknown: true},
			expected: nil,
		},
		"true_bool": {
			val:      types.Bool{Value: true},
			expected: &expectedTrue,
		},
		"false_bool": {
			val:      types.Bool{Value: false},
			expected: &expectedFalse,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := toBoolPtr(test.val)
			assert.Equal(t, test.expected, actual)
		})
	}
}
