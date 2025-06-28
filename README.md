# Nexa Auto

**Nexa Auto** is your all-in-one, session-aware CLI and TUI tool for fine-tuning Hugging Face-compatible models. It makes secure, repeatable, and portable LLM fine-tuning as simple as following a guided workflow—no notebooks, no cloud lock-in, no headaches.

---

## Why Nexa Auto?

- **CLI & TUI:** Choose your interface—modern terminal UI (Go/BubbleTea) or a sleek CLI (Python Rich).
- **Secure by design:** Your Hugging Face token is kept in memory, never written to disk.
- **Guided workflow:** Select model, dataset, and output step-by-step.
- **Hardware smart:** Detects and displays your CPU/GPU resources.
- **LoRA/PEFT ready:** Efficient adapter-based fine-tuning out of the box.
- **Extensible:** Modular for new training modes, hardware checks, and logging.

---

## Installation

**Requirements:**
- Python 3.8+ (backend & CLI)
- Go 1.18+ (for TUI, optional)
- CUDA GPU (recommended)
- Hugging Face access token

**Get started:**
```bash
git clone https://github.com/your-org/nexa-auto.git
cd nexa-auto
pip install -r requirements.txt
# (Optional) For TUI:
cd go_cli
go mod tidy
```

---

## Quickstart

1. **Start the session server:**
   ```bash
   python session_server.py
   ```
2. **Launch your interface:**
   - **CLI:** `python cli.py`
   - **TUI:** `cd go_cli && go run main.go`
3. **Follow the prompts:**
   - Enter your Hugging Face token (first run)
   - Pick your base model and dataset
   - Name your output
   - Confirm hardware
   - Start training!
4. **Monitor progress:**  
   Watch logs and training status live in your chosen interface.

---

## How It Works

- **Session Server:** Local FastAPI server keeps your Hugging Face token safe in memory.
- **CLI/TUI:** Guides you through model/dataset/output selection and training.
- **Trainer Backend:** Handles model loading, tokenization, LoRA/PEFT, and artifact saving.

---

## Project Structure

```
nexa_auto/
├── cli.py              # Python CLI (Rich)
├── session_server.py   # Secure token server (FastAPI)
├── trainer_server.py   # Training backend (REST)
├── go_cli/             # Go TUI (BubbleTea)
├── doc.md              # Technical docs
└── README.md
```

---

## Security

- **No secrets on disk:** Tokens are only in memory, encrypted.
- **Local-only:** Session server listens only on localhost.
- **Clear on exit:** Tokens wiped at session end or on request.

---

## Example Workflow

```bash
python session_server.py
python cli.py
# or
cd go_cli && go run main.go
```
- Authenticate with your Hugging Face token.
- Select model, dataset, and output name.
- Confirm hardware.
- Start and monitor training.
- Retrieve your fine-tuned model—ready for Hugging Face Hub!

---

## Extending Nexa Auto

- Add new training modes: edit `remote.py` and update the UI.
- Add hardware checks: extend `hardware.py`.
- Add logging/metrics: hook into `logging.py`.

---

## Documentation

See [doc.md](./doc.md) for full technical details, architecture, and extension notes.

---

## Contributing

We welcome issues, feature requests, and PRs!  
Open an issue or pull request to get involved.

---

## License

MIT License. See [LICENSE](LICENSE).

---

