package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/timunas/ldt/server/middleware"
	"github.com/urfave/negroni"
)

type TestMock struct {
	mock.Mock
}

func (t *TestMock) DoSomething(rw http.ResponseWriter, r *http.Request) {
	t.Called(rw, r)
}

func TestTodoCreation(t *testing.T) {
	testObj := new(TestMock)
	handler := func(rw http.ResponseWriter, r *http.Request) {
		testObj.DoSomething(rw, r)
	}

	r := httptest.NewRequest("GET", "http://example.com/foo", nil)
	rw := negroni.NewResponseWriter(httptest.NewRecorder())

	testObj.On("DoSomething", rw, r).Return()

	middleware.LoggingMiddleware(rw, r, handler)

	testObj.AssertCalled(t, "DoSomething", rw, r)
}
