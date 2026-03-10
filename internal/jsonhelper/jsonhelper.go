package jsonhelper

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// json response
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	// packet limiter
	maxBytes := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// decoder
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON object")
	}

	return nil
}

// sends error response
func WriteJSONError(w http.ResponseWriter, status int, message string) error {
	type errorResponse struct {
		Error   string `json:"error"`
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	return WriteJSON(w, status, errorResponse{
		Error:   http.StatusText(status),
		Status:  status,
		Message: message,
	})
}
