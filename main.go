package main

import (
    "github.com/cjroth/even/lib"
    "github.com/cjroth/even/routes"
    "github.com/cjroth/even/models"
    "net/http"
    "log"
    "os"
    "github.com/gorilla/mux"
    "github.com/coopernurse/gorp"
    _ "github.com/lib/pq"
    "database/sql"
)

func main() {

    dbmap := initDb()

    context := lib.Context{Dbmap: dbmap}

    r := mux.NewRouter()

    r.HandleFunc("/register", lib.GetHandlerWithContext(routes.RegisterHandler, context)).
        Methods("POST")

    users := r.PathPrefix("/users").Subrouter()
    users.HandleFunc("/{id}/shifts/{start}/{end}", lib.GetHandlerWithContext(routes.UsersShiftHandler, context)).
        Methods("GET")

    http.Handle("/", r)

    log.Fatal(http.ListenAndServe(":8080", nil))

}

func initDb() *gorp.DbMap {
    // connect to db using standard Go database/sql API
    // use whatever database/sql driver you wish
    db, err := sql.Open("postgres", "user=chris dbname=even sslmode=disable")
    lib.HandleError(err, "sql.Open failed")

    // construct a gorp DbMap
    dbmap := &gorp.DbMap{ Db: db, Dialect: gorp.PostgresDialect{} }

    dbmap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds)) 


    // add a table, setting the table name to 'posts' and
    // specifying that the Id property is an auto incrementing PK
    dbmap.AddTableWithName(models.User{}, "users")
        // .SetKeys(true, "Id")

    // create the table. in a production system you'd generally
    // use a migration tool, or create the tables via scripts
    // err = dbmap.CreateTablesIfNotExists()
    // lib.HandleError(err, "Create tables failed")

    return dbmap
}
