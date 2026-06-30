package repos

import(
    "time"
	"context"
	"log/slog"
	"github.com/AndreyKorol/subscriptions/internal/errs"
	"github.com/AndreyKorol/subscriptions/internal/models"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
)

type SubscriptionRepo struct {
    pool   *pgxpool.Pool
    logger *slog.Logger
}

func NewSubscriptionRepo(pool *pgxpool.Pool, logger *slog.Logger) *SubscriptionRepo {
    return &SubscriptionRepo{
        pool:   pool,
        logger: logger,
    }
}

func (r *SubscriptionRepo) Index(ctx context.Context, filters models.Filter) ([]*models.Subscription, error) {
    sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
    sb.Select("id", "service_name", "price", "user_id", `TO_CHAR(start_date, 'MM-YYYY') AS start_date`)
    sb.From("subscriptions")

    if filters.UserId != nil {
        sb.Where(sb.Equal("user_id", *filters.UserId))
    }
    if filters.ServiceName != nil {
        sb.Where(sb.Equal("service_name", *filters.ServiceName))
    }
    if filters.StartDate != nil {
        t, _ := time.Parse("01-2006", *filters.StartDate)
        sb.Where(sb.GE("start_date", t))
    }
    if filters.EndDate != nil {
        t, _ := time.Parse("01-2006", *filters.EndDate)
        sb.Where(sb.LE("start_date", t))
    }

    sql, args := sb.Build()

    r.logger.Debug("Index query", "sql", sql, "args", args)

    rows, err := r.pool.Query(ctx, sql, args...)
    if err != nil {
        r.logger.Error("Index query failed", "error", err)
        return nil, err
    }
    defer rows.Close()

    subs, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.Subscription])
    if err != nil {
        r.logger.Error("Index collect rows failed", "error", err)
        return nil, err
    }
    return subs, nil
}

func (r *SubscriptionRepo) Show(ctx context.Context, id uint) (*models.Subscription, error) {
    sql := `SELECT id,
                   service_name,
                   price,
                   user_id,
                   TO_CHAR(start_date, 'MM-YYYY') AS start_date
            FROM subscriptions
            WHERE subscriptions.id = $1
            ORDER BY subscriptions.id`

    r.logger.Debug("Show query", "sql", sql, "args", []any{id})

    rows, err := r.pool.Query(ctx, sql, id)
    if err != nil {
        r.logger.Error("Show query failed", "id", id, "error", err)
        return nil, err
    }
    defer rows.Close()

    sub, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.Subscription])
    if err != nil {
        r.logger.Error("Show collect row failed", "id", id, "error", err)
        return nil, err
    }
    return sub, nil
}

func (r *SubscriptionRepo) Create(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {
    sql := `INSERT INTO subscriptions (service_name, price, user_id, start_date)
            VALUES($1, $2, $3, TO_DATE($4, 'MM-YYYY'))
            RETURNING id, service_name, price, user_id, TO_CHAR(start_date, 'MM-YYYY') AS start_date;`

    args := []any{subscription.ServiceName, subscription.Price, subscription.UserId, subscription.StartDate}
    r.logger.Debug("Create query", "sql", sql, "args", args)

    rows, err := r.pool.Query(ctx, sql, args...)
    if err != nil {
        r.logger.Error("Create query failed", "error", err)
        return nil, err
    }
    defer rows.Close()

    sub, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.Subscription])
    if err != nil {
        r.logger.Error("Create collect row failed", "error", err)
        return nil, err
    }
    return sub, nil
}

func (r *SubscriptionRepo) Update(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {
    sql := `UPDATE subscriptions
            SET service_name = $1, price = $2, user_id = $3, start_date = TO_DATE($4, 'MM-YYYY')
            WHERE subscriptions.id = $5
            RETURNING id, service_name, price, user_id, TO_CHAR(start_date, 'MM-YYYY') AS start_date;`

    args := []any{subscription.ServiceName, subscription.Price, subscription.UserId, subscription.StartDate, subscription.Id}
    r.logger.Debug("Update query", "sql", sql, "args", args)

    rows, err := r.pool.Query(ctx, sql, args...)
    if err != nil {
        r.logger.Error("Update query failed", "id", subscription.Id, "error", err)
        return nil, err
    }
    defer rows.Close()

    sub, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.Subscription])
    if err != nil {
        r.logger.Error("Update collect row failed", "id", subscription.Id, "error", err)
        return nil, err
    }
    return sub, nil
}

func (r *SubscriptionRepo) Destroy(ctx context.Context, id uint) error {
    sql := `DELETE FROM subscriptions
            WHERE subscriptions.id = $1;`

    r.logger.Debug("Destroy query", "sql", sql, "args", []any{id})

    result, err := r.pool.Exec(ctx, sql, id)
    if err != nil {
        r.logger.Error("Destroy query failed", "id", id, "error", err)
        return err
    }
	if result.RowsAffected() == 0 {
		r.logger.Warn("Destroy no rows affected", "id", id)
		return errs.NotFound("subscription not found")
	}
	return nil
}

func (r *SubscriptionRepo) Aggregate(ctx context.Context, filters models.Filter) (*models.AggSubscriptions, error) {
    sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
    sb.Select("COALESCE(SUM(price)::bigint, 0) AS sum_price")
    sb.From("subscriptions")

    if filters.UserId != nil {
        sb.Where(sb.Equal("user_id", *filters.UserId))
    }
    if filters.ServiceName != nil {
        sb.Where(sb.Equal("service_name", *filters.ServiceName))
    }
    if filters.StartDate != nil {
        t, _ := time.Parse("01-2006", *filters.StartDate)
        sb.Where(sb.GE("start_date", t))
    }
    if filters.EndDate != nil {
        t, _ := time.Parse("01-2006", *filters.EndDate)
        sb.Where(sb.LE("start_date", t))
    }

    sb.Limit(1)

    sql, args := sb.Build()

    r.logger.Debug("Aggregate query", "sql", sql, "args", args)

    rows, err := r.pool.Query(ctx, sql, args...)
    if err != nil {
        r.logger.Error("Aggregate query failed", "error", err)
        return nil, err
    }
    defer rows.Close()

    aggSubs, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.AggSubscriptions])
    if err != nil {
        r.logger.Error("Aggregate collect row failed", "error", err)
        return nil, err
    }
    return aggSubs, nil
}
