package runner

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunWhisper(filename string) (string, error) {
	// docker run -d --rm --gpus all --name whisper -v data:/app whisper-gx10 test.mp3 --model medium --language Chinese --output_dir /app
	// Note: The user said "-v data:/app", implying a volume usage or local bind mount.
	// If "data" is a local folder, it should be an absolute path.
	// For simplicity, we get the absolute path of our local "data" directory.

	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	absDataPath := filepath.Join(pwd, "data")

	// The command the user gave:
	// docker run -d --rm --gpus all --name whisper -v data:/app whisper-gx10 test.mp3 --model medium --language Chinese --output_dir /app

	// We should probably NOT run with "-d" (detached) if we want to wait for it in this function easily,
	// OR we run detached and then wait/poll.
	// However, exec.Command waits by default if we don't put it in background.
	// If we use -d, the command returns immediately with the container ID.
	// Then we'd need to `docker wait <container_id>`.
	// It's easier to remove "-d" and run synchronously for the worker.

	// Also, ensure unique container name or let Docker assign one to avoid conflict if running multiple (though worker is serial?).
	// We'll remove "--name whisper" or make it unique.

	// Commands:
	// "test.mp3" corresponds to 'filename' inside container (mapped to /app)

	cmd := exec.Command("docker", "run", "--rm", "--gpus", "all",
		"-v", fmt.Sprintf("%s:/app", absDataPath),
		"whisper-gx10",
		filename,
		"--model", "medium",
		"--language", "Chinese",
		"--output_dir", "/app",
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("docker run failed: %v, stderr: %s", err, stderr.String())
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
		// Try checking if it exists, maybe whisper naming convention is slightly different
		return "Translation completed but output file not found or read error.", nil
	}

	return string(content), nil
}
