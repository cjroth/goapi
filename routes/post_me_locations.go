package routes

import (
    "github.com/cjroth/even/lib"
    "github.com/cjroth/even/models"
    "net/http"
    "fmt"
    "encoding/json"
    _ "github.com/lib/pq"
    "code.google.com/p/go-uuid/uuid"
)

// @todo integrate authentication system and assign correct user id based on token
// @todo use postgres native timestamp data type instead of integers

// curl localhost:8080/me/locations -d '{ "time": 1419022733, "lat": 123.123, "lon": 123.123 }' -H "Authorization: Bearer mytoken"
// 204 <empty body>
func PostMeLocationsHandler(c lib.Context, w http.ResponseWriter, r *http.Request) {

    decoder := json.NewDecoder(r.Body)
    var loc models.Location
    err := decoder.Decode(&loc)
    lib.HandleError(err, "json decoding failed")

    loc.Id = uuid.New()
    loc.UserId = uuid.New()

    c.LocationQueue <- loc

    // insert to database for long-term storage
    // would be useful if we wanted to re-aggregate shifts from
    // location data with a different or improved algorithm
    err = c.Dbmap.Insert(&loc)
    lib.HandleError(err, "Insert failed")

    w.WriteHeader(http.StatusNoContent)

    fmt.Println("user location received", loc)

}