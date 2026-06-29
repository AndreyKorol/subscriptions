package controllers

import(
    "errors"
    "github.com/go-playground/validator/v10"
    "github.com/AndreyKorol/subscriptions/internal/models"
)

type CreateSubscriptionRequest struct {
    ServiceName *string `json:"service_name" validate:"required,ne="`
    Price       *uint   `json:"price"        validate:"required,gte=0"`
    UserId      *string `json:"user_id"      validate:"required,uuid4"`
    StartDate   *string `json:"start_date"   validate:"required,datetime=01-2006"`
}

func (r *CreateSubscriptionRequest) Validate() error {
    return validator.New().Struct(r)
}

func (r *CreateSubscriptionRequest) ToModel() (*models.Subscription, error) {
    if err := r.Validate(); err != nil {
        return nil, errors.Join(errors.New("invalid request for conversion"), err)
    }
    return &models.Subscription{
        ServiceName: *r.ServiceName,
        Price:       *r.Price,
        UserId:      *r.UserId,
        StartDate:   *r.StartDate,
    }, nil
}

type UpdateSubscriptionRequest struct {
    ServiceName *string `json:"service_name" validate:"omitempty,ne="`
    Price       *uint   `json:"price"        validate:"omitempty,gte=0"`
    UserId      *string `json:"user_id"      validate:"omitempty,uuid4"`
    StartDate   *string `json:"start_date"   validate:"omitempty,datetime=01-2006"`
}

func (r *UpdateSubscriptionRequest) Validate() error {
    return validator.New().Struct(r)
}

func (r *UpdateSubscriptionRequest) ToModel() (*models.Subscription, error) {
    if err := r.Validate(); err != nil {
        return nil, errors.Join(errors.New("invalid request for conversion"), err)
    }
    return &models.Subscription{
        ServiceName: *r.ServiceName,
        Price:       *r.Price,
        UserId:      *r.UserId,
        StartDate:   *r.StartDate,
    }, nil
}
