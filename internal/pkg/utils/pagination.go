// Package utils holds small, dependency-free helpers shared across services.
package utils

const (
	defaultOffset = 0
	defaultLimit  = 20
	maxLimit      = 100
)

// ParsePage validates offset/limit query params and falls back to sane
// defaults. Use it in handlers before calling the service layer, e.g.:
//
//	offset, limit := utils.ParsePage(rawOffset, rawLimit)
func ParsePage(offset, limit int64) (int64, int64) {
	if offset < 0 {
		offset = defaultOffset
	}
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	return offset, limit
}
