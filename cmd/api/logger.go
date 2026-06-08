package main

import (
	"log/slog"
)

type Logger interface {
	SetHandler(handler string)
	SetMethod(method string)
	Info(message string, extras ...any)
	Error(message string, extras ...any)
}

type logger struct {
	log     *slog.Logger
	handler *string
	method  *string
}

func NewLogger(log *slog.Logger) Logger {
	return &logger{
		log: log,
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

func (l *logger) SetHandler(handler string) {
	l.handler = &handler
}

func (l *logger) SetMethod(method string) {
	l.method = &method
}

func (l *logger) setContent(extras ...any) []any {

	if l.log == nil || len(extras) == 0 {
		return []any{}
	}

	attrs := make([]any, 0, len(extras)-1+4)
	if l.handler != nil {
		attrs = append(attrs, "handler", *l.handler)
	}
	if l.method != nil {
		attrs = append(attrs, "method", *l.method)
	}

	for _, msg := range extras {
		attrs = append(attrs, msg)
	}

	return attrs
}
