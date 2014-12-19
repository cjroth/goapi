package main

import (
    "github.com/cjroth/even/lib"
    "github.com/cjroth/even/routes"
    "github.com/cjroth/even/models"
    "net/http"
    "log"
    "os"
    // "fmt"
    // "time"
    "code.google.com/p/go-uuid/uuid"
    "github.com/gorilla/mux"
    "github.com/coopernurse/gorp"
    _ "github.com/lib/pq"
    "database/sql"
)

var dbmap *gorp.DbMap

func main() {


    var LocationQueue = make(chan models.Location, 100)

    dbmap = initDb()

    context := lib.Context{Dbmap: dbmap, LocationQueue: LocationQueue}

    r := mux.NewRouter()

    r.HandleFunc("/register", lib.GetHandlerWithContext(routes.RegisterHandler, context)).
        Methods("POST")

    // create sub router for /users/*
    users := r.PathPrefix("/users").Subrouter()
    users.HandleFunc("/{id}/shifts/{start}/{end}", lib.GetHandlerWithContext(routes.GetUsersShiftsHandler, context)).
        Methods("GET")

    // create sub router for /me/*
    me := r.PathPrefix("/me").Subrouter()
    me.HandleFunc("/locations", lib.GetHandlerWithContext(routes.PostMeLocationsHandler, context)).
        Methods("POST")

    http.Handle("/", r)


    // create worker pool for processing location data
    // for w := 1; w <= 3; w++ {
    //     go locationQueueWorker(w, LocationQueue)
    // }


    log.Fatal(http.ListenAndServe(":8080", nil))

}

func locationQueueWorker(userId string) {

    // fetch all rows
    var locations []models.Location
    _, err := dbmap.Select(&locations, "select * from locations where user_id = $1", userId)
    lib.HandleError(err, "Select failed")

    y := 0
    n := 0
    t := 3 // number of "at work" or "not at work" pings in a row before detecting
           // that a user has begun or ended a shift at work

    start := 0
    end := 0
    duration := 0

    shifts := []models.Shift{}

    for i, loc := range locations {
        
        // increment the threshold counter
        // t "at work" data points in a row indicates beginning of shift
        // t "not at work" data points in a row indicates end of shift
        if isAtWork(loc) {
            y++
        } else {
            n++
        }

        // detected start of a shift
        if y == t {
            n = 0
            start = locations[i-t].Timestamp
        }

        // detected end of a shift; push shift to slice and zero start and end
        if start != 0 && n == t {
            y = 0
            end = locations[i-t].Timestamp
            duration = end - start
            shifts = append(shifts, models.Shift{
                Id: uuid.New(),
                StartTime: start,
                Duration: duration,
                UserId: loc.UserId,
            })
            start = 0
            end = 0
        }

    }

    // why doesn't `dbmap.Insert(shifts...)` work?
    for _, shift := range shifts {
        dbmap.Insert(shift)
    }

}

// return true if user is within x meters of the coordinates of their place of work
func isAtWork(loc models.Location) bool {
    return true
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
    dbmap.AddTableWithName(models.Location{}, "locations")
        // .SetKeys(true, "Id")

    // create the table. in a production system you'd generally
    // use a migration tool, or create the tables via scripts
    // err = dbmap.CreateTablesIfNotExists()
    // lib.HandleError(err, "Create tables failed")

    return dbmap
}

