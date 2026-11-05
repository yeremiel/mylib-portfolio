package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// Log is the global logger instance
	Log *zap.Logger
)

// Config holds the logger configuration
type Config struct {
	Level       string   // debug, info, warn, error, fatal
	Environment string   // Zap encoding mode: development/production; independent of APP_ENV
	OutputPaths []string // stdout, stderr, or file paths
	File        *FileConfig
}

// FileConfig holds the file logging configuration with rotation
type FileConfig struct {
	Enabled    bool   // Enable file logging
	Path       string // Log file path
	MaxSize    int    // Maximum size in megabytes before rotation
	MaxBackups int    // Maximum number of old log files to retain
	MaxAge     int    // Maximum number of days to retain old log files
	Compress   bool   // Whether to compress rotated files
}

// Init initializes the zap logger
func Init(cfg Config) error {
	// Parse log level
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	// Configure encoder based on environment
	var encoderConfig zapcore.EncoderConfig
	var encoder zapcore.Encoder

	if cfg.Environment == "production" {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Build output writers
	var writers []zapcore.WriteSyncer

	// Add console output if specified
	for _, path := range cfg.OutputPaths {
		switch path {
		case "stdout":
			writers = append(writers, zapcore.AddSync(os.Stdout))
		case "stderr":
			writers = append(writers, zapcore.AddSync(os.Stderr))
		default:
			// For file paths without rotation, use simple file
			if cfg.File == nil || !cfg.File.Enabled || path != cfg.File.Path {
				file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return err
				}
				writers = append(writers, zapcore.AddSync(file))
			}
		}
	}

	// Add file output with rotation if enabled
	if cfg.File != nil && cfg.File.Enabled {
		lj := &lumberjack.Logger{
			Filename:   cfg.File.Path,
			MaxSize:    cfg.File.MaxSize,
			MaxBackups: cfg.File.MaxBackups,
			MaxAge:     cfg.File.MaxAge,
			Compress:   cfg.File.Compress,
		}
		writers = append(writers, zapcore.AddSync(lj))
	}

	// Default to stdout if no writers configured
	if len(writers) == 0 {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}

	// Create multi-writer
	multiWriter := zapcore.NewMultiWriteSyncer(writers...)

	// Create core
	core := zapcore.NewCore(encoder, multiWriter, level)

	// Build logger
	Log = zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return nil
}

// InitDefault initializes the logger with default settings
func InitDefault() error {
	return Init(Config{
		Level:       "info",
		Environment: "development",
		OutputPaths: []string{"stdout"},
	})
}

// InitWithWriter initializes the logger with a custom writer (for testing)
func InitWithWriter(w io.Writer, level string) error {
	lvl, err := zapcore.ParseLevel(level)
	if err != nil {
		lvl = zapcore.InfoLevel
	}

	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, zapcore.AddSync(w), lvl)
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

// Sync flushes any buffered log entries
func Sync() error {
	if Log != nil {
		return Log.Sync()
	}
	return nil
}

// Info logs at Info level
func Info(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Info(msg, fields...)
	}
}

// Debug logs at Debug level
func Debug(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Debug(msg, fields...)
	}
}

// Warn logs at Warn level
func Warn(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Warn(msg, fields...)
	}
}

// Error logs at Error level
func Error(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Error(msg, fields...)
	}
}

// Fatal logs at Fatal level and exits the program
func Fatal(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Fatal(msg, fields...)
	} else {
		// fallback to os.Exit if logger is not initialized
		os.Exit(1)
	}
}

// Base returns a *zap.Logger to inject into other components and call directly.
// Log skips one caller frame because it is normally reached through the package
// functions (Info, etc.); injecting it as-is would report the wrong call site.
func Base() *zap.Logger {
	if Log == nil {
		return zap.NewNop()
	}
	return Log.WithOptions(zap.AddCallerSkip(-1))
}

// With returns a new logger with additional fields
func With(fields ...zap.Field) *zap.Logger {
	if Log != nil {
		return Log.With(fields...)
	}
	return zap.NewNop()
}

// GetLevel returns the current log level
func GetLevel() string {
	if Log != nil {
		return Log.Level().String()
	}
	return "info"
}
