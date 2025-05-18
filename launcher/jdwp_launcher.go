package launcher

import (
	"os"
	"os/exec"
)

type JdwpLauncher struct {
	jarPath  string
	jdwpPort int
	cmd      *exec.Cmd
}

func NewJdwpLauncher(jarPath string, jdwpPort int) *JdwpLauncher {
	return &JdwpLauncher{jarPath: jarPath, jdwpPort: jdwpPort}
}

func (l *JdwpLauncher) Start() error {
	l.cmd = exec.Command(
		"java",
		"-agentlib:jdwp=transport=dt_socket,server=y,suspend=y,address=5005",
		"-jar", l.jarPath,
	)
	l.cmd.Stdout = os.Stdout
	l.cmd.Stderr = os.Stderr
	if err := l.cmd.Start(); err != nil {
		return err
	}
	return nil
}

func (l *JdwpLauncher) Stop() error {
	if l.cmd != nil && l.cmd.Process != nil {
		err := l.cmd.Process.Kill()
		l.cmd.Wait()
		return err
	}
	return nil
}
