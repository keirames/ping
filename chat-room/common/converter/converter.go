package converter

import (
	"net/http"
	"strconv"
)

func GetParam(r *http.Request, field string) string {
	return r.URL.Query().Get(field)
}

func StringToInt(p string) (int, error) {
	return strconv.Atoi(p)
}

func StringToInt64(p string) (int64, error) {
	return strconv.ParseInt(p, 10, 64)
}
