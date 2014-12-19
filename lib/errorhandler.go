package lib

import (
    "log"
)

func HandleError(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}