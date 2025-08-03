package manager

import (
	"encoding/json"
	"fmt"
	"golambda/models"
	"os"
	"os/exec"
)

func RegisterFunction(name string, trigger_event string, code []byte) error {
	sourcePath := fmt.Sprintf("user_functions/%s.go", name)
	os.WriteFile(sourcePath, code, 0644)

	execPath := fmt.Sprintf("plugins/%s.exe", name)
	cmd := exec.Command("go", "build", "-o", execPath, sourcePath)
	err := cmd.Run()
	if err != nil {
		fmt.Print(err.Error())
		return err
	}

	metadata := models.FunctionMetadata{
		Name:    name,
		Trigger: trigger_event,
		Path:    execPath,
	}
	jsonData, err := json.MarshalIndent(metadata, "", "	")
	if err != nil {
		return err
	}
	err = os.WriteFile("registry/registry.json", jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}
