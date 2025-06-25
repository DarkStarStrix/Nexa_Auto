# Nexa Auto: Technical Documentation

## Overview

**Nexa Auto** is a CLI-first, session-aware fine-tuning engine for Hugging Face models. It is designed to make LLM fine-tuning repeatable, secure, and portable, with minimal configuration and maximum transparency. Nexa Auto abstracts away infrastructure complexity while giving developers full control over their training workflows.

---

## How It Works

### 1. Session Server

- A lightweight Flask server (`session_server.py`) runs in the background.
- It securely stores your Hugging Face token in memory (never on disk), encrypted with Fernet.
- The CLI communicates with this server to fetch, set, or clear the token as needed.

### 2. CLI Orchestration

- The CLI (`cli.py`) provides a rich, interactive terminal interface using [Rich](https://github.com/Textualize/rich).
- On first run, it prompts for your Hugging Face token and stores it securely for the session.
- You select the base model, dataset, and output name via prompts.
- Hardware is auto-detected (GPU/CPU) and displayed for transparency.
- The CLI supports multiple modes (local, remote API, SSH), with local training implemented and others scaffolded for future extension.

### 3. Training Pipeline

- The CLI loads the selected model and dataset, applies tokenization, and configures the Hugging Face Trainer.
- LoRA/PEFT support is built-in for efficient adapter-based fine-tuning.
- Training outputs are saved in a format ready for Hugging Face Hub upload.

### 4. Security

- Tokens are never written to disk.
- Tokens are cleared at session end or on user request.
- All secrets are handled in-memory and over localhost only.

---

## Key Components

- **cli.py**: Main CLI entrypoint and orchestration logic.
- **session_server.py**: Secure, in-memory token management.
- **config.py**: (Planned) Central config for all run parameters.
- **trainer.py**: (Planned) Modular training logic, LoRA/PEFT integration.
- **hardware.py**: (Planned) Hardware detection and validation.
- **remote.py**: (Planned) Remote/SSH/API training orchestration.
- **logging.py**: (Planned) Structured logging and metrics.

---

## Example Workflow

1. Run `python cli.py`
2. Enter your Hugging Face token (stored securely for the session)
3. Select model, dataset, and output name
4. Confirm hardware and start training
5. Artifacts are saved and ready for upload

---

## Extending Nexa Auto

- Add new training modes by implementing in `remote.py` and updating the CLI menu.
- Add new hardware checks in `hardware.py`.
- Add new logging or metrics hooks in `logging.py`.
- All core logic is modular and ready for extension.

---

## Assumptions

- Python 3.8+ and a CUDA-capable GPU for local training.
- A valid Hugging Face access token.
- Dataset is on Hugging Face Hub or in a local JSON/text format.
- For remote/API/SSH modes, you have network access and credentials (future).

---

## Security Notes

- All token handling is local and in-memory.
- No secrets are printed or stored on disk.
- The session server is local-only and cleared on demand or exit.

---

## Roadmap

- [ ] Add Docker support
- [ ] Add remote/cloud training
- [ ] Add advanced logging and metrics
- [ ] Add model card auto-generation

---

## Support

For issues, feature requests, or contributions, please open an issue or PR on the [GitHub repository](https://github.com/your-org/nexa-auto).

---
````

</file>
