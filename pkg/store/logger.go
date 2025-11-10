package store

import (
	"context"
	"time"

	"gorm.io/gorm/logger"
	"k8s.io/klog/v2"
)

// KlogLogger 实现 gorm.Logger 接口
type KlogLogger struct {
	LogLevel logger.LogLevel
}

// LogMode 设置日志级别
func (l KlogLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := l
	newLogger.LogLevel = level
	return newLogger
}

// Info 记录 info 级别的日志
func (l KlogLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info && klog.V(4).Enabled() {
		klog.Infof(msg, data...)
	}
}

// Warn 记录 warn 级别的日志
func (l KlogLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn && klog.V(1).Enabled() {
		klog.Warningf(msg, data...)
	}
}

// Error 记录 error 级别的日志
func (l KlogLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error && klog.V(1).Enabled() {
		klog.Errorf(msg, data...)
	}
}

// Trace 记录 trace 级别的日志
func (l KlogLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		// 错误日志输出在 v1 级别
		if klog.V(1).Enabled() {
			klog.Errorf("Trace Error: %v | SQL: %s | Rows affected: %d | Time: %s", err, sql, rows, elapsed)
		}
	} else {
		// 正常 SQL 日志输出在 v4 级别
		if klog.V(4).Enabled() {
			klog.Infof("Trace Success | SQL: %s | Rows affected: %d | Time: %s", sql, rows, elapsed)
		}
	}
}
