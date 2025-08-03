package manager

import (
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func ExecuteBin(name string) ([]byte, error) {
	execPath := fmt.Sprintf("plugins/%s.exe", name)
	log.WithField("binary", execPath).Info("Executing function binary")

	cmd := exec.Command(execPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.WithFields(logrus.Fields{
			"binary": execPath,
			"error":  err,
			"output": string(output),
		}).Error("Execution failed")
		return nil, err
	}

	log.WithFields(logrus.Fields{
		"binary": execPath,
		"output": string(output),
	}).Info("Execution successful")

	return output, nil
}
