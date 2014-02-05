package fcomposite

import (
	"github.com/ricallinson/forgery"
	. "github.com/ricallinson/simplebdd"
	"github.com/ricallinson/stackr"
	"net/http"
	"net/url"
	"testing"
)

func TestFcomposite(t *testing.T) {

	var mock *MockResponseWriter
	var app *f.Server
	var req *f.Request
	var res *f.Response

	BeforeEach(func() {
		mock = NewMockResponseWriter(false)
		app = f.CreateServer()
		req = &f.Request{
			Request: &stackr.Request{
				Request: &http.Request{
					URL:    &url.URL{},
					Header: map[string][]string{},
				},
				Query: map[string]string{},
			},
		}
		res = &f.Response{
			Response: &stackr.Response{
				Writer: mock,
			},
			Charset: "utf-8",
			Locals:  map[string]string{},
		}
		res.SetApplication(app)
		res.SetRequest(req)
		res.SetNext(func() {})
	})

	Describe("Dispatch()", func() {

		It("should return win", func() {

			composite := Map{
				"header": func(req *f.Request, res *f.Response, next func()) {
					res.Send("Header")
				},
				"empty": func(req *f.Request, res *f.Response, next func()) {
					res.Locals["append"] = "Bar"
					res.End("")
				},
				"body": func(req *f.Request, res *f.Response, next func()) {
					res.Locals["title"] = "Foo"
					res.Render("test.html", "Body")
				},
				"footer": func(req *f.Request, res *f.Response, next func()) {
					res.Locals["append"] = "Bar"
					res.End("Footer")
				},
				"tail": func(req *f.Request, res *f.Response, next func()) {
					res.Write("Tail")
				},
				"close": func(req *f.Request, res *f.Response, next func()) {
					res.WriteBytes([]byte("Close"))
				},
			}

			fn := func() {}

			app.Engine(".html", &MockRenderer{})

			data := composite.Dispatch(req, res, fn)

			AssertEqual(string(data["header"]), "Header")
			AssertEqual(string(data["empty"]), "")
			AssertEqual(string(data["body"]), "Body")
			AssertEqual(res.Locals["title"], "Foo")
			AssertEqual(string(data["footer"]), "Footer")
			AssertEqual(string(data["tail"]), "Tail")
			AssertEqual(string(data["close"]), "Close")
			AssertEqual(res.Locals["append"], "BarBar")
		})
	})

	Report(t)
}
