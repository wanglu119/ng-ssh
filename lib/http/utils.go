package http

import (
	"encoding/json"
	"net/http"
	"fmt"
)

func renderJSON(w http.ResponseWriter, r *http.Request, data interface{}) (int, error) {
	marsh, err := json.Marshal(data)

	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write(marsh); err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return http.StatusInternalServerError, err
	}

	return 0, nil
}




