package generatename

import (
	"fmt"
	"path/filepath"
	"time"
)

func GenerateUniqueFilename(filename string) string {
	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}
