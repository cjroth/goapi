package models

type Location struct {
    Id string           `db:"id" json:"id"`
    Lat float32         `db:"lat" json:"lat"`
    Lon float32         `db:"lon" json:"lon"`
    Timestamp int       `db:"timestamp" json:"timestamp"`
    UserId string       `db:"user_id" json:"user_id"`
}