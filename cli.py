"""
Nexa Auto CLI

A simple CLI main menu for fine-tuning orchestration, logs, help, and token management.
"""

import os
import requests

from rich.console import Console
from rich.panel import Panel
from rich.prompt import Prompt

console = Console()
VERSION = "0.3.0"

def print_nexa_splash():
    nexa_lines = [
        "███╗   ██╗███████╗██╗  ██╗ █████╗      █████╗ ██╗   ██╗████████╗ ██████╗ ",
        "████╗  ██║██╔════╝╚██╗██╔╝██╔══██╗    ██╔══██╗██║   ██║╚══██╔══╝██╔═══██╗",
        "██╔██╗ ██║█████╗   ╚███╔╝ ███████║    ███████║██║   ██║   ██║   ██║   ██║",
        "██║╚██╗██║██╔══╝   ██╔██╗ ██╔══██║    ██╔══██║██║   ██║   ██║   ██║   ██║",
        "██║ ╚████║███████╗██╔╝ ██╗██║  ██║    ██║  ██║╚██████╔╝   ██║   ╚██████╔╝",
        "╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝    ╚═╝  ╚═╝ ╚═════╝    ╚═╝    ╚═════╝ ",
    ]
    from rich.text import Text
    from rich.style import Style
    green = "#39FF14"
    styled_lines = [Text(line, style=Style(color=green, bold=True)) for line in nexa_lines]
    splash_text = "\n".join(str(line) for line in styled_lines)
    panel = Panel(
        splash_text,
        border_style=green,
        style="on #101820",
        padding=(1, 4),
        expand=False,
    )
    console.print(panel)

def show_main_menu():
    console.print("\n[bold magenta]Main Menu[/bold magenta]\n")
    console.print("[bold yellow]1.[/bold yellow] Fine-tune Model")
    console.print("[bold yellow]2.[/bold yellow] View Logs")
    console.print("[bold yellow]3.[/bold yellow] Help")
    console.print("[bold yellow]4.[/bold yellow] Token Management")
    console.print("[bold yellow]5.[/bold yellow] Exit\n")

def fine_tune_menu():
    hf_token = Prompt.ask("[bold yellow]Enter your Hugging Face token[/bold yellow]", password=True)
    os.environ["HF_TOKEN"] = hf_token
    console.print(f"[bold green]Token received. Starting fine-tune flow...[/bold green]")

    model = Prompt.ask("[bold yellow]Enter the model name[/bold yellow]")
    dataset = Prompt.ask("[bold yellow]Enter the dataset name[/bold yellow]")
    output = Prompt.ask("[bold yellow]Enter the output model name[/bold yellow]")
    local_str = Prompt.ask("[bold yellow]Is this a local job? (y/n)[/bold yellow]", choices=["y", "n"])
    local = local_str.lower() == "y"

    data = {
        "model": model,
        "dataset": dataset,
        "output": output,
        "local": local
    }

    try:
        response = requests.post("http://localhost:8770/train", json=data)
        response.raise_for_status()  # Raise HTTPError for bad responses (4xx or 5xx)
        job_id = response.json().get("job_id")
        console.print(f"[bold green]Training job started with job ID: {job_id}[/bold green]")
    except requests.exceptions.RequestException as e:
        console.print(f"[bold red]Error starting training job: {e}[/bold red]")

    Prompt.ask("\n[bold yellow]Press Enter to return to main menu[/bold yellow]")

def view_logs_menu():
    log_path = os.path.join(os.path.dirname(__file__), "go_cli", "Tune.log")
    if not os.path.exists(log_path):
        console.print("[bold red]No logs found.[/bold red]")
    else:
        with open(log_path, "r") as f:
            logs = f.read()
        console.print(Panel(logs if logs else "[No logs]", title="Tune.log", style="bold yellow"))
    Prompt.ask("\n[bold yellow]Press Enter to return to main menu[/bold yellow]")

def help_menu():
    help_text = """
[bold magenta]Nexa Auto CLI Help[/bold magenta]

[bold yellow]Menu Options:[/bold yellow]
  1  Fine-tune Model
  2  View Logs
  3  Help
  4  Token Management
  5  Exit

[bold yellow]Note:[/bold yellow]
  Hugging Face tokens must have read and write permissions to your account for fine-tuning and pushing models.

[bold yellow]Docs:[/bold yellow]
  https://github.com/your-org/nexa-auto
"""
    console.print(Panel(help_text, style="bold yellow"))
    Prompt.ask("\n[bold yellow]Press Enter to return to main menu[/bold yellow]")

def token_management_menu():
    while True:
        console.print("\n[bold magenta]Token Management[/bold magenta]")
        console.print("[bold yellow]1.[/bold yellow] Get Token")
        console.print("[bold yellow]2.[/bold yellow] Set Token")
        console.print("[bold yellow]3.[/bold yellow] Clear Token")
        console.print("[bold yellow]4.[/bold yellow] Back\n")
        choice = Prompt.ask("[bold yellow]Select an option[/bold yellow]", choices=["1", "2", "3", "4"])
        if choice == "1":
            token = os.environ.get("HF_TOKEN")
            if token:
                console.print(f"[bold green]Token: {token[:4]}...{token[-4:]}[/bold green]")
            else:
                console.print("[bold red]No token found.[/bold red]")
        elif choice == "2":
            token = Prompt.ask("[bold yellow]Enter new token[/bold yellow]", password=True)
            os.environ["HF_TOKEN"] = token
            console.print("[bold green]Token set for this session.[/bold green]")
        elif choice == "3":
            os.environ.pop("HF_TOKEN", None)
            console.print("[bold green]Token cleared.[/bold green]")
        elif choice == "4":
            break
        Prompt.ask("\n[bold yellow]Press Enter to return to token menu[/bold yellow]")

def main():
    while True:
        console.clear()
        print_nexa_splash()
        show_main_menu()
        choice = Prompt.ask("[bold yellow]Select an option[/bold yellow]", choices=["1", "2", "3", "4", "5"])
        if choice == "1":
            fine_tune_menu()
        elif choice == "2":
            view_logs_menu()
        elif choice == "3":
            help_menu()
        elif choice == "4":
            token_management_menu()
        elif choice == "5":
            console.print("[bold red]Goodbye![/bold red]")
            break

if __name__ == "__main__":
    main()
