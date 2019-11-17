package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"form3/api"
	"form3/business"
	"form3/business/sqlite"

	logrusTest "github.com/sirupsen/logrus/hooks/test"
	"gopkg.in/mgo.v2/bson"
)

// todo this ensures some basic protection against regressions, but it needs more edge cases
// todo integration testing
func TestCreateRetrieve(t *testing.T) {
	// implementation of model.Storage
	store, err := sqlite.New(":memory:")
	if err != nil {
		log.Fatalf("can't initialize storage: %w", err)
	}
	logger, _ := logrusTest.NewNullLogger()
	service := api.NewService(store, *logger)

	// some of these tests can also work without opening this local listener
	// by testing individual handlers rather than this server instance
	ts := httptest.NewServer(service.GetMux())
	defer ts.Close()

	// this will be used to keep some state between the tests below
	var obj business.Employee

	// isolating tests into closures with deps clearly outlined. order matters.
	t.Run("create", create(ts, &obj))
	t.Run("retrieve", retrieve(ts, obj.Id.Hex()))
}

func create(ts *httptest.Server, created *business.Employee) func(t *testing.T) {
	return func(t *testing.T) {
		req := business.Employee{Version: 13, Attributes: &business.Attributes{Amount: 123}}
		jsonRepresentation := getJSON(req)
		resp, err := http.Post(
			getURL(ts, api.RootPath),
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

		if *created, err = fromJSON(body); err != nil {
			t.Fatalf("could not decode response json (%s): %s", reqInfo, err)
		}

		id := created.Id
		created.Id = bson.ObjectId("") //  don't compare ids
		if !reflect.DeepEqual(req, *created) {
			t.Fatalf("expected response obj (%s) to be \n%+v\n, got\n%+v", reqInfo, req, *created)
		}
		created.Id = id
	}
}

func retrieve(ts *httptest.Server, createdObjID string) func(t *testing.T) {
	return func(t *testing.T) {
		res, err := http.Get(getURL(ts, api.RootPath+"/doesnt-exist"))
		if err != nil {
			t.Fatalf("error performing http request in %s:\n%s", t.Name(), err)
		}
		if res.StatusCode != http.StatusNotFound {
			t.Fatalf("expected status %d, got %d", http.StatusNotFound, res.StatusCode)
		}

		targetURL := getURL(ts, api.RootPath+"/"+createdObjID)
		res, err = http.Get(targetURL)
		if err != nil {
			t.Fatalf("error performing http request in %s:\n%s", t.Name(), err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, res.StatusCode)
		}

		body, _ := ioutil.ReadAll(res.Body)

		t.Logf("got reply in %s:\n%s", t.Name(), string(body))

		obj := new(business.Employee)
		if err := json.Unmarshal(body, obj); err != nil {
			errStr := err.Error()
			t.Fatalf("could not decode %s response json: %s", t.Name(), errStr)
		}

		t.Logf("retrieved %v", obj)
	}
}

// helpers:

func getURL(ts *httptest.Server, relativePath string) string {
	return ts.URL + relativePath
}

func getJSON(obj business.Employee) string {
	b, _ := json.Marshal(obj)
	return string(b)
}

func fromJSON(jsonned []byte) (obj business.Employee, err error) {
	err = json.Unmarshal(jsonned, &obj)
	return
}

func getRequestInfo(req http.Request) string {
	return fmt.Sprintf("%s %s", req.Method, req.URL)
}
