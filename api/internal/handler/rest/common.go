package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type HandlerError struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

func (e HandlerError) Error() string {
	return fmt.Sprintf("Code: %d, Description:%s", e.Code, e.Description)
}

func ErrorHandler(handlerFunc func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Catch handler error
		if err := handlerFunc(w, r); err != nil {
			// Write HandlerError if err is HandlerError
			herr, ok := err.(HandlerError)
			if ok {
				w.WriteHeader(herr.Code)

				json.NewEncoder(w).Encode(HandlerError{
					Code:        herr.Code,
					Description: herr.Description,
				})
				log.Println(err.Error())

				return
			}

			// Write internal error
			w.WriteHeader(http.StatusInternalServerError)

			json.NewEncoder(w).Encode(HandlerError{
				Code:        http.StatusInternalServerError,
				Description: "Internal Server Error",
			})
			log.Println(err.Error())
		}
	}
}
