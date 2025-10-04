package sl

import "log/slog"

// Err Функция для форматирования стндартной ошибки в Slog атрибут
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
