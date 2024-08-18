package bybit_provider

import "errors"

var (
	ErrEmptyResponse       = errors.New("empty response")
	ErrInvalidResponseType = errors.New("invalid response type")
)
