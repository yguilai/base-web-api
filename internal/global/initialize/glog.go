package initialize

import (
	"context"
	"fmt"
	"github.com/wam-lab/base-web-api/internal/global"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"log"
	"os"
	"time"
)

var (
	Default = NewGormLogger(log.New(os.Stdout, "\r\n", log.LstdFlags), GormConfig{
		SlowThreshold: 2 * time.Second,
		LogLevel:      logger.Warn,
	})
)

type GormConfig struct {
	SlowThreshold time.Duration
	LogLevel      logger.LogLevel
}

type Writer interface {
	Printf(string, ...interface{})
}

type GormLogger struct {
	GormConfig
	Writer
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func NewGormLogger(writer Writer, config GormConfig) logger.Interface {
	var (
		infoStr      = "[INFO] "
		warnStr      = "[WARN] "
		errStr       = "[ERROR] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	return &GormLogger{
		Writer:       writer,
		GormConfig:   config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *g
	newLogger.LogLevel = level
	return &newLogger
}

func (g *GormLogger) Info(ctx context.Context, message string, data ...interface{}) {
	if g.LogLevel > logger.Info {
		g.Printf(
			g.infoStr+message,
			append([]interface{}{utils.FileWithLineNum()}, data...)...,
		)
	}
}

func (g *GormLogger) Warn(ctx context.Context, message string, data ...interface{}) {
	if g.LogLevel > logger.Warn {
		g.Printf(
			g.warnStr+message,
			append([]interface{}{utils.FileWithLineNum()}, data...)...,
		)
	}
}

func (g *GormLogger) Error(ctx context.Context, message string, data ...interface{}) {
	if g.LogLevel >= logger.Error {
		g.Printf(
			g.errStr+message,
			append([]interface{}{utils.FileWithLineNum()}, data...)...,
		)
	}
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if g.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && g.LogLevel >= logger.Error:
			sql, rows := fc()
			if rows == -1 {
				g.Printf(g.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				g.Printf(g.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.LogLevel >= logger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", g.SlowThreshold)
			if rows == -1 {
				g.Printf(g.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				g.Printf(g.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case g.LogLevel >= logger.Info:
			sql, rows := fc()
			if rows == -1 {
				g.Printf(g.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				g.Printf(g.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		}
	}
}

func (g *GormLogger) Printf(message string, data ...interface{}) {
	if len(data) == 2 {
		global.Log.Error(
			"[GORM] "+message,
			zap.Any("src", data[0]),
			zap.Any("error", data[1]),
		)
		return
	}
	global.Log.Info(
		"[GORM] "+message,
		zap.String("type", "sql"),
		zap.Any("src", data[0]),
		zap.Any("duration", data[1]),
		zap.Any("rows", data[2]),
		zap.Any("sql", data[3]),
	)
}
