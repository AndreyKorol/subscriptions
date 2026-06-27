package repos

import(
    "context"
    "github.com/AndreyKorol/subscriptions/internal/models"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jackc/pgx/v5"
    "regexp"
    "errors"
)

type SubscriptionRepo struct {
    pool *pgxpool.Pool
}

func NewSubscriptionRepo(pool *pgxpool.Pool) *SubscriptionRepo {
    return &SubscriptionRepo{pool: pool}
}

func (r *SubscriptionRepo) Query(ctx context.Context, filters models.Filter) ([]*models.Subscription, error) {
    rows, err := r.pool.Query(
        ctx,
        `SELECT id,
                service_name,
                price,
                user_id,
                TO_CHAR(start_date, 'MM-YYYY') AS start_date
        FROM subscriptions;`,
    )
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

// func (r *SubscriptionRepo) Show(ctx context.Context, id int) (*models.Subscription, error) {
    
// }

func (r *SubscriptionRepo) Create(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {
    re := regexp.MustCompile(`(\d\d)-(\d\d\d\d)`)
    match := re.FindString(subscription.StartDate)
    if match == "" {
        return nil, errors.New("invalid date format")
    }

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

// func (r *SubscriptionRepo) Update(ctx context.Context, id int, subscription *models.Subscription) (*models.Subscription, error) {

// }

// func (r *SubscriptionRepo) Destroy(ctx context.Context, id int) error {

// }

// func (r *SubscriptionRepo) Aggregate(ctx context.Context, filters Filter) (*models.AggSubscriptions, error) {

// }
