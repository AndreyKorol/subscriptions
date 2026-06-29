package main

import(
    "os"
    "time"
    "context"
    "net/http"
    "log/slog"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/AndreyKorol/subscriptions/internal/config"
    "github.com/AndreyKorol/subscriptions/internal/services"
    "github.com/AndreyKorol/subscriptions/internal/controllers"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()

    cfg := config.Load()

    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.AppLogLevel}))
    slog.SetDefault(logger)

    logger.Info("config loaded", "env", cfg.AppEnv, "port", cfg.AppPort)

    pool, err := pgxpool.New(ctx, cfg.DSN())
    if err != nil {
        logger.Error("connection to database failed", "error", err)
        return
    }
    logger.Info("connected to database")

    services := services.NewManager(pool, logger)
    subController := controllers.NewSubscriptionsController(ctx, services, logger)

    mux := http.NewServeMux()
    mux.HandleFunc("GET /subscriptions", subController.Index)
    mux.HandleFunc("POST /subscriptions", subController.Create)
    mux.HandleFunc("GET /subscriptions/{id}", subController.Show)
    mux.HandleFunc("PATCH /subscriptions/{id}", subController.Update)
    mux.HandleFunc("DELETE /subscriptions/{id}", subController.Destroy)
    mux.HandleFunc("GET /subscriptions/agg", subController.Aggregate)

    s := http.Server{
        Addr:    cfg.Addr(),
        Handler: controllers.LoggingMiddleware(logger, mux),
    }

    logger.Info("server starting", "addr", cfg.Addr())

    if err := s.ListenAndServe(); err != nil {
        logger.Error("server failed", "error", err)
    }
}
