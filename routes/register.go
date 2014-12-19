package routes

import (
    "github.com/cjroth/even/lib"
    "github.com/cjroth/even/models"
    "net/http"
    "fmt"
    "encoding/json"
    _ "github.com/lib/pq"
    "code.google.com/p/go-uuid/uuid"
    "golang.org/x/crypto/bcrypt"
)

// curl localhost:8080/register -d '{ "id": "123", "email": "chris@cjroth.com", "password": "test" }'
// {"id":"a35ac905-d777-43cb-b5f9-0579648ce38f","email":"chris@cjroth.com","password":"$2a$10$GshhwSaL9ZipvDw6.a.xbuRoh6zdIRYl37bqalYzRdJBB61MDkGH."}
func RegisterHandler(c lib.Context, w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var u models.User
    err := decoder.Decode(&u)
    lib.HandleError(err, "json decoding failed")

    hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)

    u.Id = uuid.New()
    u.Password = string(hash)

    err = c.Dbmap.Insert(&u)
    lib.HandleError(err, "Insert failed")

    w.Header().Set("Content-Type", "application/json")

    encoder := json.NewEncoder(w)
    err = encoder.Encode(&u)
    lib.HandleError(err, "json encoding failed")

    fmt.Println("user registered", u)

}