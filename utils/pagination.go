package utils

import (
	"net/http"
	"strconv"
)

func GetPagination(r *http.Request) (int64, int64, error) {
	offset, limit := r.URL.Query().Get("offset"), r.URL.Query().Get("limit")

	if offset == "" {
		offset = "0"
	}

	if limit == "" {
		limit = "10"
	}

	offsetInt, err := strconv.ParseInt(offset, 10, 64)

	if err != nil {
		return 0, 0, err
	}

	limitInt, err := strconv.ParseInt(limit, 10, 64)

	if err != nil {
		return 0, 0, err
	}

	return offsetInt, limitInt, nil
}
