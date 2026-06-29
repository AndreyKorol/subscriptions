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
    c.logger.Info("handling Index",
        "method", r.Method,
        "path", r.URL.Path,
        "query", r.URL.RawQuery,
    )

    filters := models.Filter{}
    schema.NewDecoder().Decode(&filters, r.URL.Query())

    if err := validator.New().Struct(filters); err != nil {
        c.logger.Error("invalid filters", "error", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    subs, err := c.services.SubService.Index(r.Context(), filters)
    if err != nil {
        c.logger.Error("Index service error", "error", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    resp := CollectionData{Data: Items{Items: subs}}
    json.NewEncoder(w).Encode(resp)

    c.logger.Info("Index completed", "count", len(subs))
}

func (c *SubscriptionsController) Show(w http.ResponseWriter, r *http.Request) {
    c.logger.Info("handling Show",
        "method", r.Method,
        "path", r.URL.Path,
        "id", r.PathValue("id"),
    )

    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.logger.Error("invalid id format", "id", idStr, "error", err)
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }
    if err := validator.New().Var(id, "gte=1"); err != nil {
        c.logger.Error("id validation failed", "id", id, "error", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    sub, err := c.services.SubService.Show(r.Context(), uint(id))
    if err != nil {
        c.logger.Error("Show service error", "id", id, "error", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if sub == nil {
        c.logger.Warn("subscription not found", "id", id)
        http.Error(w, "subscription not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    resp := Data{Data: sub}
    json.NewEncoder(w).Encode(resp)

    c.logger.Info("Show completed", "id", id)
}

func (c *SubscriptionsController) Create(w http.ResponseWriter, r *http.Request) {
    c.logger.Info("handling Create",
        "method", r.Method,
        "path", r.URL.Path,
    )

    createSubRequest := &CreateSubscriptionRequest{}

    err := json.NewDecoder(r.Body).Decode(createSubRequest)
    if err != nil {
        c.logger.Error("failed to decode request body", "error", err)
        http.Error(w, "Invalid JSON body", http.StatusBadRequest)
        return
    }
    c.logger.Debug("request decoded", "body", createSubRequest)

    if err = validator.New().Struct(createSubRequest); err != nil {
        c.logger.Error("validation failed", "error", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    subscription, err := createSubRequest.ToModel()
    if err != nil {
        c.logger.Error("ToModel conversion failed", "error", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    sub, err := c.services.SubService.Create(r.Context(), subscription)
    if err != nil {
        c.logger.Error("Create service error", "error", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    resp := Data{Data: sub}
    json.NewEncoder(w).Encode(resp)

    c.logger.Info("Create completed", "id", sub.Id)
}

func (c *SubscriptionsController) Update(w http.ResponseWriter, r *http.Request) {
    c.logger.Info("handling Update",
        "method", r.Method,
        "path", r.URL.Path,
        "id", r.PathValue("id"),
    )

    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.logger.Error("invalid id format", "id", idStr, "error", err)
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }
    if err := validator.New().Var(id, "gte=1"); err != nil {
        c.logger.Error("id validation failed", "id", id, "error", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    updateSubRequest := &UpdateSubscriptionRequest{}
    err = json.NewDecoder(r.Body).Decode(updateSubRequest)
    if err != nil {
        c.logger.Error("failed to decode request body", "error", err)
        http.Error(w, "Invalid JSON body", http.StatusBadRequest)
        return
    }
    c.logger.Debug("request decoded", "body", updateSubRequest)

    if err = validator.New().Struct(updateSubRequest); err != nil {
        c.logger.Error("validation failed", "error", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    existingSub, err := c.services.SubService.Show(r.Context(), uint(id))
    if err != nil {
        c.logger.Error("Show service error (pre-update)", "id", id, "error", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if existingSub == nil {
        c.logger.Warn("subscription not found for update", "id", id)
        http.Error(w, "subscription not found", http.StatusNotFound)
        return
    }

    updateSubRequest.ApplyTo(existingSub)

    sub, err := c.services.SubService.Update(r.Context(), existingSub)
    if err != nil {
        c.logger.Error("Update service error", "id", id, "error", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    resp := Data{Data: sub}
    json.NewEncoder(w).Encode(resp)

    c.logger.Info("Update completed", "id", id)
}

func (c *SubscriptionsController) Destroy(w http.ResponseWriter, r *http.Request) {
    c.logger.Info("handling Destroy",
        "method", r.Method,
        "path", r.URL.Path,
        "id", r.PathValue("id"),
    )

    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.logger.Error("invalid id format", "id", idStr, "error", err)
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }
    if err = validator.New().Var(id, "gte=1"); err != nil {
        c.logger.Error("id validation failed", "id", id, "error", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err = c.services.SubService.Destroy(r.Context(), uint(id))
    if err != nil {
        c.logger.Error("Destroy service error", "id", id, "error", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)

    c.logger.Info("Destroy completed", "id", id)
}

func (c *SubscriptionsController) Aggregate(w http.ResponseWriter, r *http.Request) {
    c.logger.Info("handling Aggregate",
        "method", r.Method,
        "path", r.URL.Path,
        "query", r.URL.RawQuery,
    )

    filters := models.Filter{}
    schema.NewDecoder().Decode(&filters, r.URL.Query())

    if err := validator.New().Struct(filters); err != nil {
        c.logger.Error("invalid filters", "error", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    aggSubs, err := c.services.SubService.Aggregate(r.Context(), filters)
    if err != nil {
        c.logger.Error("Aggregate service error", "error", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    resp := Data{Data: aggSubs}
    json.NewEncoder(w).Encode(resp)

    c.logger.Info("Aggregate completed", "sum_price", aggSubs.SumPrice)
}

type Data struct {
    Data any `json:"data"`
}

type CollectionData struct {
    Data Items `json:"data"`
}

type Items struct {
    Items []*models.Subscription `json:"items"`
}
