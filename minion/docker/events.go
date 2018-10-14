package docker

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"
)

type Events struct {
	Status string `json:"status,omitempty"`
	ID     string `json:"id,omitempty"`
	From   string `json:"from,omitempty"`
	Type   string `json:"Type,omitempty"`
	Action string `json:"Action,omitempty"`
	Actor  Actor  `json:"Actor,omitempty"`
}

type Actor struct {
	ID         string     `json:ID,omitempty`
	Attributes Attributes `json:"Attributes,omitempty"`
}

type Attributes struct {
	Image string `json:"image,omitempty"`
	Name  string `json:"name,omitempty"`
}

func DockerEvents(ctx context.Context, event chan<- []byte) error {
	var (
		cmd     *exec.Cmd
		err     error
		stdout  io.ReadCloser
		scanner *bufio.Scanner
	)

	cmd = exec.CommandContext(ctx, "docker", "events", "--format", "{{json .}}")
	cmd.Stderr = os.Stderr

	if stdout, err = cmd.StdoutPipe(); err != nil {
		return err
	}

	scanner = bufio.NewScanner(bufio.NewReader(stdout))

	go func() {
		for scanner.Scan() {
			event <- scanner.Bytes()
		}
	}()

	if err = cmd.Start(); err != nil {
		return err
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	cmd.Wait()

	return nil
}
