package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ExtensionName = filepath.Base(os.Args[0]) // extension name has to match the filename
	PrintPrefix   = fmt.Sprintf("[%s] ", ExtensionName)
)

// Method for pretty printing objects in logs
func PrettyPrint(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return ""
	}
	return string(data)
}
