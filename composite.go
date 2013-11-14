package fcomposite

import(
    "github.com/ricallinson/forgery"
)

type Map map[string]func(*f.Request, *f.Response, func())

/*
    composite := fcomposite.Map{
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

    data := composite.Dispatch(req, res, next)
*/
func (this *Map) Dispatch(req *f.Request, res *f.Response, next func()) (map[string]string) {

    // Replace res.Writer with something to buffer the data.

    return map[string]string{}
}