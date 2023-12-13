package common

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/rs/zerolog"
)

func NewLogger(output string) (*zerolog.Logger, error) {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return fmt.Sprintf("%s:L%d", runtime.FuncForPC(pc).Name(), line)
	}

	if _, err := os.Stat(output); err != nil {
		saveDir := filepath.Dir(output)
		os.MkdirAll(saveDir, 0755)
	}

	f, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	logs := zerolog.New(f).With().Timestamp().Caller().Logger()

	return &logs, nil
}
