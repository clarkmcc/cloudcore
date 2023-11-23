package utils

import (
	"encoding/json"
	"fmt"
)

func PrintStruct(v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
