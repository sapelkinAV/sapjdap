package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"sapelkinav/javadap/utils"
	"time"

	"github.com/rs/zerolog"
)

type JavaLauncher struct {
	jarPath  string
	jdwpPort int
	cmd      *exec.Cmd
	logger   zerolog.Logger
}

func NewJavaLauncher(jarPath string, jdwpPort int) *JavaLauncher {
	logger, err := utils.GetComponentLogger("launcher", "java")
	if err != nil {
		// Fallback to global logger if component logger can't be created
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	return &JavaLauncher{
		jarPath:  jarPath,
		jdwpPort: jdwpPort,
		logger:   logger,
	}
}

func (l *JavaLauncher) Start() error {
	// Create logs directory for Java process output
	javaLogDir := "./.logs"
	if err := os.MkdirAll(javaLogDir, 0755); err != nil {
		return fmt.Errorf("failed to create Java log directory: %w", err)
	}

	// Create Java stdout and stderr log files
	stdoutLog, err := os.OpenFile(fmt.Sprintf("%s/java_stdout.log", javaLogDir), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to create Java stdout log file: %w", err)
	}

	stderrLog, err := os.OpenFile(fmt.Sprintf("%s/java_stderr.log", javaLogDir), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		stdoutLog.Close()
		return fmt.Errorf("failed to create Java stderr log file: %w", err)
	}

	l.logger.Info().
		Str("jar", l.jarPath).
		Int("jdwpPort", l.jdwpPort).
		Msg("Starting Java application with JDWP enabled")

	// Create command with JDWP options
	l.cmd = exec.Command(
		"java",
		fmt.Sprintf("-agentlib:jdwp=transport=dt_socket,server=y,suspend=y,address=%d", l.jdwpPort),
		"-jar", l.jarPath,
	)

	// Redirect stdout and stderr to log files
	l.cmd.Stdout = stdoutLog
	l.cmd.Stderr = stderrLog

	// Start the process
	if err := l.cmd.Start(); err != nil {
		stdoutLog.Close()
		stderrLog.Close()
		return utils.LogError(l.logger, err, "Failed to start Java process")
	}

	l.logger.Info().Int("pid", l.cmd.Process.Pid).Msg("Java process started successfully")

	// Allow time for JDWP to initialize
	time.Sleep(500 * time.Millisecond)

	return nil
}

func (l *JavaLauncher) Stop() error {
	if l.cmd == nil || l.cmd.Process == nil {
		l.logger.Info().Msg("No Java process to stop")
		return nil
	}

	l.logger.Info().Int("pid", l.cmd.Process.Pid).Msg("Stopping Java process")

	// Attempt graceful termination first
	err := l.cmd.Process.Signal(os.Interrupt)
	if err != nil {
		l.logger.Warn().Err(err).Msg("Failed to send interrupt signal, attempting to kill")
		err = l.cmd.Process.Kill()
	}

	// Wait for the process to exit
	_, waitErr := l.cmd.Process.Wait()
	if waitErr != nil {
		l.logger.Error().Err(waitErr).Msg("Error waiting for Java process to exit")
	}

	l.logger.Info().Msg("Java process terminated")
	return err
}

// IsRunning checks if the Java process is still running
func (l *JavaLauncher) IsRunning() bool {
	if l.cmd == nil || l.cmd.Process == nil {
		return false
	}

	// Try to signal the process with signal 0, which doesn't actually send a signal
	// but checks if the process exists and we have permission to signal it
	err := l.cmd.Process.Signal(os.Signal(nil))
	return err == nil
}
