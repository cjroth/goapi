package lib

import (
    "net/http"
    "github.com/cjroth/even/models"
    "github.com/coopernurse/gorp"
)

type Context struct {
    Dbmap *gorp.DbMap
    LocationQueue chan models.Location
}

func GetHandlerWithContext(h func(Context, http.ResponseWriter, *http.Request), c Context) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        h(c, w, r)
    }
}