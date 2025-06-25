# Nexa Auto

**Nexa Auto** is a session-aware, developer-centric CLI orchestration engine for fine-tuning Hugging Face-compatible models. It abstracts away ML infrastructure pain while preserving power-user control, making repeatable, portable, and secure LLM fine-tuning as simple as a single command.

---

## 🚀 Features

- **CLI-first, infra-aware:** Rich terminal UI, hardware detection, and resource validation.
- **Session-aware token management:** Secure, in-memory Hugging Face token handling via a local session server.
- **Minimal config:** Only model, dataset, and output name are required; everything else is auto-inferred.
- **Local and remote training:** Train on your local GPU, or (future) launch jobs on remote APIs or SSH nodes.
- **LoRA/PEFT support:** Out-of-the-box support for LoRA adapters and quantized training.
- **Reproducible artifacts:** All outputs are ready for Hugging Face Hub upload.
- **No notebooks required:** Run and manage experiments entirely from the CLI.

---

## 🧑‍💻 Use Cases

- **Domain adaptation:** Fine-tune open LLMs (e.g., Mistral, Llama) on your own scientific, legal, or business datasets.
- **Research workflows:** Run repeatable, isolated experiments without Jupyter or cloud lock-in.
- **Infra abstraction:** Seamlessly switch between local, SSH, or (future) cloud API training.
- **Secure collaboration:** Share models and configs without leaking tokens or credentials.

---

## ⚡ Quickstart

```bash
# Install dependencies
pip install -r requirements.txt

# Start the CLI
python cli.py

# Example: Fine-tune a model
# (Follow the interactive prompts for model, dataset, and output name)
```

---

## 🛠️ Assumptions

- You have Python 3.8+ and a CUDA-capable GPU (for local training).
- You have a valid Hugging Face access token.
- Your dataset is either on the Hugging Face Hub or in a local JSON/text format.
- For remote/API/SSH modes, you have network access and credentials (future).

---

## 🏗️ Project Structure

```
nexa_auto/
├── cli.py              # Main CLI entrypoint and orchestration
├── session_server.py   # Local Flask server for secure token storage
├── config.py           # Config class for model/dataset/params (future)
├── trainer.py          # Training logic (future)
├── hardware.py         # Hardware detection (future)
├── remote.py           # Remote/SSH/API logic (future)
├── logging.py          # Logging utilities (future)
├── requirements.txt
├── README.md
└── doc.md
```

---

## 🧩 Key Design Principles

- **Stateful, repeatable, isolated:** Each session is secure and reproducible.
- **CLI-first, notebook-free:** No Jupyter required.
- **Minimal config, maximal power:** Only specify what matters.
- **Portable and hackable:** Open, extensible, and not cloud-locked.

---

## 📦 Outputs

- `adapter_model.safetensors`
- `adapter_config.json`
- `tokenizer.json` and config
- `training_args.bin`
- Logs and metrics

---

## 🔒 Security

- Hugging Face tokens are never written to disk; they're stored encrypted in memory via the session server.
- Tokens are cleared at session end or on demand.

---

## 🧠 Why Nexa Auto?

- **Open-source alternative** to cloud-locked tools like AutoTrain.
- **Transparent and extensible** for research and production.
- **Scales with you**: from local dev to remote clusters.

---

## 📝 License

MIT License. See [LICENSE](LICENSE).

---

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Please open an issue or PR.

---

## 📚 Documentation

See [doc.md](./doc.md) for a technical overview and developer notes.
