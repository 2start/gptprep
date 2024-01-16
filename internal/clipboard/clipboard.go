package clipboard

import (
	"fmt"
	"runtime"

	"github.com/atotto/clipboard"
)

func LoadToClipboard(content string) error {
	switch runtime.GOOS {
	case "darwin", "linux", "windows":
		return clipboard.WriteAll(content)
	default:
		return fmt.Errorf("unsupported platform")
	}
}
