package numeral

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
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
			ok, _, err := test.Input.Validate("file://../resources/request_schema.json")
			if test.ExpectedSuccess {
				assert.Nil(err)
			}
			assert.Equal(test.ExpectedOutput.ok, ok)
		})
	}
}

// TODO fix
func TestInsert(t *testing.T) {
	assert := assert.New(t)

	var db *gorm.DB
	_, mock, err := sqlmock.NewWithDSN("sqlmock_db_0")
	assert.NoError(err)

	db, err = gorm.Open("sqlmock", "sqlmock_db_0")
	assert.NoError(err)
	defer db.Close()

	type output struct {
		error error
	}

	tests := []struct {
		Description    string
		Input          Payment
		ExpectedOutput output
		ExpectedResult bool
	}{
		{
			Description: "",
			Input:       Payment{},
			ExpectedOutput: output{
				error: errors.New(""),
			},
			ExpectedResult: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {

			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO payments").WithArgs(test.Input).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			err := test.Input.Insert(db)
			if test.ExpectedResult {
				assert.NoError(err)
			} else {
				assert.NotNil(err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
