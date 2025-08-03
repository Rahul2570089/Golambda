package main

import (
	"encoding/json"
	"fmt"
	"golambda/manager"
	"golambda/models"
	"golambda/orchestrator"
	"net/http"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	log.Info("Starting server on :8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server listening on :8080"))
		log.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Info("Health check hit")
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Info("Register function endpoint called")

		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			log.Warn("Invalid method used for /register")
			return
		}

		var function models.Function
		err := json.NewDecoder(r.Body).Decode(&function)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid payload: %s", err), http.StatusBadRequest)
			log.WithError(err).Error("Failed to decode function payload")
			return
		}

		log.WithFields(logrus.Fields{
			"name":    function.Name,
			"trigger": function.Trigger,
		}).Info("Registering function")

		err = manager.RegisterFunction(function.Name, function.Trigger, function.Code)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error registering the function: %s", err), http.StatusNotImplemented)
			log.WithError(err).Error("Failed to register function")
			return
		}

		if function.Trigger == "http" {
			orchestrator.RegisterHttpRoute(function.Name)
			log.WithField("name", function.Name).Info("HTTP route registered")
		} else if function.Trigger != "" {
			orchestrator.RegisterCron(function.Name, function.Trigger)
			log.WithFields(logrus.Fields{
				"name":     function.Name,
				"schedule": function.Trigger,
			}).Info("Cron job registered")
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "Function registered"})
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.WithError(err).Fatal("Server failed to start")
	}
}
