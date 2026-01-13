package runner

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func RunWhisper(filename string) (string, error) {
	// docker run -d --rm --gpus all --name whisper -v data:/app whisper-gx10 test.mp3 --model medium --language Chinese --output_dir /app
	// Note: The user said "-v data:/app", implying a volume usage or local bind mount.
	// If "data" is a local folder, it should be an absolute path.
	// For simplicity, we get the absolute path of our local "data" directory.

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Minute)
	defer cancel()

	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	absDataPath := filepath.Join(pwd, "data")

	// The command the user gave:
	// docker run -d --rm --gpus all --name whisper -v data:/app whisper-gx10 test.mp3 --model medium --language Chinese --output_dir /app

	cmd := exec.CommandContext(ctx, "docker", "run", "--rm", "--gpus", "all",
		"-v", fmt.Sprintf("%s:/app", absDataPath),
		"whisper-gx10",
		filename,
		"--model", "medium",
		"--language", "Chinese",
		"--output_dir", "/app",
	)

	// Capture both stdout and stderr for debugging
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	fmt.Printf("Executing Docker command: %s\n", cmd.String())
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("docker run timed out after 10 minutes")
		}
		return "", fmt.Errorf("docker run failed: %v, stdout: %s, stderr: %s", err, stdout.String(), stderr.String())
	}

	// Read output file. Whisper normally produces .txt, .srt, .vtt.
	// Assuming we want the text content.
	// The input file extension might be .mp3, .m4a etc.
	// Whisper usually outputs [filename].txt

	baseName := filename
	if idx := strings.LastIndex(filename, "."); idx != -1 {
		baseName = filename[:idx]
	}

	// Let's try to read the .txt file
	outputFile := filepath.Join(absDataPath, baseName+".txt")
	content, err := ioutil.ReadFile(outputFile)
	if err != nil {
		// Log the stdout/stderr even if successful run but missing file, might help debug
		fmt.Printf("Whisper finished but output file missing. Stdout: %s\nStderr: %s\n", stdout.String(), stderr.String())
		return "Translation completed but output file not found or read error.", nil
	}

	return string(content), nil
}
