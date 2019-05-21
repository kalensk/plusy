package testutils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type FakeTelegramServer2 struct {
	testing    *testing.T
	testServer *httptest.Server
	router     *mux.Router
	handlers   map[string][]func(http.ResponseWriter, *http.Request)
}

func NewFakeTelgramServer2(testing *testing.T) *FakeTelegramServer2 {
	return &FakeTelegramServer2{
		testing:  testing,
		router:   mux.NewRouter(),
		handlers: make(map[string][]func(http.ResponseWriter, *http.Request))}
}

func (f *FakeTelegramServer2) Start() string {
	//for path, handler := range f.handlers {
	//f.router.HandleFunc(path, handler)
	//}

	f.testServer = httptest.NewUnstartedServer(f.router)

	f.testServer.Start()
	return f.testServer.URL
}

func (f *FakeTelegramServer2) Url() string {
	return f.testServer.URL
}

func (f *FakeTelegramServer2) Close() {
	f.testServer.Close()
}

func (f *FakeTelegramServer2) AddHandler(path string, handlerFunc func(http.ResponseWriter, *http.Request)) {
	f.handlers[path] = append(f.handlers[path], handlerFunc)
}
