package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"form3/api"
	"form3/business"
	"form3/business/memory"

	logrusTest "github.com/sirupsen/logrus/hooks/test"
	"gopkg.in/mgo.v2/bson"
)

// todo this ensures some basic protection against regressions, but it needs more edge cases
func TestCRUD(t *testing.T) {
	// implementation of model.Storage
	col := memory.Storage{}

	logger, _ := logrusTest.NewNullLogger()
	service := api.NewService(col, *logger)

	// some of these tests can also work without opening this local listener
	// by testing individual handlers rather than this server instance
	ts := httptest.NewServer(service.GetMux())
	defer ts.Close()

	// this will be used to keep some state between the tests below
	var payment business.Payment

	// isolating tests into closures with deps clearly outlined.
	t.Run("create", create(ts, &payment))
	t.Run("list", list(ts, payment))
	t.Run("retrieve", retrieve(ts, payment.Id.Hex()))
	t.Run("update", update(ts, &payment))
	t.Run("delete", del(ts, payment.Id.Hex(), col))
	// order matters ^
}

func del(ts *httptest.Server, paymentID string, paymentsCollection memory.Storage) func(t *testing.T) {
	return func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, getURL(ts, api.PaymentsPath+"/"+paymentID), nil)
		req.Header.Set("Content-Type", "application/json")

		reqInfo := getRequestInfo(*req)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("error performing http request in %s (%s):\n%s", t.Name(), reqInfo, err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d (%s)", http.StatusOK, res.StatusCode, reqInfo)
		}

		body, _ := ioutil.ReadAll(res.Body)
		t.Logf("got reply in %s (%s):\n%s", t.Name(), reqInfo, string(body))

		// short circuiting this in order not to call list again
		if len(paymentsCollection) != 0 {
			t.Fatalf("item %s was not deleted (%s)", paymentID, reqInfo)
		}
	}
}

func update(ts *httptest.Server, existingPayment *business.Payment) func(t *testing.T) {
	return func(t *testing.T) {
		paymentToUpdate := *existingPayment

		paymentToUpdate.Version = 7
		jsonRepresentation := getJSON(paymentToUpdate)

		req, _ := http.NewRequest(
			http.MethodPut,
			getURL(ts, api.PaymentsPath+"/"+paymentToUpdate.Id.Hex()),
			bytes.NewBufferString(jsonRepresentation),
		)
		req.Header.Set("Content-Type", "application/json")

		reqInfo := getRequestInfo(*req)

		t.Logf("requesting in %s (%s):\n%s", t.Name(), reqInfo, jsonRepresentation)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("error performing http request in %s (%s):\n%s", t.Name(), reqInfo, err)
		}

		body, _ := ioutil.ReadAll(resp.Body)
		if len(body) == 0 {
			t.Fatalf("response was empty in %s (%s)", t.Name(), reqInfo)
		}

		t.Logf("got reply in %s (%s):\n%s", t.Name(), reqInfo, string(body))

		if *existingPayment, err = paymentFromJSON(body); err != nil {
			t.Fatalf("could not decode [%s] response json: %s", reqInfo, err)
		}

		if !reflect.DeepEqual(paymentToUpdate, *existingPayment) {
			t.Fatalf("expected response payment to be \n%+v\n, got\n%+v", paymentToUpdate, *existingPayment)
		}
	}
}

func retrieve(ts *httptest.Server, createdPaymentID string) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := http.Get(getURL(ts, api.PaymentsPath+"/doesnt-exist"))
		if err != nil {
			t.Fatalf("error performing http request in %s:\n%s", t.Name(), err)
		}
		if res.StatusCode != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, res.StatusCode)
		}

		targetURL := getURL(ts, api.PaymentsPath+"/"+createdPaymentID)
		res, err = http.Get(targetURL)
		if err != nil {
			t.Fatalf("error performing http request in %s:\n%s", t.Name(), err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, res.StatusCode)
		}

		body, _ := ioutil.ReadAll(res.Body)

		t.Logf("got reply in %s:\n%s", t.Name(), string(body))

		payment := new(business.Payment)
		if err := json.Unmarshal(body, payment); err != nil {
			errStr := err.Error()
			t.Fatalf("could not decode %s response json: %s", t.Name(), errStr)
		}

		t.Logf("retrieved %v", payment)
	}
}

func list(ts *httptest.Server, existingPayment business.Payment) func(t *testing.T) {
	return func(t *testing.T) {
		var (
			listPayments []*business.Payment
			err          error
		)

		res, err := http.Get(getURL(ts, api.PaymentsPath))
		if err != nil {
			t.Fatalf("error performing http request in %s:\n%s", t.Name(), err)
		}

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, res.StatusCode)
		}

		body, _ := ioutil.ReadAll(res.Body)

		if err = json.Unmarshal(body, &listPayments); err != nil {
			t.Fatalf("could not decode list response json: %s", err)
		}

		if len(listPayments) != 1 {
			t.Fatalf("retrieved list has %d items instead of %d", len(listPayments), 1)
		}

		for _, payment := range listPayments {
			if !reflect.DeepEqual(*payment, existingPayment) {
				t.Fatalf("expected payment retrieved from the list endpoint to be \n%v\n, got\n%v",
					existingPayment, payment)
			}
			break
		}
	}
}

func create(ts *httptest.Server, createdPayment *business.Payment) func(t *testing.T) {
	return func(t *testing.T) {
		paymentReq := business.Payment{Version: 13, Attributes: &business.Attributes{Amount: 123}}
		jsonRepresentation := getJSON(paymentReq)
		resp, err := http.Post(
			getURL(ts, api.PaymentsPath),
			"application/json",
			bytes.NewBufferString(jsonRepresentation),
		)
		if err != nil {
			t.Fatalf("error performing http request in %s:\n%s", t.Name(), err)
		}

		reqInfo := getRequestInfo(*resp.Request)

		t.Logf("requesting in %s:\n%s", t.Name(), jsonRepresentation)

		body, _ := ioutil.ReadAll(resp.Body)
		if len(body) == 0 {
			t.Fatalf("response was empty in %s (%s)", t.Name(), reqInfo)
		}

		t.Logf("got reply in %s (%s):\n%s", t.Name(), string(body), reqInfo)

		if *createdPayment, err = paymentFromJSON(body); err != nil {
			t.Fatalf("could not decode response json (%s): %s", reqInfo, err)
		}

		id := createdPayment.Id
		createdPayment.Id = bson.ObjectId("") //  don't compare ids
		if !reflect.DeepEqual(paymentReq, *createdPayment) {
			t.Fatalf("expected response payment (%s) to be \n%+v\n, got\n%+v", reqInfo, paymentReq, *createdPayment)
		}
		createdPayment.Id = id
	}
}

func getURL(ts *httptest.Server, relativePath string) string {
	return ts.URL + relativePath
}

func getJSON(payment business.Payment) string {
	b, _ := json.Marshal(payment)
	return string(b)
}

func paymentFromJSON(jsonned []byte) (payment business.Payment, err error) {
	err = json.Unmarshal(jsonned, &payment)
	return
}

func getRequestInfo(req http.Request) string {
	return fmt.Sprintf("%s %s", req.Method, req.URL)
}
