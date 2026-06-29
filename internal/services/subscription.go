package services

import(
    "context"
    "log/slog"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/AndreyKorol/subscriptions/internal/repos"
    "github.com/AndreyKorol/subscriptions/internal/models"
)

type SubscriptionService struct {
    repo   *repos.SubscriptionRepo
    logger *slog.Logger
}

func NewSubscriptionService(pool *pgxpool.Pool, logger *slog.Logger) *SubscriptionService {
    return &SubscriptionService{
        repo:   repos.NewSubscriptionRepo(pool, logger),
        logger: logger,
    }
}

func (s *SubscriptionService) Index(ctx context.Context, filters models.Filter) ([]*models.Subscription, error) {
    s.logger.Debug("Index called", "filters", filters)
    subs, err := s.repo.Index(ctx, filters)
    if err != nil {
        s.logger.Error("Index repo error", "error", err)
        return nil, err
    }
    return subs, nil
}

func (s *SubscriptionService) Show(ctx context.Context, id uint) (*models.Subscription, error) {
    s.logger.Debug("Show called", "id", id)
    sub, err := s.repo.Show(ctx, id)
    if err != nil {
        s.logger.Error("Show repo error", "id", id, "error", err)
        return nil, err
    }
    return sub, nil
}

func (s *SubscriptionService) Create(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {
    s.logger.Debug("Create called", "subscription", subscription)
    sub, err := s.repo.Create(ctx, subscription)
    if err != nil {
        s.logger.Error("Create repo error", "error", err)
        return nil, err
    }
    return sub, nil
}

func (s *SubscriptionService) Update(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {
    s.logger.Debug("Update called", "id", subscription.Id, "subscription", subscription)
    sub, err := s.repo.Update(ctx, subscription)
    if err != nil {
        s.logger.Error("Update repo error", "id", subscription.Id, "error", err)
        return nil, err
    }
    return sub, nil
}

func (s *SubscriptionService) Destroy(ctx context.Context, id uint) error {
    s.logger.Debug("Destroy called", "id", id)
    err := s.repo.Destroy(ctx, id)
    if err != nil {
        s.logger.Error("Destroy repo error", "id", id, "error", err)
        return err
    }
    return nil
}

func (s *SubscriptionService) Aggregate(ctx context.Context, filters models.Filter) (*models.AggSubscriptions, error) {
    s.logger.Debug("Aggregate called", "filters", filters)
    agg, err := s.repo.Aggregate(ctx, filters)
    if err != nil {
        s.logger.Error("Aggregate repo error", "error", err)
        return nil, err
    }
    return agg, nil
}

