package fcomposite

import(
    "fmt"
    "bytes"
    "net/http"
    "github.com/ricallinson/forgery"
)

// The Writer used in place of stackr.Writer for buffering.
type BufferedResponseWriter struct {
    Headers http.Header
    Buffer *bytes.Buffer
    Status int
}

// Header returns the header map that will be sent by WriteHeader.
// Changing the header after a call to WriteHeader (or Write) has
// no effect.
func (this *BufferedResponseWriter) Header() (http.Header) {
    if this.Headers == nil {
        this.Headers = http.Header{}
    }
    return this.Headers
}

// Write writes the data to the connection as part of an HTTP reply.
// If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
// before writing the data.  If the Header does not contain a
// Content-Type line, Write adds a Content-Type set to the result of passing
// the initial 512 bytes of written data to DetectContentType.
func (this *BufferedResponseWriter) Write(b []byte) (int, error) {
    if this.Buffer == nil {
        this.Buffer = &bytes.Buffer{}
    }
    len, err := this.Buffer.Write(b)
    fmt.Println(string(this.Buffer.Bytes()))
    return len, err
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
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
func (this Map) Dispatch(req *f.Request, res *f.Response, next func()) (map[string]string) {

    headers := map[string]http.Header{}
    renders := map[string]string{}

    // Grab the res.Writer so we can put it back later.
    w := res.Response.Writer

    c := make(chan int, len(this))
    for id, fn := range this {
        go func(mapId string, mapFn func(*f.Request, *f.Response, func())) {
            // Clone the res so it can be changed in isolation.
            response := res.Clone()
            // Create a buffer.
            buffer := &BufferedResponseWriter{}
            // Replace res.Writer with BufferedResponseWriter so all the output can be captured.
            response.Writer = buffer
            // Call the function.
            mapFn(req, response, next)
            // Add the buffered data to the renders map.
            if buffer.Buffer != nil {
                renders[mapId] = buffer.Buffer.String()
            }
            // Add the buffered headers to the headers map.
            headers[mapId] = buffer.Headers
            // Return the channel
            c <- 1
        }(id, fn)
    }
    // Wait for all the channels to close.
    <-c

    // Put the res.Writer back.
    res.Response.Writer = w

    // Write the headers to the real response.
    // TODO

    return renders
}