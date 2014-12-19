package routes

import (
    "github.com/cjroth/even/lib"
    "net/http"
    "fmt"
    "github.com/gorilla/mux"
)

func GetUsersShiftsHandler(c lib.Context, res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    id := vars["id"]
    start := vars["start"]
    end := vars["end"]
    fmt.Println("vars", id, start, end)
    res.Header().Set("Content-Type", "text/plain")
    res.Write([]byte("This is an example server.\n"))
}