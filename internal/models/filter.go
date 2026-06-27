package models

import(
    "github.com/google/uuid"
)

type Filter struct {
    UserId uuid.UUID
    ServiceName string
    StartDate string
    EndDate string
}
