# Nexa Auto: Fine-Tuning Made Simple

**Nexa Auto** is a developer-first tool for fine-tuning Hugging Face-compatible language models. It provides a secure, session-aware workflow that abstracts away infrastructure complexity, letting you focus on model and data selection while ensuring your secrets are never written to disk.

---

## What is Nexa Auto?

Nexa Auto is a CLI and TUI (Terminal User Interface) application that orchestrates the entire fine-tuning process for LLMs (Large Language Models) using Hugging Face libraries.
It is designed for both beginners and advanced users who want a repeatable, portable, and secure way to fine-tune models on their own hardware or, in the future, on remote/cloud resources.

---

## Key Features

- **Session-aware token management:** Your Hugging Face token is stored securely in memory, never on disk.
- **Interactive CLI/TUI:** Choose between a modern TUI (Go/BubbleTea) or a sleek CLI (Python Rich) for orchestration.
- **Guided workflow:** Step-by-step selection of model, dataset, and output configuration.
- **Hardware detection:** Automatically detects and displays available CPU/GPU resources.
- **LoRA/PEFT support:** Efficient adapter-based fine-tuning out of the box.
- **Extensible:** Modular design for adding new training modes, hardware checks, and logging.

---

## How Nexa Auto Works

### 1. Session Server (`session_server.py`)
- Runs locally using FastAPI.
- Manage your Hugging Face token in memory, encrypted with AES-GCM.
- CLI/TUI communicates with this server to fetch, set, or clear your token.

### 2. User Interface (CLI/TUI)
- **Python CLI (`cli.py`):** Uses Python Rich for a command-line experience.
- **Go TUI (`go_cli/main.go`):** Uses BubbleTea for a modern terminal UI.
- On the first run, prompts for your Hugging Face token and stores it securely for the session.
- Guides you through:
  - Selecting a base model (from Hugging Face Hub or local path)
  - Choosing a dataset (from Hugging Face Hub or local file)
  - Naming your output model
  - Confirming hardware resources
- Supports switching between TUI and CLI modes.

### 3. Training Backend (`trainer_server.py`)
- Exposes a REST API for launching and monitoring training jobs.
- Loads the selected model and dataset, applies tokenization, and configures the Hugging Face Trainer.
- Supports LoRA/PEFT for efficient fine-tuning.
- Saves outputs ready for Hugging Face Hub upload.

---

## Quick Start Guide

### 1. Prerequisites

- Python 3.8+ (for backend and CLI)
- Go 1.18+ (for TUI, optional)
- CUDA-capable GPU (recommended for local training)
- Valid Hugging Face access token
- Dataset available on Hugging Face Hub or as a local file

### 2. Installation

Clone the repository:

```sh
git clone https://github.com/your-org/nexa-auto.git
cd nexa-auto
```

Install Python dependencies:

```sh
pip install -r requirements.txt
```

(Optional) Install Go dependencies for TUI:

```sh
cd go_cli
go mod tidy
```

### 3. Start the Session Server

```sh
python session_server.py
```

This will start a local FastAPI server to securely manage your Hugging Face token.

### 4. Launch the Interface

#### Python CLI

```sh
python cli.py
```

#### Go TUI

```sh
cd go_cli
go run main.go
```

### 5. Follow the Prompts

- Enter your Hugging Face token (prompted on first run)
- Select your base model and dataset
- Name your output model
- Confirm detected hardware (CPU/GPU)
- Start training

### 6. Monitor Progress

- View logs and training progress in the interface.
- Retrieve output artifacts for upload to Hugging Face Hub.

---

## Example Workflow

1. **Start the session server:**  
   `python session_server.py`
2. **Run the CLI or TUI:**  
   `python cli.py` or `go run main.go`
3. **Authenticate:**  
   Enter your Hugging Face token when prompted.
4. **Configure training:**  
   Select model, dataset, and output name.
5. **Review hardware:**  
   Confirm detected resources.
6. **Start and monitor training:**  
   Watch logs and progress in real time.
7. **Retrieve outputs:**  
   Find your fine-tuned model ready for upload.

---

## Security Model

- **No secrets on disk:** Tokens are only stored in memory, never written to disk.
- **Local-only server:** The session server only listens to localhost.
- **Clear on exit:** Tokens are cleared at session end or on user request.

---

## Extending Nexa Auto

- **Add new training modes:** Implement in `remote.py` and update the UI menu.
- **Add hardware checks:** Extend `hardware.py`.
- **Add logging/metrics:** Hook into `logging.py`.
- **Modular design:** All components are ready for extension.

---

## Roadmap

- [x] Initial release with CLI and TUI
- [x] Session-aware token management
- [x] Local training with Hugging Face Trainer
- [x] LoRA/PEFT support
- [x] Basic hardware detection
- [x] Remote training support (SSH, API)
- [x] Remote/cloud training (Kaggle, SSH, etc.)
- [x] Advanced logging and metrics
- [ ] Model card auto-generation

---

## Troubleshooting & Support

- If you encounter issues, check the logs in the interface.
- For feature requests or bug reports, open an issue or PR on [GitHub](https://github.com/your-org/nexa-auto).

---

## FAQ

**Q: Is my Hugging Face token safe?**  
A: Yes.
It is only stored in memory, encrypted, and never written to disk.

**Q: Can I use my own dataset?**  
A: Yes. You can select a local file or a dataset from the Hugging Face Hub.

**Q: Can I run this on a remote server?**  
A: Remote/SSH/API modes are planned for future releases.

---

**Nexa Auto** — Secure, repeatable, and developer-friendly fine-tuning for LLMs.
