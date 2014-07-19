package vagrant

import (
	"fmt"
	"os"

	"github.com/mlafeldt/chef-runner/exec"
)

const (
	DefaultMachine = "default"
)

func init() {
	os.Setenv("VAGRANT_NO_PLUGINS", "1")
}

type Client struct {
	Machine string
}

func NewClient(machine string) *Client {
	if machine == "" {
		machine = DefaultMachine
	}
	return &Client{Machine: machine}
}

func (c Client) String() string {
	return fmt.Sprintf("Vagrant (machine: %s)", c.Machine)
}

func (c Client) SSHCommand(command string) []string {
	return []string{"vagrant", "ssh", c.Machine, "-c", command}
}

func (c Client) RunCommand(command string) error {
	return exec.RunCommand(c.SSHCommand(command))
}