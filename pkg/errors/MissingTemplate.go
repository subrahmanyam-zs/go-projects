package errors

import (
	"fmt"
	"strings"
)

type MissingTemplate struct {
	FileLocation string
	FileName     string
}

func (e MissingTemplate) Error() string {
	if strings.TrimSpace(e.FileName) == "" {
		return "Filename not provided"
	}

	return fmt.Sprintf("Template %v does not exist at location: %v", e.FileName, e.FileLocation)
}
