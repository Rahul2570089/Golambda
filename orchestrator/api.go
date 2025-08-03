package orchestrator

import (
	"encoding/json"
	"fmt"
	"golambda/manager"
	"net/http"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func RegisterHttpRoute(name string) {
	route := fmt.Sprintf("/%s", name)

	log.WithField("route", route).Info("Registering HTTP route")

	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(logrus.Fields{
			"route":  route,
			"method": r.Method,
		}).Info("Incoming HTTP request")

		output, err := manager.ExecuteBin(name)
		if err != nil {
			log.WithFields(logrus.Fields{
				"function": name,
				"error":    err,
			}).Error("Error executing function")

			http.Error(w, fmt.Sprintf("Unable to execute binary: %s", err), http.StatusNotImplemented)
			return
		}

		log.WithFields(logrus.Fields{
			"function": name,
			"output":   string(output),
		}).Info("Function executed successfully")

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": string(output)})
	})
}
