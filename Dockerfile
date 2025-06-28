# ---- Python/uv base ----
FROM python:3.10-slim AS python-base

# System deps
RUN apt-get update && apt-get install -y git gcc curl && rm -rf /var/lib/apt/lists/*

# Install uv (fast Python package manager)
RUN pip install --upgrade pip && pip install uv

# Set workdir
WORKDIR /app

# Copy Python project files
COPY pyproject.toml uv.lock ./
COPY requirements.txt ./

# Install Python deps with uv
RUN uv pip install --system --no-cache-dir

# Copy Python source
COPY . .

# ---- Go base ----
FROM golang:1.21-alpine AS go-base

WORKDIR /go_cli
COPY go_cli/go.mod go_cli/go.sum ./
RUN go mod download

COPY go_cli/ ./
RUN go build -o /usr/local/bin/nexa-tui main.go splash.go

# ---- Final image ----
FROM python:3.10-slim

WORKDIR /app

# System deps
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

# Copy Python env from python-base
COPY --from=python-base /usr/local/lib/python3.10 /usr/local/lib/python3.10
COPY --from=python-base /usr/local/bin /usr/local/bin
COPY --from=python-base /app /app

# Copy Go TUI binary
COPY --from=go-base /usr/local/bin/nexa-tui /usr/local/bin/nexa-tui

# Expose FastAPI/uvicorn port
EXPOSE 8000

# Entrypoint: run uvicorn server (session_server.py) by default
CMD ["uvicorn", "session_server:app", "--host", "0.0.0.0", "--port", "8000"]
