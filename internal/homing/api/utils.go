package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/mainflux/mainflux/logger"
	"github.com/mainflux/mainflux/pkg/errors"
)

// BearerPrefix represents the token prefix for Bearer authentication scheme.
const BearerPrefix = "Bearer "

// ErrorRes represents the HTTP error response body.
type ErrorRes struct {
	Err string `json:"error"`
}

// LoggingErrorEncoder is a go-kit error encoder logging decorator.
func LoggingErrorEncoder(logger logger.Logger, enc kithttp.ErrorEncoder) kithttp.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		switch err {
		case ErrBearerToken,
			ErrLimitSize,
			ErrOffsetSize:
			logger.Error(err.Error())
		}

		enc(ctx, err, w)
	}
}

// ExtractBearerToken returns value of the bearer token. If there is no bearer token - an empty value is returned.
func ExtractBearerToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if !strings.HasPrefix(token, BearerPrefix) {
		return ""
	}

	return strings.TrimPrefix(token, BearerPrefix)
}

// ReadUintQuery reads the value of uint64 http query parameters for a given key
func ReadUintQuery(r *http.Request, key string, def uint64) (uint64, error) {
	vals := bone.GetQuery(r, key)
	if len(vals) > 1 {
		return 0, errors.ErrInvalidQueryParams
	}

	if len(vals) == 0 {
		return def, nil
	}

	strval := vals[0]
	val, err := strconv.ParseUint(strval, 10, 64)
	if err != nil {
		return 0, errors.ErrInvalidQueryParams
	}

	return val, nil
}

// ReadStringQuery reads the value of string http query parameters for a given key
func ReadStringQuery(r *http.Request, key string, def string) (string, error) {
	vals := bone.GetQuery(r, key)
	if len(vals) > 1 {
		return "", errors.ErrInvalidQueryParams
	}

	if len(vals) == 0 {
		return def, nil
	}

	return vals[0], nil
}
