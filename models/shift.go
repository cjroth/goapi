package models

type Shift struct {
    Id string           `db:"id" json:"id"`
    StartTime int      `db:"start_time" json:"start_time"`
    Duration int       `db:"duration" json:"duration"`
    UserId string       `db:"user_id" json:"user_id"`
}