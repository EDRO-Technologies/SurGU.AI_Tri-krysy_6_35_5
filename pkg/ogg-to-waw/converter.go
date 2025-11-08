package ogg_to_waw

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

func Convert(input []byte) ([]byte, error) {
	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-ac", "1", "-ar", "16000", "-f", "wav", "pipe:1")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("cmd.StdinPipe: %w", err)
	}
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = &bytes.Buffer{}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("cmd.Start: %w", err)
	}

	go func() {
		defer stdin.Close()
		_, _ = io.Copy(stdin, bytes.NewReader(input))
	}()

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("ffmpeg error: %w, stderr: %s", err, cmd.Stderr.(*bytes.Buffer).String())
	}

	return stdout.Bytes(), nil
}
