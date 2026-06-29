package config

import (
    "fmt"
    "os"
    "strconv"
    "log/slog"

    "github.com/joho/godotenv"
)

type Config struct {
    DBUser      string
    DBPassword  string
    DBName      string
    DBHost      string
    DBPort      int
    DBSSLMode   string
    AppPort     string
    AppEnv      string
    AppLogLevel slog.Level
}

func Load() *Config {
    _ = godotenv.Load()

    return &Config{
        DBUser:      getEnv("POSTGRES_USER", "postgres"),
        DBPassword:  getEnv("POSTGRES_PASSWORD", ""),
        DBName:      getEnv("POSTGRES_DB", "subscriptions"),
        DBHost:      getEnv("DB_HOST", "postgres"),
        DBPort:      getEnvAsInt("DB_PORT", 5432),
        DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
        AppPort:     getEnv("APP_PORT", "8080"),
        AppEnv:      getEnv("APP_ENV", "development"),
        AppLogLevel: getSlogLevel(getEnv("APP_LOG_LEVEL", "debug")),
    }
}

func (c *Config) DSN() string {
    return fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s?sslmode=%s",
        c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode,
    )
}

func (c *Config) Addr() string {
    return fmt.Sprintf(":%s", c.AppPort)
}

func getSlogLevel(level string) slog.Level {
    switch level {
    case "debug":
        return slog.LevelDebug
    case "info":
        return slog.LevelInfo
    case "warn":
        return slog.LevelWarn
    case "error":
        return slog.LevelError
    default:
        return slog.LevelDebug
    }
}

func getEnv(key, defaultVal string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
    if val := os.Getenv(key); val != "" {
        if i, err := strconv.Atoi(val); err == nil {
            return i
        }
    }
    return defaultVal
}
