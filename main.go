package main

import (
    "net/http"
    "log"
    "fmt"
    "os"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/coopernurse/gorp"
    _ "github.com/lib/pq"
    "code.google.com/p/go-uuid/uuid"
    "golang.org/x/crypto/bcrypt"
    "database/sql"
)

var dbmap *gorp.DbMap

type User struct {
    Id string           `db:"id" json:"id"`
    Email string        `db:"email" json:"email"`
    Password string     `db:"password" json:"password"`
}

func UsersShiftHandler(res http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    id := vars["id"]
    start := vars["start"]
    end := vars["end"]
    fmt.Println("vars", id, start, end)
    res.Header().Set("Content-Type", "text/plain")
    res.Write([]byte("This is an example server.\n"))
}

func RegisterHandler(res http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)
    var u User
    err := decoder.Decode(&u)
    if err != nil {
        panic(err)
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)

    u.Id = uuid.New()
    u.Password = string(hash)

    log.Println(u)

    err = dbmap.Insert(&u)
    handleError(err, "Insert failed")

    res.Header().Set("Content-Type", "application/json")

    encoder := json.NewEncoder(res)
    err = encoder.Encode(&u)

    fmt.Println("user registered", u)

}

func main() {

    dbmap = initDb()

    r := mux.NewRouter()

    r.HandleFunc("/register", RegisterHandler).
        Methods("POST")

    users := r.PathPrefix("/users").Subrouter()
    users.HandleFunc("/{id}/shifts/{start}/{end}", UsersShiftHandler)

    http.Handle("/", r)

    log.Fatal(http.ListenAndServe(":8080", nil))

}

func initDb() *gorp.DbMap {
    // connect to db using standard Go database/sql API
    // use whatever database/sql driver you wish
    db, err := sql.Open("postgres", "user=chris dbname=even sslmode=disable")
    handleError(err, "sql.Open failed")

    // construct a gorp DbMap
    dbmap := &gorp.DbMap{ Db: db, Dialect: gorp.PostgresDialect{} }

    dbmap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds)) 


    // add a table, setting the table name to 'posts' and
    // specifying that the Id property is an auto incrementing PK
    dbmap.AddTableWithName(User{}, "users")
        // .SetKeys(true, "Id")

    // create the table. in a production system you'd generally
    // use a migration tool, or create the tables via scripts
    // err = dbmap.CreateTablesIfNotExists()
    // handleError(err, "Create tables failed")

    return dbmap
}

func handleError(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}
