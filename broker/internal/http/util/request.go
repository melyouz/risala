package util

import (
	"encoding/json"
	"net/http"
)

func Decode(r *http.Request, dst interface{}) {
	_ = json.NewDecoder(r.Body).Decode(&dst)
}
