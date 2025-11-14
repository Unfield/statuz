package logger

import (
	"go.uber.org/zap"
)

type ProdLogger struct {
	l *zap.Logger
}

func NewProdLogger() *ProdLogger {
	zl, _ := zap.NewProduction()
	return &ProdLogger{l: zl}
}

func (p *ProdLogger) Info(msg string, kv ...any) {
	p.l.Sugar().Infow(msg, kv...)
}
func (p *ProdLogger) Warn(msg string, kv ...any) {
	p.l.Sugar().Warnw(msg, kv...)
}
func (p *ProdLogger) Error(msg string, kv ...any) {
	p.l.Sugar().Errorw(msg, kv...)
}
func (p *ProdLogger) Debug(msg string, kv ...any) {
	p.l.Sugar().Debugw(msg, kv...)
}

func (p *ProdLogger) Infof(template string, args ...any) {
	p.l.Sugar().Infof(template, args...)
}
func (p *ProdLogger) Warnf(template string, args ...any) {
	p.l.Sugar().Warnf(template, args...)
}
func (p *ProdLogger) Errorf(template string, args ...any) {
	p.l.Sugar().Errorf(template, args...)
}
func (p *ProdLogger) Debugf(template string, args ...any) {
	p.l.Sugar().Debugf(template, args...)
}
