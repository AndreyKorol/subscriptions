package repos

import(
    "context"
    "github.com/AndreyKorol/subscriptions/internal/models"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jackc/pgx/v5"
)

type SubscriptionRepo struct {
    pool *pgxpool.Pool
}

func NewSubscriptionRepo(pool *pgxpool.Pool) *SubscriptionRepo {
    return &SubscriptionRepo{pool: pool}
}

func (r *SubscriptionRepo) Query(ctx context.Context, filters models.Filter) ([]*models.Subscription, error) {
    rows, err := r.pool.Query(ctx, "SELECT id, service_name, price, user_id, start_date FROM subscriptions;")
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

// func (r *SubscriptionRepo) Create(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {

// }

// func (r *SubscriptionRepo) Update(ctx context.Context, id int, subscription *models.Subscription) (*models.Subscription, error) {

// }

// func (r *SubscriptionRepo) Destroy(ctx context.Context, id int) error {

// }

// func (r *SubscriptionRepo) Aggregate(ctx context.Context, filters Filter) (*models.AggSubscriptions, error) {

// }
