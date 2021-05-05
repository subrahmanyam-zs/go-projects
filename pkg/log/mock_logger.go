package log

import (
	"io"
	"sync"
)

func NewMockLogger(output io.Writer) Logger {
	rls.level = Debug

	return &logger{
		out: output,
		app: appInfo{
			Data:      make(map[string]interface{}),
			Framework: "gofr-" + GofrVersion,
			syncData:  &sync.Map{},
		},
	}
}
