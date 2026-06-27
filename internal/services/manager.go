package services

import(
    "github.com/jackc/pgx/v5/pgxpool"
)

type Manager struct {
    SubService *SubscriptionService
}

func NewManager(pool *pgxpool.Pool) *Manager {
    return &Manager{
        SubService: NewSubscriptionService(pool),
    }
}
