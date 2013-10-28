package composite

import(
    "github.com/ricallinson/forgery"
)

type Map map[string]func(*Request, *Response, func())

/*
    cfg := composite.Map{
        "header": func(req, res, next) {
            res.End("Header string")
        },
        "body": func(req, res, next) {
            res.End("Body string")
        },
        "footer": func(req, res, next) {
            res.End("Footer string")
        },
    }

    data := cfg.Dispatch(req, res, next)
*/
func (this *Map) Dispatch(req *forgery.Request, res *forgery.Response, next func()) (map[string]string) {
    // Replace res.Writer with something to buffer the data.
}