package models

type Subscription struct {
    Id          uint   `json:"id"           db:"id"`
    ServiceName string `json:"service_name" db:"service_name"`
    Price       uint   `json:"price"        db:"price"`
    UserId      string `json:"user_id"      db:"user_id"`
    StartDate   string `json:"start_date"   db:"start_date"`
}

type AggSubscriptions struct {
    sumPrice int `json:"sum_price"`
}

type Filter struct {
    ServiceName *string `schema:"service_name" validate:"omitempty,ne="`
    UserId      *string `schema:"user_id"      validate:"omitempty,uuid4"`
    StartDate   *string `schema:"start_date"   validate:"omitempty,datetime=01-2006"`
    EndDate     *string `schema:"end_date"     validate:"omitempty,datetime=01-2006"`
}
