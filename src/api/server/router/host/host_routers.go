package host

import (
	"net/http"
	"api/server/httputils"
	"src/golang.org/x/net/context"
)

func (s *containerRouter) getHostJSON(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {

	host := "111111111"

	return httputils.WriteJSON(w, http.StatusOK, host)
}