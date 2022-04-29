package numeral

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type EnvVars struct {
	RequestSchema    string
	BankFolder       string
	SqliteDbLocation string
}

type Handler struct {
	Database Database
	EnvVars  EnvVars
}

func (h *Handler) AddEntity() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		payment, err := GetPaymentFromRequest(req)
		if err != nil {
			responseUnprocessable(w, err.Error())
			return
		}

		ok, customerror, err := payment.Validate(h.EnvVars.RequestSchema)
		if err != nil {
			responseUnprocessable(w, err.Error())
			return
		}

		if !ok {
			errormsg := ""
			for _, e := range customerror {
				errormsg += e.String()
			}
			responseUnprocessable(w, errormsg)
			return
		}

		if err := payment.Insert(h.Database.Db); err != nil {
			log.Println(err)
			responseError(w, fmt.Sprintf("error inserting payment %s", payment.IdempotencyUK))
			return
		}

		if err := GenerateXML(*payment, h.EnvVars.BankFolder); err != nil {
			log.Println(err)
			responseError(w, fmt.Sprintf("error generating XML %s", payment.IdempotencyUK))
			return
		}
		responseOk(w, true, payment.IdempotencyUK)

		time.Sleep(5 * time.Second)
		bankresp, err := GetBankResponse()
		if err != nil {
			log.Println(err)
			responseError(w, fmt.Sprintf("error generating XML %s", payment.IdempotencyUK))
			return
		}

		payment.BankProcessedStatus = bankresp.Status
		if err := payment.Update(h.Database.Db); err != nil {
			log.Println(err)
			responseError(w, fmt.Sprintf("error updating payment %s", payment.IdempotencyUK))
			return
		}
	})
}
