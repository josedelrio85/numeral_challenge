package numeral

import (
	"fmt"
	"log"
	"net/http"
)

type Handler struct {
	Database Database
}

func (h *Handler) AddEntity() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		payment, err := GetPaymentFromRequest(req)
		if err != nil {
			responseUnprocessable(w, err.Error())
			return
		}

		ok, _, err := payment.Validate()
		if err != nil {
			responseUnprocessable(w, err.Error())
			return
		}

		if !ok {
			log.Println("no ok!") // TODO
			responseUnprocessable(w, err.Error())
			return
		}

		if err := payment.Insert(h.Database.Db); err != nil {
			log.Println(err)
			responseError(w, fmt.Sprintf("error inserting payment %s", payment.IdempotencyUK))
			return
		}

		// TODO xml
		responseOk(w, true, "")
	})
}
