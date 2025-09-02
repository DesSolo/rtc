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

	queryAction := queryOr(r, "action", "config_updated")
	limit := queryOr(r, "limit", 100)
	offset := queryOr(r, "offset", 0)

	action := models.ConvertAuditActionToModel(queryAction)
	if action == models.AuditActionUnknown {
		respondError(ctx, w, http.StatusBadRequest, "invalid action")
		return
	}

	audits, err := s.provider.Audits(ctx, action, limit, offset)
	if err != nil {
		slog.ErrorContext(ctx, "provider.Audits", "err", err, "action", action)
		respondError(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}

	respondData(ctx, w, http.StatusOK, listAuditsResponse{
		Audits: convertModelsToAudits(audits),
	})

}
