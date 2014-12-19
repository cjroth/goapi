package models

type Shift struct {
    Id string           `db:"id" json:"id"`
    StartTime uint      `db:"start_time" json:"start_time"`
    Duration uint       `db:"duration" json:"duration"`
    UserId string       `db:"user_id" json:"user_id"`
}