package errors

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func Handler(endpoint func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		err := endpoint(w, r)

		if err == nil {
			return
		}

		log.Error(err)

		clientError, ok := err.(ClientError)

		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		body, err := clientError.ResponseBody()
		if err != nil {
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		status, headers := clientError.ResponseHeaders()
		for k, v := range headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(status)
		w.Write(body)
	}
}
