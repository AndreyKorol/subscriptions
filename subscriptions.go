package main

import(
    "net/http"
    "log/slog"
    "context"
    "time"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/AndreyKorol/subscriptions/internal/config"
    "github.com/AndreyKorol/subscriptions/internal/services"
    "github.com/AndreyKorol/subscriptions/internal/controllers"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()

    logger := slog.Default()

    cfg := config.Load()

    pool, err := pgxpool.New(ctx, cfg.DSN())
    if err != nil {
        logger.Error("connection to database failed", "error", err)
        return
    }

    services := services.NewManager(pool)
    subController := controllers.NewSubscriptionsController(ctx, services, logger)

    mux := http.NewServeMux()
    mux.HandleFunc("GET /subscriptions", subController.Index)
    mux.HandleFunc("POST /subscriptions", subController.Create)
    mux.HandleFunc("GET /subscriptions/{id}", subController.Show)
    mux.HandleFunc("PATCH /subscriptions/{id}", subController.Update)
    // mux.HandleFunc("DELETE /subscriptions/{id}", subController.Destroy)
    // mux.HandleFunc("GET /subscriptions/agg", subController.Aggregate)

    s := http.Server{
        Addr: cfg.Addr(),
        Handler: mux,
    }

    if err := s.ListenAndServe(); err != nil {
        slog.Error("server failed", "error", err)
        return
    }
}
