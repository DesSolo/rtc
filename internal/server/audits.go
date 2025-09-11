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

	queryAction := queryOr(r, "action", "config_updated")
	actor := queryOr(r, "actor", "")
	fromDate := queryOr(r, "from", now.Add(-24*time.Hour))
	toDate := queryOr(r, "to", now)

	action := models.ConvertAuditActionToModel(queryAction)
	if action == models.AuditActionUnknown {
		respondError(ctx, w, http.StatusBadRequest, "invalid action")
		return
	}

	audits, err := s.provider.AuditsSearch(ctx, models.AuditFilter{
		Action:   action,
		Actor:    actor,
		FromDate: fromDate,
		ToDate:   toDate,
	})
	if err != nil {
		slog.ErrorContext(ctx, "provider.Audits", "err", err, "action", action)
		respondError(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}

	respondData(ctx, w, http.StatusOK, listAuditsResponse{
		Audits: convertModelsToAudits(audits),
	})
}
