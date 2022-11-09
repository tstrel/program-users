package api

import (
	"encoding/json"
	"io"
)

func NewIndentEncoder(w io.Writer) *json.Encoder {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc
}
