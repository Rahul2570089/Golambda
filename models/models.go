package models

type Function struct {
	Name    string `json:"name"`
	Trigger string `json:"trigger"`
	Code    []byte `json:"code"`
}

type FunctionMetadata struct {
	Name    string `json:"name"`
	Trigger string `json:"trigger"`
	Path    string `json:"path"`
}
