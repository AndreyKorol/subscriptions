package models

type Subscription struct {
    Id uint `json:"id" db:"id"`
    ServiceName string `json:"service_name" db:"service_name"`
    Price uint `json:"price" db:"price"`
    UserId string `json:"user_id" db:"user_id"`
    StartDate string `json:"start_date" db:"start_date"`
}
