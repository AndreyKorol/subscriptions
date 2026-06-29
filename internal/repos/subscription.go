package repos

import(
    "context"
    "github.com/AndreyKorol/subscriptions/internal/models"
    "github.com/huandu/go-sqlbuilder"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jackc/pgx/v5"
    "time"
)

type SubscriptionRepo struct {
    pool *pgxpool.Pool
}

func NewSubscriptionRepo(pool *pgxpool.Pool) *SubscriptionRepo {
    return &SubscriptionRepo{pool: pool}
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

    rows, err := r.pool.Query(ctx, sql, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    subs, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.Subscription])
    if err != nil {
        return nil, err
    }
    return subs, nil
}

func (r *SubscriptionRepo) Show(ctx context.Context, id uint) (*models.Subscription, error) {
    rows, err := r.pool.Query(
        ctx,
        `SELECT id,
                service_name,
                price,
                user_id,
                TO_CHAR(start_date, 'MM-YYYY') AS start_date
        FROM subscriptions
        WHERE subscriptions.id = $1`,
        id,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    sub, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.Subscription])
    if err != nil {
        return nil, err
    }
    return sub, nil
}

func (r *SubscriptionRepo) Create(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {
    rows, err := r.pool.Query(
        ctx,
        `INSERT INTO subscriptions (service_name, price, user_id, start_date)
         VALUES($1, $2, $3, TO_DATE($4, 'MM-YYYY'))
         RETURNING id, service_name, price, user_id, TO_CHAR(start_date, 'MM-YYYY') AS start_date;`,
        subscription.ServiceName, subscription.Price, subscription.UserId, subscription.StartDate,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    sub, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.Subscription])
    if err != nil {
        return nil, err
    }
    return sub, nil
}

// func (r *SubscriptionRepo) Update(ctx context.Context, id uint, subscription *models.Subscription) (*models.Subscription, error) {

// }

// func (r *SubscriptionRepo) Destroy(ctx context.Context, id uint) error {

// }

// func (r *SubscriptionRepo) Aggregate(ctx context.Context, filters models.Filter) (*models.AggSubscriptions, error) {

// }
