package utils

import "log/slog"

// ErrLog adds the err param and changes it to a slog attr
// Example: slog.Warn("Some warning", ErrLog(err))
func ErrLog(err error) slog.Attr {
	return slog.String("err", err.Error())
}
