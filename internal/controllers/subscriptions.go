package controllers

import(
    "net/http"
    "log/slog"
    "context"
    "encoding/json"
    "github.com/google/uuid"
    "github.com/AndreyKorol/subscriptions/internal/models"
    "github.com/AndreyKorol/subscriptions/internal/services"
)

type SubscriptionsController struct {
    ctx      context.Context
    services *services.Manager
    logger   *slog.Logger
}

func NewSubscriptionsController(ctx context.Context, services *services.Manager, logger *slog.Logger) *SubscriptionsController {
    return &SubscriptionsController{
        ctx: ctx,
        services: services,
        logger: logger,
    }
}

func (c *SubscriptionsController) Query(w http.ResponseWriter, r *http.Request) {
    values := r.URL.Query()

    uuid, err := uuid.Parse(values.Get("user_id"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    filters := models.Filter{
        UserId: uuid,
        ServiceName: values.Get("service_name"),
        StartDate: values.Get("start_date"),
        EndDate: values.Get("end_date"),
    }

    subs, err := c.services.SubService.Query(r.Context(), filters)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    resp := CollectionData{Data: Items{Items: subs}}
    json.NewEncoder(w).Encode(resp)
}

// func (c *SubscriptionsController) Show(w http.ResponseWriter, r *http.Request) {
    
// }

func (c *SubscriptionsController) Create(w http.ResponseWriter, r *http.Request) {
    var subscription models.Subscription

    err := json.NewDecoder(r.Body).Decode(&subscription)
    if err != nil {
        http.Error(w, "Invalid JSON body", http.StatusBadRequest)
        return
    }

    sub, err := c.services.SubService.Create(r.Context(), &subscription)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    resp := Data{Data: sub}
    json.NewEncoder(w).Encode(resp)
}

// func (c *SubscriptionsController) Update(w http.ResponseWriter, r *http.Request) {
    
// }

// func (c *SubscriptionsController) Destroy(w http.ResponseWriter, r *http.Request) {
    
// }

// func (c *SubscriptionsController) Aggregate(w http.ResponseWriter, r *http.Request) {
    
// }

type Data struct {
    Data any `json:"data"`
}

type CollectionData struct {
    Data Items `json:"data"`
}

type Items struct {
    Items []*models.Subscription `json:"items"`
}
