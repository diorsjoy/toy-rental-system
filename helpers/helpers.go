package helpers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type envelope map[string]any

type Helpers interface {
	readIdParam(r *http.Request) (int64, error)
	writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error
	readJSON(w http.ResponseWriter, r *http.Request, dst any) error
}

func readIdParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, nil
	}

	return id, nil
}
