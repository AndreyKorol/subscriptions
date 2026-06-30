package services

import(
	"context"
	"errors"
	"log/slog"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/AndreyKorol/subscriptions/internal/errs"
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
        return nil, errs.Internal("")
    }
    return subs, nil
}

func (s *SubscriptionService) Show(ctx context.Context, id uint) (*models.Subscription, error) {
	s.logger.Debug("Show called", "id", id)
	sub, err := s.repo.Show(ctx, id)
	if err != nil {
		s.logger.Error("Show repo error", "id", id, "error", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("subscription not found")
		}
		return nil, errs.Internal("")
	}
	return sub, nil
}

func (s *SubscriptionService) Create(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {
    s.logger.Debug("Create called", "subscription", subscription)
    sub, err := s.repo.Create(ctx, subscription)
    if err != nil {
        s.logger.Error("Create repo error", "error", err)
        return nil, errs.Internal("")
    }
    return sub, nil
}

func (s *SubscriptionService) Update(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {
    s.logger.Debug("Update called", "id", subscription.Id, "subscription", subscription)
    sub, err := s.repo.Update(ctx, subscription)
    if err != nil {
        s.logger.Error("Update repo error", "id", subscription.Id, "error", err)
        return nil, errs.Internal("")
    }
    return sub, nil
}

func (s *SubscriptionService) Destroy(ctx context.Context, id uint) error {
	s.logger.Debug("Destroy called", "id", id)
	err := s.repo.Destroy(ctx, id)
	if err != nil {
		s.logger.Error("Destroy repo error", "id", id, "error", err)
		var appErr *errs.Error
		if errors.As(err, &appErr) {
			return appErr
		}
		return errs.Internal("")
	}
	return nil
}

func (s *SubscriptionService) Aggregate(ctx context.Context, filters models.Filter) (*models.AggSubscriptions, error) {
    s.logger.Debug("Aggregate called", "filters", filters)
    agg, err := s.repo.Aggregate(ctx, filters)
    if err != nil {
        s.logger.Error("Aggregate repo error", "error", err)
        return nil, errs.Internal("")
    }
    return agg, nil
}

