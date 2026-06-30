package controllers

import(
    "errors"
    "context"
    "strconv"
    "net/http"
    "log/slog"
    "encoding/json"
    "github.com/gorilla/schema"
    "github.com/go-playground/validator/v10"
    "github.com/AndreyKorol/subscriptions/internal/errs"
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
        respondError(w, errs.BadRequest("request validation failed", errs.FromValidator(err)))
        return
    }

    subs, err := c.services.SubService.Index(r.Context(), filters)
    if err != nil {
        c.logger.Error("Index service error", "error", err)
        respondError(w, err)
        return
    }

    resp := CollectionData{Data: Items{Items: subs}}
    respondJSON(w, http.StatusOK, resp)

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
        respondError(w, errs.BadRequest("invalid id"))
        return
    }
    if err := validator.New().Var(id, "gte=1"); err != nil {
        c.logger.Error("id validation failed", "id", id, "error", err)
        respondError(w, errs.BadRequest("request validation failed", errs.FromValidator(err)))
        return
    }

    sub, err := c.services.SubService.Show(r.Context(), uint(id))
    if err != nil {
        c.logger.Error("Show service error", "id", id, "error", err)
        respondError(w, err)
        return
    }

    resp := Data{Data: sub}
    respondJSON(w, http.StatusOK, resp)

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
        respondError(w, errs.BadRequest("invalid JSON body"))
        return
    }
    c.logger.Debug("request decoded", "body", createSubRequest)

    if err = validator.New().Struct(createSubRequest); err != nil {
        c.logger.Error("validation failed", "error", err)
        respondError(w, errs.BadRequest("request validation failed", errs.FromValidator(err)))
        return
    }

    subscription, err := createSubRequest.ToModel()
    if err != nil {
        c.logger.Error("ToModel conversion failed", "error", err)
        respondError(w, errs.BadRequest(err.Error()))
        return
    }

    sub, err := c.services.SubService.Create(r.Context(), subscription)
    if err != nil {
        c.logger.Error("Create service error", "error", err)
        respondError(w, err)
        return
    }

    resp := Data{Data: sub}
    respondJSON(w, http.StatusOK, resp)

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
        respondError(w, errs.BadRequest("invalid id"))
        return
    }
    if err := validator.New().Var(id, "gte=1"); err != nil {
        c.logger.Error("id validation failed", "id", id, "error", err)
        respondError(w, errs.BadRequest("request validation failed", errs.FromValidator(err)))
        return
    }

    updateSubRequest := &UpdateSubscriptionRequest{}
    err = json.NewDecoder(r.Body).Decode(updateSubRequest)
    if err != nil {
        c.logger.Error("failed to decode request body", "error", err)
        respondError(w, errs.BadRequest("invalid JSON body"))
        return
    }
    c.logger.Debug("request decoded", "body", updateSubRequest)

    if err = validator.New().Struct(updateSubRequest); err != nil {
        c.logger.Error("validation failed", "error", err)
        respondError(w, errs.BadRequest("request validation failed", errs.FromValidator(err)))
        return
    }

    existingSub, err := c.services.SubService.Show(r.Context(), uint(id))
    if err != nil {
        c.logger.Error("Show service error (pre-update)", "id", id, "error", err)
        respondError(w, err)
        return
    }

    updateSubRequest.ApplyTo(existingSub)

    sub, err := c.services.SubService.Update(r.Context(), existingSub)
    if err != nil {
        c.logger.Error("Update service error", "id", id, "error", err)
        respondError(w, err)
        return
    }

    resp := Data{Data: sub}
    respondJSON(w, http.StatusOK, resp)

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
        respondError(w, errs.BadRequest("invalid id"))
        return
    }
    if err = validator.New().Var(id, "gte=1"); err != nil {
        c.logger.Error("id validation failed", "id", id, "error", err)
        respondError(w, errs.BadRequest("request validation failed", errs.FromValidator(err)))
        return
    }

    err = c.services.SubService.Destroy(r.Context(), uint(id))
    if err != nil {
        c.logger.Error("Destroy service error", "id", id, "error", err)
        respondError(w, err)
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
        respondError(w, errs.BadRequest("request validation failed", errs.FromValidator(err)))
        return
    }

    aggSubs, err := c.services.SubService.Aggregate(r.Context(), filters)
    if err != nil {
        c.logger.Error("Aggregate service error", "error", err)
        respondError(w, err)
        return
    }

    resp := Data{Data: aggSubs}
    respondJSON(w, http.StatusOK, resp)

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

type errorResponse struct {
	Error *errs.Error `json:"error"`
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, err error) {
	if error, ok := errors.AsType[*errs.Error](err); ok {
		respondJSON(w, error.Status(), errorResponse{Error: error})
		return
	}
	respondJSON(w, http.StatusInternalServerError, errorResponse{
		Error: errs.Internal("internal server error"),
	})
}
