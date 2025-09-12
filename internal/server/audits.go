package server

import (
	"log/slog"
	"net/http"
	"time"

	"rtc/internal/models"
)

type audit struct {
	Action  string    `json:"action"`
	Actor   string    `json:"actor"`
	Payload []byte    `json:"payload"`
	Ts      time.Time `json:"ts"`
}

type listAuditsResponse struct {
	Audits []audit `json:"audits"`
}

func (s *Server) handleListAudits(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	now := time.Now().UTC()

	queryAction := queryOr(r, "action", "")
	actor := queryOr(r, "actor", "")
	fromDate := queryOr(r, "from", now.Add(-24*time.Hour))
	toDate := queryOr(r, "to", now)

	filter := models.AuditFilter{
		Actor:    actor,
		FromDate: fromDate,
		ToDate:   toDate,
	}

	if queryAction != "" {
		action := models.ConvertAuditActionToModel(queryAction)
		if action == models.AuditActionUnknown {
			respondError(ctx, w, http.StatusBadRequest, "invalid action")
			return
		}
		filter.Action = action
	}

	audits, err := s.provider.AuditsSearch(ctx, filter)
	if err != nil {
		slog.ErrorContext(ctx, "provider.Audits", "err", err, "filter", filter)
		respondError(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}

	respondData(ctx, w, http.StatusOK, listAuditsResponse{
		Audits: convertModelsToAudits(audits),
	})
}

func (s *Server) handleAuditActions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	actions, err := s.provider.AuditActions(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "provider.AuditActions", "err", err)
		respondError(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}

	respondData(ctx, w, http.StatusOK, convertModelsToAuditActions(actions))
}
