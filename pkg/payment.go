package numeral

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/xeipuuv/gojsonschema"
)

type Payment struct {
	DebtorIban          string  `json:"debtor_iban"`
	DebtorName          string  `json:"debtor_name"`
	CreditorIban        string  `json:"creditor_iban"`
	CreditorName        string  `json:"creditor_name"`
	Ammount             float32 `json:"ammount"`
	IdempotencyUK       string  `json:"idempotency_unique_key" gorm:"column:idempotency_unique_key"`
	BankProcessedStatus string  `json:"-"`
}

var singletonSchemaLoader gojsonschema.JSONLoader
var once sync.Once

func GetSchemaLoader(schemaPath string) gojsonschema.JSONLoader {
	once.Do(func() {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		if strings.Contains(dir, "/pkg") {
			dir = strings.ReplaceAll(dir, "pkg", "")
		}
		path := fmt.Sprintf("file://%s/%s", dir, schemaPath)
		singletonSchemaLoader = gojsonschema.NewReferenceLoader(path)
	})
	return singletonSchemaLoader
}

func (p *Payment) Validate(schemaPath string) (bool, []gojsonschema.ResultError, error) {
	loader := gojsonschema.NewGoLoader(p)
	schemaLoader := GetSchemaLoader(schemaPath)
	result, err := gojsonschema.Validate(schemaLoader, loader)
	if err != nil {
		return false, nil, err
	}
	if len(result.Errors()) > 0 {
		return false, result.Errors(), nil
	}
	return true, nil, nil
}

func (p *Payment) Insert(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Debug().Create(p).Error; err != nil {
			return err
		}
		return nil
	})
}

func (p *Payment) Update(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Debug().Model(Payment{}).Where("idempotency_unique_key = ?", p.IdempotencyUK).Updates(p).Error; err != nil {
			return err
		}
		return nil
	})
}
