"""
Nexa Auto CLI

A session-aware, developer-centric fine-tuning orchestration engine for Hugging Face models.
Supports local, remote API, and SSH training modes with secure token management.
"""

import os
import platform
import subprocess
import sys
import time
import requests
import torch
from rich import box
from rich.align import Align
from rich.console import Console
from rich.live import Live
from rich.panel import Panel
from rich.prompt import Prompt
from rich.table import Table
from rich.text import Text
from transformers import AutoModelForCausalLM, AutoTokenizer, TrainingArguments, Trainer, AutoConfig

# === Future stubs for modularity and optimization ===
# from hardware import detect_hardware
# from remote import launch_remote_training
# from logging import setup_logging
# from optim import fast_tokenizer, fast_loader

console = Console()
SESSION_SERVER_URL = "http://127.0.0.1:8765"
VERSION = "0.3.0"

def start_session_server():
    """
    Launch the session server as a background process and wait until it's ready.
    """
    import psutil

    script_path = os.path.join(os.path.dirname(os.path.abspath(__file__)), "session_server.py")
    for proc in psutil.process_iter(['pid', 'name', 'cmdline']):
        cmdline = proc.info.get('cmdline') or []
        if 'session_server.py' in ' '.join(cmdline):
            return

    if platform.system() == "Windows":
        creationflags = subprocess.CREATE_NEW_PROCESS_GROUP
        subprocess.Popen(
            [sys.executable, script_path],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
            creationflags=creationflags
        )
    else:
        subprocess.Popen(
            [sys.executable, script_path],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
            start_new_session=True
        )

    # Wait for the server to be ready
    for _ in range(20):
        try:
            resp = requests.get(f"{SESSION_SERVER_URL}/get_token", timeout=0.5)
            if resp.status_code in (200, 404):
                return
        except requests.RequestException:
            time.sleep(0.25)
    console.print("[bold red]Session server failed to start or is not responding.[/bold red]")
    sys.exit(1)

def store_token(token):
    """
    Store the Hugging Face token securely in the session server.
    """
    requests.post(f"{SESSION_SERVER_URL}/set_token", json={"token": token})

def clear_token():
    """
    Clear the stored Hugging Face token from the session server.
    """
    requests.post(f"{SESSION_SERVER_URL}/clear_token")

def get_token():
    """
    Retrieve the Hugging Face token from the session server.
    Returns None if not set.
    """
    resp = requests.get(f"{SESSION_SERVER_URL}/get_token")
    if resp.status_code == 200:
        return resp.json()["token"]
    elif resp.status_code == 404:
        return None
    else:
        console.print(f"[bold red]Unexpected error from session server: {resp.text}[/bold red]")
        return None

def print_hf_logo():
    """
    Print the Hugging Face ASCII logo and Nexa Auto welcome message.
    """
    logo = '''
[bold yellow]██████╗ ██╗   ██╗ ██████╗  ██████╗  ██████╗ 
██╔══██╗██║   ██║██╔═══██╗██╔════╝ ██╔═══██╗
██████╔╝██║   ██║██║   ██║██║  ███╗██║   ██║
██╔═══╝ ██║   ██║██║   ██║██║   ██║██║   ██║
██║     ╚██████╔╝╚██████╔╝╚██████╔╝╚██████╔╝
╚═╝      ╚═════╝  ╚═════╝  ╚═════╝  ╚═════╝[/bold yellow]
'''
    console.print(logo)
    console.print(Panel(Text("[bold magenta]Welcome to nexa_auto[/bold magenta]", justify="center"), box=box.DOUBLE, style="bold yellow"))
    console.print(Panel("[bold cyan]A Hugging Face-aligned fine-tuning engine for models in a dev-friendly CLI.[/bold cyan]", style="bold yellow"))

def spinner_loading(message, duration=2):
    """
    Show a spinner with a message for a given duration.
    """
    import itertools
    spinner = itertools.cycle(["/", "-", "\\", "|"])
    with Live(Align.center("", vertical="middle"), refresh_per_second=10, console=console) as live:
        start = time.time()
        while time.time() - start < duration:
            live.update(Align.center(f"[bold yellow]{next(spinner)}[/bold yellow] [bold cyan]{message}[/bold cyan]", vertical="middle"))
            time.sleep(0.1)

def render_main_menu():
    """
    Render the main command menu.
    """
    table = Table(title="[bold magenta]Nexa Auto CLI[/bold magenta]", box=box.ROUNDED, style="bold yellow", show_lines=True)
    table.add_column("Command", style="bold green", justify="center")
    table.add_column("Description", style="bold white")
    table.add_row("[bold green]1[/bold green]", "[cyan]Train a model (local)[/cyan]")
    table.add_row("[bold green]2[/bold green]", "[cyan]Train on remote API (Kaggle, Lambda, etc.)[/cyan]")
    table.add_row("[bold green]3[/bold green]", "[cyan]Train via SSH on remote server[/cyan]")
    table.add_row("[bold green]4[/bold green]", "[cyan]Show hardware info[/cyan]")
    table.add_row("[bold green]5[/bold green]", "[cyan]Clear session token[/cyan]")
    table.add_row("[bold green]6[/bold green]", "[cyan]Help[/cyan]")
    table.add_row("[bold green]7[/bold green]", "[red]Exit[/red]")
    return table

def render_hardware_panel():
    """
    Render a panel with detected hardware info (GPU or CPU).
    """
    try:
        result = subprocess.run(["nvidia-smi"], capture_output=True, text=True, check=True)
        return Panel(result.stdout, title="[bold cyan]NVIDIA-SMI Output[/bold cyan]", style="bold yellow", box=box.ROUNDED)
    except (subprocess.CalledProcessError, FileNotFoundError):
        cpu = platform.processor() or platform.machine()
        return Panel(f"[bold cyan]CPU:[/bold cyan] {cpu}", title="[bold cyan]CPU Info[/bold cyan]", style="bold yellow", box=box.ROUNDED)

def run_trainer(model_name, dataset_name, new_model_name):
    """
    Run the Hugging Face Trainer for the selected model and dataset.
    """
    console.print(f"[bold magenta]Fine-tuning [yellow]{model_name}[/yellow] on [cyan]{dataset_name}[/cyan] as [green]{new_model_name}[/green]...[/bold magenta]")
    hf_token = get_token()
    if not hf_token:
        console.print("[bold red]No Hugging Face token found in session. Aborting.[/bold red]")
        return
    console.print(f"[bold cyan]Loading model and tokenizer from Hugging Face Hub...[/bold cyan]")
    model_id = f"{model_name.lower()}"
    config = AutoConfig.from_pretrained(model_id, use_auth_token=hf_token)
    model = AutoModelForCausalLM.from_pretrained(model_id, config=config, use_auth_token=hf_token)
    tokenizer = AutoTokenizer.from_pretrained(model_id, use_auth_token=hf_token)
    from datasets import load_dataset
    dataset = load_dataset(dataset_name, split="train")
    def tokenize_function(examples):
        return tokenizer(examples['text'], truncation=True, padding='max_length', max_length=128)
    tokenized_dataset = dataset.map(tokenize_function, batched=True)
    output_dir = os.path.join(os.getcwd(), new_model_name)
    training_args = TrainingArguments(
        output_dir=output_dir,
        num_train_epochs=1,
        per_device_train_batch_size=2,
        save_steps=10,
        save_total_limit=1,
        logging_steps=5,
        report_to=[],
        push_to_hub=False,
    )
    trainer = Trainer(
        model=model,
        args=training_args,
        train_dataset=tokenized_dataset,
        tokenizer=tokenizer,
    )
    trainer.train()
    trainer.save_model(output_dir)
    tokenizer.save_pretrained(output_dir)
    console.print(f"[bold green]Model and tokenizer saved to {output_dir}![/bold green]")
    del model
    del trainer
    torch.cuda.empty_cache()

def train_flow_local():
    """
    Interactive flow for local training.
    """
    models = ["mistralai/Mistral-7B-v0.1", "meta-llama/Llama-2-7b-hf"]
    model_panel = Panel(
        "\n".join([f"[yellow]{idx}.[/yellow] [cyan]{model}[/cyan]" for idx, model in enumerate(models, 1)]),
        title="[bold magenta]Available models to fine-tune[/bold magenta]",
        box=box.ROUNDED,
        style="bold yellow"
    )
    console.print(model_panel)
    model_choice = Prompt.ask("[bold yellow]Select a model by number[/bold yellow]", choices=[str(i) for i in range(1, len(models)+1)])
    selected_model = models[int(model_choice)-1]
    dataset_panel = Panel(
        "[bold yellow]Enter the name of the dataset to use (e.g. 'wikitext' or './mydata.json')[/bold yellow]",
        title="[bold magenta]Dataset Selection[/bold magenta]",
        box=box.ROUNDED,
        style="bold yellow"
    )
    console.print(dataset_panel)
    dataset = Prompt.ask("[bold yellow]Dataset[/bold yellow]")
    name_panel = Panel(
        "[bold yellow]Enter a new name for your fine-tuned model[/bold yellow]",
        title="[bold magenta]Model Name[/bold magenta]",
        box=box.ROUNDED,
        style="bold yellow"
    )
    console.print(name_panel)
    new_model_name = Prompt.ask("[bold yellow]New Model Name[/bold yellow]")
    console.print("[bold magenta]Detecting available hardware...[/bold magenta]")
    console.print(render_hardware_panel())
    cont = Prompt.ask("[bold yellow]Continue with these settings?[/bold yellow]", choices=["y", "n"], default="y")
    if cont.lower() != "y":
        console.print("[bold red]Aborted by user.[/bold red]")
        return
    spinner_loading("Loading model and preparing for fine-tuning...", duration=2)
    run_trainer(selected_model, dataset, new_model_name)

def train_flow_remote():
    """
    Stub for remote API training (future).
    """
    console.print(Panel("[bold cyan]Remote API training is not yet implemented. Coming soon![/bold cyan]", style="bold yellow"))

def train_flow_ssh():
    """
    Stub for SSH remote training (future).
    """
    console.print(Panel("[bold cyan]SSH remote training is not yet implemented. Coming soon![/bold cyan]", style="bold yellow"))

def show_help():
    """
    Show help and usage information.
    """
    help_text = """
[bold magenta]Nexa Auto CLI Help[/bold magenta]

[bold yellow]Usage:[/bold yellow]
  python cli.py

[bold yellow]Commands:[/bold yellow]
  1  Train a model locally
  2  Train on remote API (Kaggle, Lambda, etc.) [future]
  3  Train via SSH on remote server [future]
  4  Show hardware info
  5  Clear session token
  6  Help
  7  Exit

[bold yellow]Flags:[/bold yellow]
  --help      Show this help message
  --version   Show version

[bold yellow]Docs:[/bold yellow]
  https://github.com/your-org/nexa-auto
"""
    console.print(Panel(help_text, style="bold yellow"))

def scaffold_optimizations():
    """
    Scaffold for future low-level optimizations (tokenizer, dataloader, etc.).
    """
    # Example: Use fast_tokenizer or fast_loader if available
    # tokenizer = fast_tokenizer(...) if use_fast else AutoTokenizer.from_pretrained(...)
    # dataloader = fast_loader(...) if use_fast else DataLoader(...)
    pass

def launch_cli():
    """
    Main CLI loop: handles session, menu, and command dispatch.
    """
    start_session_server()
    print_hf_logo()
    hf_token = get_token()
    while not hf_token:
        console.print("[bold red]No Hugging Face token found for this session.[/bold red]")
        hf_key = Prompt.ask("[bold yellow]Enter your Hugging Face access token[/bold yellow]", password=True)
        if hf_key.strip():
            store_token(hf_key)
            console.print("[bold yellow]Token stored for this session only.[/bold yellow]\n")
            hf_token = get_token()
        else:
            console.print("[bold red]Token cannot be empty. Please try again.[/bold red]")
    while True:
        console.clear()
        print_hf_logo()
        console.print(render_main_menu())
        cmd = Prompt.ask("[bold yellow]Enter command number[/bold yellow]", choices=[str(i) for i in range(1, 8)])
        if cmd == "1":
            train_flow_local()
            Prompt.ask("\n[bold yellow]Press Enter to return to main menu[/bold yellow]")
        elif cmd == "2":
            train_flow_remote()
            Prompt.ask("\n[bold yellow]Press Enter to return to main menu[/bold yellow]")
        elif cmd == "3":
            train_flow_ssh()
            Prompt.ask("\n[bold yellow]Press Enter to return to main menu[/bold yellow]")
        elif cmd == "4":
            console.print(render_hardware_panel())
            Prompt.ask("\n[bold yellow]Press Enter to return to main menu[/bold yellow]")
        elif cmd == "5":
            clear_token()
            console.print("[bold yellow]Session token cleared.[/bold yellow]")
            Prompt.ask("\n[bold yellow]Press Enter to return to main menu[/bold yellow]")
        elif cmd == "6":
            show_help()
            Prompt.ask("\n[bold yellow]Press Enter to return to main menu[/bold yellow]")
        elif cmd == "7":
            clear_token()
            console.print("[bold red]Goodbye![/bold red]")
            break

if __name__ == "__main__":
    if "--help" in sys.argv:
        show_help()
        sys.exit(0)
    if "--version" in sys.argv:
        console.print(f"[bold magenta]Nexa Auto CLI version {VERSION}[/bold magenta]")
        sys.exit(0)
    launch_cli()

# Note:
# For a true external GUI or TUI, consider using:
# - [Textual](https://textual.textualize.io/) for a modern TUI (Python)
# - PyQt/Tkinter for desktop GUI (Python)
# - Streamlit/Gradio for web UI (Python)
# - Rust/Go/C++/Node.js for native CLI/GUI if you want a different toolchain
# This script is a visually rich CLI using Python and Rich, and does not open a new terminal window.
