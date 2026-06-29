package controllers

import(
    "errors"
    "github.com/AndreyKorol/subscriptions/internal/models"
)

type CreateSubscriptionRequest struct {
    ServiceName *string `json:"service_name" validate:"required,ne="`
    Price       *uint   `json:"price"        validate:"required,gte=0"`
    UserId      *string `json:"user_id"      validate:"required,uuid4"`
    StartDate   *string `json:"start_date"   validate:"required,datetime=01-2006"`
}

func (r *CreateSubscriptionRequest) ToModel() (*models.Subscription, error) {
    sub := &models.Subscription{}

    if r.ServiceName == nil {
        return nil, errors.New("ServiceName is missing")
    }
    if r.Price == nil {
        return nil, errors.New("Price is missing")
    }
    if r.UserId == nil {
        return nil, errors.New("UserId is missing")
    }
    if r.StartDate == nil {
        return nil, errors.New("StartDate is missing")
    }

    sub.ServiceName = *r.ServiceName
    sub.Price = *r.Price
    sub.UserId = *r.UserId
    sub.StartDate = *r.StartDate

    return sub, nil
}

type UpdateSubscriptionRequest struct {
    ServiceName *string `json:"service_name" validate:"omitempty,ne="`
    Price       *uint   `json:"price"        validate:"omitempty,gte=0"`
    UserId      *string `json:"user_id"      validate:"omitempty,uuid4"`
    StartDate   *string `json:"start_date"   validate:"omitempty,datetime=01-2006"`
}

func (r *UpdateSubscriptionRequest) ApplyTo(sub *models.Subscription) {
	if r.ServiceName != nil {
		sub.ServiceName = *r.ServiceName
	}
	if r.Price != nil {
		sub.Price = *r.Price
	}
	if r.UserId != nil {
		sub.UserId = *r.UserId
	}
	if r.StartDate != nil {
		sub.StartDate = *r.StartDate
	}
}
