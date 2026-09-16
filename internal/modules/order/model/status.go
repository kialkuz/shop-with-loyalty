package model

import "fmt"

const (
	StatusNew        = "new"
	StatusProcessing = "processing"
	StatusInvalid    = "invalid"
	StatusProcessed  = "processed"
)

var statuses = []string{
	StatusNew,
	StatusProcessing,
	StatusInvalid,
	StatusProcessed,
}

var statusValueForView = map[string]string{
	StatusNew:        "NEW",
	StatusProcessing: "PROCESSING",
	StatusInvalid:    "INVALID",
	StatusProcessed:  "PROCESSED",
}

type Status struct {
	Value string
}

func NewStatus(value string) (*Status, error) {
	status := Status{Value: value}

	if !status.isValid() {
		return nil, fmt.Errorf("invalid status: %q", value)
	}

	return &status, nil
}

func CreateStatusFromView(inputValue string) (*Status, error) {
	for value, viewValue := range statusValueForView {
		if viewValue == inputValue {
			return &Status{Value: value}, nil
		}
	}

	return nil, fmt.Errorf("invalid status view value: %q", inputValue)
}

func (s Status) IsValid() bool {
	_, ok := statusValueForView[s.Value]
	return ok
}

func (s Status) GetValueForView() string {
	return statusValueForView[s.Value]
}

func (s Status) isValid() bool {
	for _, status := range statuses {
		if s.Value == status {
			return true
		}
	}
	return false
}
