package logger

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"sutjin/go-rest-template/internal/pkg/config"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var lg *zap.Logger

func InitLogger(cfg *config.LogConfig) (err error) {
	writeSyncer := getLogWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	encoder := getEncoder()

	pe := zap.NewProductionEncoderConfig()
	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(cfg.Level))

	if err != nil {
		return
	}

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, writeSyncer, l),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), l),
	)

	// core := zapcore.NewCore(encoder, writeSyncer, l)
	lg = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(lg) //  Replace zap Global in package logger example

	return
}

func getEncoder() zapcore.Encoder {

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {

	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}

	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger  receive gin The default log of the framework
// TODO: change output name
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		requestID := requestid.Get(c)

		c.Next()

		cost := time.Since(start)

		lg.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("requestId", requestID),
			zap.String("http", c.Request.Method+" "+path+" "+query),
			zap.String("From", "SOME_SERVICE"),
			zap.String("source", c.ClientIP()+", "+c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("time", cost),
		)
	}
}

// GinRecovery recover Drop the project that may appear panic, And use zap Record relevant logs

func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.

				var brokenPipe bool

				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				if brokenPipe {
					lg.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)

					// If the connection is dead, we can't write a status to it.

					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func InlineLog(c *gin.Context, message string, data interface{}) {
	// Disable logger only for test ENV
	if lg != nil {
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		requestID := requestid.Get(c)

		lg.Error("path",
			zap.String("requestId", requestID),
			zap.String("http", c.Request.Method+" "+path+" "+query),
			zap.String("From", "IAM"),
			zap.String("errors", message),
			zap.Reflect("payload", data),
		)
	}
}
