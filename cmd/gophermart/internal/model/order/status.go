package order

const (
	STATUS_NEW        = "new"
	STATUS_PROCESSING = "processing"
	STATUS_INVALID    = "invalid"
	STATUS_PROCESSED  = "processed"
)

var statusValueForView = map[string]string{
	STATUS_NEW:        "NEW",
	STATUS_PROCESSING: "PROCESSING",
	STATUS_INVALID:    "INVALID",
	STATUS_PROCESSED:  "PROCESSED",
}

type Status struct {
	Value string
}

func NewStatus(value string) Status {
	return Status{Value: value}
}

func (m Status) GetValueForView() string {
	return statusValueForView[m.Value]
}
