package helper

import (
	"encoding/json"
	"fmt"
)

type Status string

const (
	StatusSuccess      Status = "Success"
	StatusFailure      Status = "Failure"
	StatusNotSupported Status = "Not Supported"
)

type Response struct{
	Status  Status `json:"status"`
	Message string `json:"message"`
	Device  string `json:"device,omitempty"`

	Capabilities *Capabilities `json:"capabilities,omitempty"`
}

type Capabilities struct {
	Attach bool `json:"attach"`
}

func Handle(resp Response) {
	// format the output as JSON
	output, _ := json.Marshal(resp)
	fmt.Println(string(output))
}