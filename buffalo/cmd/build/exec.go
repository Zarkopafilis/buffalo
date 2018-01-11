package build

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/gobuffalo/envy"
	"github.com/sirupsen/logrus"
)

type debugWriter int

func (debugWriter) Write(data []byte) (int, error) {
	for _, l := range bytes.Split(data, []byte("\n")) {
		logrus.Debug(string(l))
	}
	return len(data), nil
}

func (b *Builder) exec(name string, args ...string) error {
	cmd := exec.CommandContext(b.ctx, name, args...)
	logrus.Debugf("running %s", strings.Join(cmd.Args, " "))
	
	cmd.Env = append(cmd.Env, "GO_ENV=" + envy.Get("GO_ENV", "development"))
	
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = debugWriter(0)
	err := cmd.Run()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
