package fcomposite

import (
	"bytes"
	"fmt"
	"github.com/ricallinson/forgery"
	"net/http"
)

// The Writer used in place of stackr.Writer for buffering.
type BufferedResponseWriter struct {
	Headers http.Header
	Buffer  *bytes.Buffer
	Status  int
}

// Header returns the header map that would have been sent by WriteHeader.
func (this *BufferedResponseWriter) Header() http.Header {
	if this.Headers == nil {
		this.Headers = http.Header{}
	}
	return this.Headers
}

// Write writes the data to the buffer.
func (this *BufferedResponseWriter) Write(b []byte) (int, error) {
	if this.Buffer == nil {
		this.Buffer = &bytes.Buffer{}
	}
	len, err := this.Buffer.Write(b)
	fmt.Println(string(this.Buffer.Bytes()))
	return len, err
}

// WriteHeader buffers the HTTP status code.
func (this *BufferedResponseWriter) WriteHeader(code int) {
	this.Status = code
}

/*
   composite := fcomposite.Map{
       "header": func(req, res, next) {
           res.Send("Header string")
       },
       "body": func(req, res, next) {
           res.Render("page.html", "Body string")
       },
       "footer": func(req, res, next) {
           res.End("Footer string")
       },
       "tail": func(req, res, next) {
           res.Write("Tail string")
       },
       "close": func(req, res, next) {
           res.WriteBytes([]byte("Close string"))
       },
   }

   data := composite.Dispatch(req, res, next)
*/
type Map map[string]func(*f.Request, *f.Response, func())

/*
   The worker.
*/
func (this Map) Dispatch(req *f.Request, res *f.Response, next func()) map[string]string {

	headers := map[string]http.Header{}
	locals := map[string]map[string]string{}
	renders := map[string]string{}

	// Grab the res.Writer so we can put it back later.
	w := res.Response.Writer

	c := make(chan int, len(this))
	// Loop over the items in the map.
	for id, fn := range this {
		// Dispatch function.
		go func(mapId string, mapFn func(*f.Request, *f.Response, func())) {
			// Clone the res so it can be changed in isolation.
			response := res.Clone()
			// Create a buffer.
			buffer := &BufferedResponseWriter{}
			// Replace res.Writer with BufferedResponseWriter so all the output can be captured.
			response.Response.Writer = buffer
			// Call the function.
			mapFn(req, response, next)
			// Add the buffered data to the renders map.
			if buffer.Buffer != nil {
				renders[mapId] = buffer.Buffer.String()
			}
			// Add the buffered Headers to the headers map.
			headers[mapId] = buffer.Headers
			// Add the buffered Locals to the locals map.
			locals[mapId] = response.Locals
			// Return the channel
			c <- 1
		}(id, fn)
	}
	// Wait for all the channels to close.
	<-c

	// Put the res.Writer back.
	res.Response.Writer = w

	// Write the headers to the real response.
	// Note, only the first value is transfered from the given http.Header map.
	// Note, no ordering is guaranteed so headers could be overridden randomly.
	// This is not going to work. Some intelligence has to be introduced for Cookies.
	for _, mapHeaders := range headers {
		for f, v := range mapHeaders {
			// res.Writer.Header()[f] = v
			res.Set(f, v[0])
		}
	}

	// Write the Locals to the real response.
	// Locals are appended to each other if they share the same key.
	// Note, no ordering is guaranteed so Locals will be appended randomly.
	for _, mapLocals := range locals {
		for f, v := range mapLocals {
			if _, ok := res.Locals[f]; ok {
				res.Locals[f] += v
			} else {
				res.Locals[f] = v
			}
		}
	}

	return renders
}
