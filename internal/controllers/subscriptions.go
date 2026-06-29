package controllers

import(
    "net/http"
    "log/slog"
    "context"
    "strconv"
    "encoding/json"
    "github.com/gorilla/schema"
    "github.com/go-playground/validator/v10"
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

func (c *SubscriptionsController) Index(w http.ResponseWriter, r *http.Request) {
    filters := models.Filter{}
    schema.NewDecoder().Decode(&filters, r.URL.Query())

    if err := validator.New().Struct(filters); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    subs, err := c.services.SubService.Index(r.Context(), filters)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    resp := CollectionData{Data: Items{Items: subs}}
    json.NewEncoder(w).Encode(resp)
}

func (c *SubscriptionsController) Show(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }
    if err := validator.New().Var(id, "gte=1"); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    sub, err := c.services.SubService.Show(r.Context(), uint(id))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if sub == nil {
        http.Error(w, err.Error(), http.StatusNotFound)
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    resp := Data{Data: sub}
    json.NewEncoder(w).Encode(resp)
}

func (c *SubscriptionsController) Create(w http.ResponseWriter, r *http.Request) {
    var createSubRequest CreateSubscriptionRequest

    err := json.NewDecoder(r.Body).Decode(createSubRequest)
    if err != nil {
        http.Error(w, "Invalid JSON body", http.StatusBadRequest)
        return
    }

    subscription, err := createSubRequest.ToModel()
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    sub, err := c.services.SubService.Create(r.Context(), subscription)
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
