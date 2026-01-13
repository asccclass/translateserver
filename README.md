# Whisper Translation System

This is a Golang-based server that accepts audio/video uploads, translates/transcribes them using OpenAI's Whisper (via Docker), and emails the results to the user.

## Prerequisites

- [Docker](https://www.docker.com/) installed and running.
- Docker image `whisper-gx10` available.
- NVIDIA Container Toolkit (if using `--gpus all`).

## Installation

### x86_64 (Standard PC)

1.  Clone the repository.
2.  Build:
    ```bash
    go build -o server .
    ```
3.  Run:
    ```bash
    ./server
    ```

### ARM Architecture (e.g., Raspberry Pi, Apple Silicon)

To compile for ARM64 Linux:

```bash
env GOOS=linux GOARCH=arm64 go build -o server-arm64 .
```

To compile for Apple Silicon (M1/M2/M3) macOS:

```bash
env GOOS=darwin GOARCH=arm64 go build -o server-mac-arm64 .
```

## Configuration

Set environment variables to configure:

- `CHECK_INTERVAL`: Seconds between checking for new jobs (default: 30).
- `SMTP_HOST`: SMTP server host (e.g., smtp.gmail.com).
- `SMTP_PORT`: SMTP port (e.g., 587).
- `SMTP_USER`: SMTP username.
- `SMTP_PASS`: SMTP password.

## Docker Setup

Ensure you have the image:
```bash
# Example if building locally
docker build -t whisper-gx10 .
```

## Usage

1. Open http://localhost:8080.
2. Enter email and upload a file.
3. Wait for the email with the result.
