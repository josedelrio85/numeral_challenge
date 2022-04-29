package numeral

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	assert := assert.New(t)

	type output struct {
		ok    bool
		error error
	}

	tests := []struct {
		Description     string
		Input           Payment
		ExpectedOutput  output
		ExpectedSuccess bool
	}{
		{
			Description: "when input payment has no values",
			Input:       Payment{},
			ExpectedOutput: output{
				ok:    false,
				error: errors.New(""),
			},
			ExpectedSuccess: false,
		},
		{
			Description: "when input payment has an invalid debtor_iban property value",
			Input: Payment{
				DebtorIban: "F11112739000504482744411A64",
			},
			ExpectedOutput: output{
				ok:    false,
				error: errors.New(""),
			},
			ExpectedSuccess: false,
		},
		{
			Description: "when input payment matches all validation requirements",
			Input: Payment{
				DebtorIban:    "FR1112739000504482744411A64",
				DebtorName:    "company1",
				CreditorIban:  "DE65500105179799248552",
				CreditorName:  "beneficiary",
				Ammount:       42.99,
				IdempotencyUK: "JXJ984XXXZ",
			},
			ExpectedOutput: output{
				ok:    true,
				error: nil,
			},
			ExpectedSuccess: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			ok, _, err := test.Input.Validate("resources/request_schema.json")
			if test.ExpectedSuccess {
				assert.Nil(err)
			}
			assert.Equal(test.ExpectedOutput.ok, ok)
		})
	}
}
