package services

import(
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/AndreyKorol/subscriptions/internal/repos"
    "github.com/AndreyKorol/subscriptions/internal/models"
)

type SubscriptionService struct {
    repo *repos.SubscriptionRepo
}

func NewSubscriptionService(pool *pgxpool.Pool) *SubscriptionService {
    return &SubscriptionService{
        repo: repos.NewSubscriptionRepo(pool),
    }
}

func (s *SubscriptionService) Index(ctx context.Context, filters models.Filter) ([]*models.Subscription, error) {
    return s.repo.Index(ctx, filters)
}

func (s *SubscriptionService) Show(ctx context.Context, id uint) (*models.Subscription, error) {
    return s.repo.Show(ctx, id)
}

func (s *SubscriptionService) Create(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {
    return s.repo.Create(ctx, subscription)
}

func (s *SubscriptionService) Update(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {
    return s.repo.Update(ctx, subscription)
}

func (s *SubscriptionService) Destroy(ctx context.Context, id uint) error {
    return s.repo.Destroy(ctx, id)
}

// func (s *SubscriptionService) Aggregate(ctx context.Context, filters models.Filter) (*models.AggSubscriptions, error) {
//     return s.repo.Aggregate(ctx, filters)
// }

