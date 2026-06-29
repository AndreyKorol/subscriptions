package services

import(
    "log/slog"
    "github.com/jackc/pgx/v5/pgxpool"
)

type Manager struct {
    SubService *SubscriptionService
}

func NewManager(pool *pgxpool.Pool, logger *slog.Logger) *Manager {
    return &Manager{
        SubService: NewSubscriptionService(pool, logger),
    }
}
