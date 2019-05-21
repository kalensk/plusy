package testutils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Behavior struct {
	handler   http.HandlerFunc
	assertion []byte
}

type FakeTelegramServer struct {
	testing     *testing.T
	testServer  *httptest.Server
	currentPath string
	behaviors   map[string][]Behavior // [path][]Behavior
}

/*
	queue := make([]int, 0)
	// Push to the queue
	queue = append(queue, 1)
	// Top (just get next element, don't remove it)
	x = queue[0]
	// Discard top element
	queue = queue[1:]
	// Is empty ?
	if len(queue) == 0 {
	fmt.Println("Queue is empty !")
	}
*/

func NewFakeTelgramServer(testing *testing.T) *FakeTelegramServer {
	return &FakeTelegramServer{testing: testing, behaviors: make(map[string][]Behavior)}
}

func (f *FakeTelegramServer) Start() string {
	f.testServer = httptest.NewServer(http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		splitPath := strings.Split(request.URL.Path, "/")
		path := splitPath[len(splitPath)-1]

		requestBodyBytes, err := ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println("Inside Main HandlerFunc")
		fmt.Println(request.Method + " " + request.URL.Path + " -> " + path)
		fmt.Println("request body received: " + string(requestBodyBytes))
		fmt.Println("======================")

		behaviors, ok := f.behaviors[path]
		if !ok {
			respWriter.WriteHeader(404)
			return
		}

		behaviors[0].handler(respWriter, request)
		f.behaviors[request.URL.Path] = behaviors[1:]
	}))

	return f.testServer.URL

}

func (f *FakeTelegramServer) Url() string {
	return f.testServer.URL
}

func (f *FakeTelegramServer) Close() {
	f.testServer.Close()
}

func (f *FakeTelegramServer) ForCallTo(path string) *FakeTelegramServer {
	f.currentPath = path
	return f
}

func (f *FakeTelegramServer) ReturnString(body []byte) *FakeTelegramServer {
	handler := func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write(body)

		requestBodyBytes, err := ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println("Inside testServer")
		fmt.Println(request.Method + " " + request.URL.Path)
		fmt.Println("request body received: " + string(requestBodyBytes))
		fmt.Println("writing body: " + string(body))
		fmt.Println("======================")

		splitPath := strings.Split(request.URL.Path, "/")
		path := splitPath[len(splitPath)-1]

		f.behaviors[path] = f.behaviors[path][1:] // try changing path to f.currentPath...
	}

	currentBehaviors := f.behaviors[f.currentPath]
	currentBehaviors = append(currentBehaviors, Behavior{handler: handler})
	f.behaviors[f.currentPath] = currentBehaviors
	return f
}

func (f *FakeTelegramServer) ReturnStringAndAssert(body []byte, expectedRequestBody []byte) *FakeTelegramServer {
	handler := func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write(body)

		requestBodyBytes, err := ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println("Inside testServer")
		fmt.Println(request.Method + " " + request.URL.Path)
		fmt.Println("request body received: " + string(requestBodyBytes))
		fmt.Println("writing body: " + string(body))
		fmt.Println("======================")

		assert.Equal(f.testing, string(expectedRequestBody), string(requestBodyBytes),
			"Request body did not match for path: %s", f.currentPath)

		splitPath := strings.Split(request.URL.Path, "/")
		path := splitPath[len(splitPath)-1]

		f.behaviors[path] = f.behaviors[path][1:] // try chaning path to f.currentPath...
	}

	currentBehaviors := f.behaviors[f.currentPath]
	currentBehaviors = append(currentBehaviors, Behavior{handler: handler, assertion: expectedRequestBody})
	f.behaviors[f.currentPath] = currentBehaviors
	return f
}

/////////////////////////////////////

//type DongleServer struct {
//	testing   *testing.T
//	server    *httptest.Server
//	behaviors map[string]http.HandlerFunc // map to a queue of funcs (pop off a func each time its called)
//	// sine you can have getUpdates twice in one call.
//
//	currentBehavior struct {
//		path      string
//		handler   http.HandlerFunc
//		assertion string
//	}
//}
//
//func NewDongleServer(testing *testing.T) *DongleServer {
//	var dongle DongleServer
//	dongle.testing = testing
//	dongle.behaviors = make(map[string]http.HandlerFunc)
//	dongle.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		handler, ok := dongle.behaviors[r.URL.Path]
//		if !ok {
//			w.WriteHeader(404)
//			return
//		}
//		handler(w, r)
//	}))
//	return &dongle
//}
//func (d *DongleServer) when(path string, handler http.HandlerFunc) {
//	d.behaviors[path] = handler
//}
//
//func (d *DongleServer) once(path string, handler http.HandlerFunc) {
//	d.when(path, func(w http.ResponseWriter, r *http.Request) {
//		handler(w, r)
//		delete(d.behaviors, path)
//	})
//}
//
//func (d *DongleServer) returnOnceString(path string, body string) {
//	d.behaviors[path] = func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusOK)
//		w.Write([]byte(body))
//		delete(d.behaviors, path)
//	}
//}
//
//func (d *DongleServer) whenString(path string, body string) {
//	d.behaviors[path] = func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusOK)
//		w.Write([]byte(body))
//	}
//}
//
//func (d *DongleServer) forCallTo(path string) *DongleServer {
//	d.currentBehavior.path = path
//	return d
//}
//
//func (d *DongleServer) returnString(body []byte) *DongleServer {
//	d.currentBehavior.handler = func(w http.ResponseWriter, r *http.Request) {
//		w.WriteHeader(http.StatusOK)
//		w.Write(body)
//		delete(d.behaviors, d.currentBehavior.path)
//	}
//
//	d.behaviors[d.currentBehavior.path] = d.currentBehavior.handler
//	return d
//}
//
//func (d *DongleServer) andAssertThat(expectedBodyResponse string) *DongleServer {
//	d.currentBehavior.assertion = expectedBodyResponse
//	actualBodyResponse := request.Body.Read()
//	assert.Equal(d.testing, expectedBodyResponse, actualBodyResponse)
//	return d
//}
//
//func (d *DongleServer) Close() {
//	d.server.Close()
//}
//
//func Test_Stuff(t *testing.T) {
//	dongle := NewDongleServer(t)
//	defer dongle.Close()
//	dongle.whenString("getUpdates", "[1,2,3]")
//	dongle.whenString("sendMessage", "message sent")
//	// plusy does stuff
//	dongle.returnOnceString("getUpdates", "[4,5,6]")
//
//	dongle.forCallTo("getUpdates").returnString("").andAssertThat("asd")
//
//	//dongle.forPath("getUpdates").thenReturnString("")
//	//dongle.forPath("getUpdates").thenReturn("")
//	//
//	//dongle.forPath("getUpdates").thenOnlyReturnString("")
//	//dongle.forPath("getUpdates").thenOnlyReturn("")
//}
