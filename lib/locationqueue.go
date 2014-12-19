package lib

import (
  "github.com/cjroth/even/models"
)

var LocationQueue = make(chan models.Location, 100)

