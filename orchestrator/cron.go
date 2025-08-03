package orchestrator

import (
	"golambda/manager"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func RegisterCron(name string, expression string) {
	log := logrus.WithFields(logrus.Fields{
		"function": name,
		"schedule": expression,
	})

	log.Info("Registering cron job")

	c := cron.New()

	_, err := c.AddFunc(expression, func() {
		log.Info("Cron job triggered")

		output, err := manager.ExecuteBin(name)
		if err != nil {
			log.WithError(err).Error("Error executing function")
			return
		}

		log.WithField("output", output).Info("Function executed successfully")
	})

	if err != nil {
		log.WithError(err).Error("Failed to schedule cron job")
		return
	}

	c.Start()
	log.Info("Cron scheduler started")
}
