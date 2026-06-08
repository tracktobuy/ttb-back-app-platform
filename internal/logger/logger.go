package logger

import (
	"log/slog"
	"os"
)

type Logger interface {
	SetHandlerName(name string)
	SetMethodName(name string)
	SetRepositoryName(name string)
	SetServiceName(name string)
	SetMiddlewareName(name string)
	SetConfigName(name string)
	SetComponentName(name string)

	Info(message string, extras ...any)
	Error(message string, extras ...any)
}

type logger struct {
	log        *slog.Logger
	handler    *string
	service    *string
	repository *string
	middleware *string
	config     *string
	component  *string
	method     *string
}

func NewLogger() Logger {
	return &logger{
		log: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (l *logger) Info(message string, extras ...any) {
	attr := l.setContent(extras...)
	l.log.Info(message, attr...)
}

func (l *logger) Error(message string, extras ...any) {
	attr := l.setContent(extras...)
	l.log.Error(message, attr...)
}

func (l *logger) SetHandlerName(name string) {
	l.handler = &name
}

func (l *logger) SetMethodName(name string) {
	l.method = &name
}

func (l *logger) SetRepositoryName(name string) {
	l.repository = &name
}

func (l *logger) SetServiceName(name string) {
	l.service = &name
}

func (l *logger) SetMiddlewareName(name string) {
	l.middleware = &name
}

func (l *logger) SetConfigName(name string) {
	l.config = &name
}

func (l *logger) SetComponentName(name string) {
	l.component = &name
}

func (l *logger) setContent(extras ...any) []any {
	if l.log == nil {
		return []any{}
	}

	attrs := make([]any, 0, len(extras)+6)
	if l.handler != nil {
		attrs = append(attrs, "handler", *l.handler)
	}
	if l.repository != nil {
		attrs = append(attrs, "repository", *l.repository)
	}
	if l.service != nil {
		attrs = append(attrs, "service", *l.service)
	}
	if l.middleware != nil {
		attrs = append(attrs, "middleware", *l.middleware)
	}
	if l.config != nil {
		attrs = append(attrs, "config", *l.config)
	}
	if l.component != nil {
		attrs = append(attrs, "component", *l.component)
	}
	if l.method != nil {
		attrs = append(attrs, "method", *l.method)
	}

	attrs = append(attrs, extras...)
	return attrs
}
