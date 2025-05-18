package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ModuleLogger represents a logger for a specific module
type ModuleLogger struct {
	Logger   zerolog.Logger
	FileDesc *os.File // Keep file descriptor to close it later
}

var (
	// Map to store module loggers
	moduleLoggers = make(map[string]*ModuleLogger)
	loggersMutex  sync.RWMutex // Mutex to protect concurrent map access

	// Global configuration
	logDir      string
	globalLevel zerolog.Level
)

// InitializeLogger sets up the global logger with appropriate configuration
func InitializeLogger(baseLogDir string, logLevel string) error {
	// Store global configuration
	logDir = baseLogDir

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(baseLogDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create global log files
	currentTime := time.Now().Format("2006-01-02")
	allLogsFile := filepath.Join(baseLogDir, fmt.Sprintf("all_%s.log", currentTime))
	errorLogsFile := filepath.Join(baseLogDir, fmt.Sprintf("error_%s.log", currentTime))

	// Open global log files
	allLogs, err := os.OpenFile(allLogsFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open all logs file: %w", err)
	}

	errorLogs, err := os.OpenFile(errorLogsFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		allLogs.Close()
		return fmt.Errorf("failed to open error logs file: %w", err)
	}

	// Set global settings for zerolog
	zerolog.TimeFieldFormat = time.RFC3339

	// Parse the log level
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	globalLevel = level
	zerolog.SetGlobalLevel(level)

	// Create console writer for stdout with colors
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.CallerFieldName,
			zerolog.MessageFieldName,
		},
	}

	// Create a multi-writer for both console and file
	multi := zerolog.MultiLevelWriter(consoleWriter, allLogs)

	// Create base logger with timestamp, caller info, and configured outputs
	log.Logger = zerolog.New(multi).
		With().
		Timestamp().
		Caller().
		Logger()

	// Create a hook for error logs
	errorHook := zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, message string) {
		if level >= zerolog.ErrorLevel {
			// Get caller info to include in error log
			_, file, line, ok := runtime.Caller(4) // Adjust the stack depth as needed
			if ok {
				errorLogger := zerolog.New(errorLogs).
					With().
					Timestamp().
					Str("file", filepath.Base(file)).
					Int("line", line).
					Logger()

				// Create an event with the appropriate level and then log the message
				errorLogger.WithLevel(level).Msg(message)
			}
		}
	})

	// Add hook to the global logger
	log.Logger = log.Hook(errorHook)

	// Initialize module loggers that you know will be used
	// This is optional - they can also be initialized on first use
	if _, err := GetModuleLogger("jdwp"); err != nil {
		return err
	}
	if _, err := GetModuleLogger("dap"); err != nil {
		return err
	}
	if _, err := GetModuleLogger("project"); err != nil {
		return err
	}

	return nil
}

// GetModuleLogger creates or retrieves a logger specific to a module
func GetModuleLogger(moduleName string) (zerolog.Logger, error) {
	loggersMutex.RLock()
	moduleLogger, exists := moduleLoggers[moduleName]
	loggersMutex.RUnlock()

	if exists {
		return moduleLogger.Logger, nil
	}

	// If doesn't exist, create it with write lock
	loggersMutex.Lock()
	defer loggersMutex.Unlock()

	// Double-check after acquiring write lock
	if moduleLogger, exists = moduleLoggers[moduleName]; exists {
		return moduleLogger.Logger, nil
	}

	// Create module-specific log file with prefix instead of subdirectory
	currentTime := time.Now().Format("2006-01-02")
	// Use prefix for the filename instead of subdirectory
	moduleLogFile := filepath.Join(logDir, fmt.Sprintf("%s_%s.log", moduleName, currentTime))

	// Open module log file
	file, err := os.OpenFile(moduleLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return zerolog.Logger{}, fmt.Errorf("failed to open module log file: %w", err)
	}

	// Create console writer
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.CallerFieldName,
			zerolog.MessageFieldName,
		},
	}

	// Create multi-writer (console + module file)
	multi := zerolog.MultiLevelWriter(consoleWriter, file)

	// Create module logger
	logger := zerolog.New(multi).
		With().
		Timestamp().
		Caller().
		Str("module", moduleName).
		Logger()

	// Store in map
	moduleLoggers[moduleName] = &ModuleLogger{
		Logger:   logger,
		FileDesc: file,
	}

	return logger, nil
}

// GetComponentLogger creates a logger for a specific component within a module
func GetComponentLogger(moduleName, componentName string) (zerolog.Logger, error) {
	moduleLogger, err := GetModuleLogger(moduleName)
	if err != nil {
		return zerolog.Logger{}, err
	}

	return moduleLogger.With().Str("component", componentName).Logger(), nil
}

// Cleanup closes all open log files
func Cleanup() {
	loggersMutex.Lock()
	defer loggersMutex.Unlock()

	for moduleName, moduleLogger := range moduleLoggers {
		if moduleLogger.FileDesc != nil {
			moduleLogger.FileDesc.Close()
			delete(moduleLoggers, moduleName)
		}
	}
}
