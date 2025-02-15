package logger

import "context"

type Logger interface {
	InfoContext(ctx context.Context, msg string, fields ...any)
	ErrorContext(ctx context.Context, msg string, fields ...any)
	DebugContext(ctx context.Context, msg string, fields ...any)
	WarnContext(ctx context.Context, msg string, fields ...any)
}
