package numeral

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	assert := assert.New(t)

	type input struct {
		req  *http.Request
		tipe interface{}
	}

	type dummy struct{}

	type output struct {
		Element interface{}
		Error   error
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	}))
	defer func() { testServer.Close() }()

	marshaledBody, err := json.Marshal(Payment{})
	assert.Nil(err)

	tests := []struct {
		Description     string
		Input           input
		ExpectedOutput  output
		ExpectedSuccess bool
	}{
		{
			Description: "when the input type is not an Payment entity, should return a nil element and no error",
			Input: input{
				req:  &http.Request{},
				tipe: dummy{},
			},
			ExpectedOutput: output{
				Element: nil,
				Error:   nil,
			},
			ExpectedSuccess: true,
		},
		{
			Description: "when the input type is an Payment entity, but the body does not represent an payment payload, should return nil element and EOF error",
			Input: input{
				req:  httptest.NewRequest(http.MethodGet, testServer.URL, nil),
				tipe: Payment{},
			},
			ExpectedOutput: output{
				Element: nil,
				Error:   errors.New("EOF"),
			},
			ExpectedSuccess: false,
		},
		{
			Description: "when the input type is an Payment entity, and the body represents a valid payment payload, should return an Payment entity and nil error",
			Input: input{
				req:  httptest.NewRequest(http.MethodPost, testServer.URL, bytes.NewBuffer(marshaledBody)),
				tipe: Payment{},
			},
			ExpectedOutput: output{
				Element: Payment{},
				Error:   nil,
			},
			ExpectedSuccess: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			output, err := Decode(test.Input.req, test.Input.tipe)
			if test.ExpectedSuccess {
				assert.NoError(err)
			} else {
				assert.NotNil(err)
			}
			assert.Equal(test.ExpectedOutput.Element, output)
		})
	}
}

func TestGetPaymentFromRequest(t *testing.T) {
	assert := assert.New(t)

	type input struct {
		req *http.Request
	}

	type output struct {
		Element *Payment
		Error   error
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	}))
	defer func() { testServer.Close() }()

	testpayment := Payment{
		DebtorIban:    "FR1112739000504482744411A64",
		DebtorName:    "company1",
		CreditorIban:  "DE65500105179799248552",
		CreditorName:  "beneficiary",
		Ammount:       42.99,
		IdempotencyUK: "JXJ984XXXZ",
	}

	marshaledBody, err := json.Marshal(testpayment)
	assert.Nil(err)

	tests := []struct {
		Description     string
		Input           input
		ExpectedOutput  output
		ExpectedSuccess bool
	}{
		{
			Description: "when the request body is empty, should return a nil element and an error",
			Input: input{
				req: httptest.NewRequest(http.MethodGet, testServer.URL, nil),
			},
			ExpectedOutput: output{
				Element: nil,
				Error:   errors.New(""),
			},
			ExpectedSuccess: false,
		},
		{
			Description: "when the request body represents a valid payment payload, should return an Payment entity and nil error",
			Input: input{
				req: httptest.NewRequest(http.MethodPost, testServer.URL, bytes.NewBuffer(marshaledBody)),
			},
			ExpectedOutput: output{
				Element: &testpayment,
				Error:   nil,
			},
			ExpectedSuccess: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			output, err := GetPaymentFromRequest(test.Input.req)
			if test.ExpectedSuccess {
				assert.NoError(err)
			} else {
				assert.NotNil(err)
			}
			assert.Equal(test.ExpectedOutput.Element, output)
		})
	}
}

func TestAuthMiddleware(t *testing.T) {
	assert := assert.New(t)

	testHandler := func(w http.ResponseWriter, r *http.Request) {}
	testdomain := "http://yourdomain.gal"

	authKo1 := fmt.Sprintf("Basic " + BasicAuth("usertest", ""))
	authKo2 := fmt.Sprintf("Basic " + BasicAuth("", "userpass"))
	authOk := fmt.Sprintf("Basic " + BasicAuth("usertest", "userpass"))

	testcases := []struct {
		Description    string
		Header         map[string]string
		ExpectedResult int
	}{
		{
			Description: "When required authentication is not present",
			Header: map[string]string{
				"": "",
			},
			ExpectedResult: http.StatusUnauthorized,
		},
		{
			Description: "When required authentication username is present but password is not",
			Header: map[string]string{
				"Authorization": authKo1,
			},
			ExpectedResult: http.StatusUnauthorized,
		},
		{
			Description: "When required authentication password is present but username is not",
			Header: map[string]string{
				"Authorization": authKo2,
			},
			ExpectedResult: http.StatusUnauthorized,
		},
		{
			Description: "When authentication is correct",
			Header: map[string]string{
				"Authorization": authOk,
			},
			ExpectedResult: http.StatusOK,
		},
	}

	for _, test := range testcases {
		log.Println(test.Description)

		req := httptest.NewRequest(http.MethodGet, testdomain, nil)
		for k, v := range test.Header {
			req.Header.Set(k, v)
		}
		res := httptest.NewRecorder()

		testHandler(res, req)
		th := http.HandlerFunc(testHandler)

		zz := AuthMiddleware(th)
		zz.ServeHTTP(res, req)
		assert.Equal(res.Result().StatusCode, test.ExpectedResult)
	}
}
