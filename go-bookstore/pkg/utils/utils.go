package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func ParseBody(r *http.Request, x interface{}) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal([]byte(body), x); err != nil {
		return
	}

	// decode := json.NewDecoder(r.Body)
	// err := decode.Decode(x)

	// if err != nil {
	// 	panic(err)
	// }
}
