package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- Splash Art ---
func CreateNexaSplash() string {
	nexaLines := []string{
		"███╗   ██╗███████╗██╗  ██╗ █████╗      █████╗ ██╗   ██╗████████╗ ██████╗ ",
		"████╗  ██║██╔════╝╚██╗██╔╝██╔══██╗    ██╔══██╗██║   ██║╚══██╔══╝██╔═══██╗",
		"██╔██╗ ██║█████╗   ╚███╔╝ ███████║    ███████║██║   ██║   ██║   ██║   ██║",
		"██║╚██╗██║██╔══╝   ██╔██╗ ██╔══██║    ██╔══██║██║   ██║   ██║   ██║   ██║",
		"██║ ╚████║███████╗██╔╝ ██╗██║  ██║    ██║  ██║╚██████╔╝   ██║   ╚██████╔╝",
		"╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝    ╚═╝  ╚═╝ ╚═════╝    ╚═╝    ╚═════╝ ",
	}
	green := lipgloss.Color("#39FF14")
	box := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(green).
		Background(lipgloss.Color("#101820")).
		Padding(1, 4).
		Align(lipgloss.Center)
	style := lipgloss.NewStyle().Foreground(green).Bold(true)
	var styledLines []string
	for _, line := range nexaLines {
		styledLines = append(styledLines, style.Render(line))
	}
	return box.Render(strings.Join(styledLines, "\n"))
}

// --- Constants and Styles ---
const (
	mainMenu state = iota
	fineTune
	logs
	help
	tokenMenu
	modelSelect
	datasetSelect
	outputName
	confirmRun
	modeSelect
)

var (
	boxStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(1, 2).
		BorderForeground(lipgloss.Color("#6C63FF")).
		Background(lipgloss.Color("#16161a"))
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#6C63FF")).
		Background(lipgloss.Color("#232946")).
		Padding(0, 1)
	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD21F")).
		Background(lipgloss.Color("#232946")).
		Bold(true).
		Padding(0, 1)
	loadingSlash = []string{"|", "/", "-", "\\"}
	modelOptions   = []string{"mistral-7b", "llama-2-7b", "custom..."}
	datasetOptions = []string{"local.jsonl", "hf-dataset", "custom..."}
	modeOptions    = []string{"TUI Mode (modern)", "Classic CLI Mode"}
)

// --- Types ---
type state int

type model struct {
	state           state
	menuIdx         int
	loading         bool
	loadingFrame    int
	output          string
	backendStatus   string
	tokenStatus     string
	tokenInput      string
	selectedModel   int
	selectedDataset int
	outputName      string
	confirmMsg      string
	mode            int // 0 = TUI, 1 = CLI
	logs            []string
	showLog         bool
	loadingMenu     bool
	cliStyle        lipgloss.Style
}

// --- Model Initialization ---
func initialModel() model {
	return model{
		state:      modeSelect,
		logs:       loadLogs(),
		loadingMenu: false,
		loadingFrame: 0,
		cliStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4A90E2")).
			Padding(1, 2),
	}
}

// --- Bubbletea Init ---
func (m model) Init() tea.Cmd {
	if m.state == fineTune && m.loading {
		return tea.Batch(checkBackendHealth, tickLoading())
	}
	return nil
}

// --- Main Update ---
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.mode == 1 {
			return m.updateCLI(msg)
		}
		return m.updateTUI(msg)
	case backendHealthMsg:
		m.loading = false
		m.backendStatus = string(msg)
		m.appendLog("Backend health checked: " + m.backendStatus)
	case tokenStatusMsg:
		m.tokenStatus = string(msg)
		m.tokenInput = ""
		m.appendLog("Token status: " + m.tokenStatus)
	case tickMsg:
		if m.loading {
			m.loadingFrame = (m.loadingFrame + 1) % len(loadingSlash)
			return m, tickLoading()
		}
	case menuLoadedMsg:
		m.loadingMenu = false
	case tickMenuMsg:
		if m.loadingMenu {
			m.loadingFrame = (m.loadingFrame + 1) % len(loadingSlash)
			return m, tickMenuLoading()
		}
	}
	return m, nil
}

// --- TUI Update ---
func (m model) updateTUI(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case modeSelect:
		switch msg.String() {
		case "j", "down":
			m.menuIdx = (m.menuIdx + 1) % len(modeOptions)
		case "k", "up":
			m.menuIdx = (m.menuIdx + len(modeOptions) - 1) % len(modeOptions)
		case "enter":
			m.mode = m.menuIdx
			m.state = mainMenu
			m.menuIdx = 0
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case mainMenu:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			m.menuIdx = (m.menuIdx + 1) % 4
		case "k", "up":
			m.menuIdx = (m.menuIdx + 3) % 4
		case "enter":
			m.loadingMenu = true
			m.loadingFrame = 0
			switch m.menuIdx {
			case 0:
				m.state = fineTune
				m.loading = true
				m.appendLog("Started fine-tune session")
				return m, tea.Batch(checkBackendHealth, tickLoading(), tickMenuLoading())
			case 1:
				m.state = logs
				m.showLog = true
				m.logs = loadLogs()
				return m, tickMenuLoading()
			case 2:
				m.state = help
				return m, tickMenuLoading()
			case 3:
				m.state = tokenMenu
				m.tokenStatus = ""
				return m, tickMenuLoading()
			}
		case "esc":
			m.state = modeSelect
		}
	case tokenMenu:
		switch msg.String() {
		case "esc", "q":
			m.state = mainMenu
			m.tokenInput = ""
			return m, nil
		case "1":
			return m, getToken
		case "2":
			m.tokenInput = ""
			m.tokenStatus = "Enter token: "
			return m, nil
		case "3":
			return m, clearToken
		}
		if m.tokenStatus == "Enter token: " && msg.Type == tea.KeyRunes {
			m.tokenInput += msg.String()
			return m, nil
		}
		if m.tokenStatus == "Enter token: " && msg.Type == tea.KeyEnter {
			return m, setToken(m.tokenInput)
		}
	case fineTune:
		if msg.String() == "esc" || msg.String() == "q" {
			m.state = mainMenu
			m.backendStatus = ""
			m.loading = false
			m.appendLog("Exited fine-tune session")
			return m, nil
		}
		if !m.loading && m.backendStatus != "" && !strings.Contains(m.backendStatus, "Token:") {
			m.state = tokenMenu
			m.tokenStatus = "No HF token found. Please set your Hugging Face token.\n"
			m.appendLog("Prompted for Hugging Face token")
			return m, nil
		}
		if !m.loading && m.backendStatus != "" && strings.Contains(m.backendStatus, "Token:") {
			m.state = modelSelect
			m.menuIdx = 0
			m.appendLog("Backend ready, proceeding to model selection")
			return m, nil
		}
	case modelSelect:
		switch msg.String() {
		case "esc":
			m.state = mainMenu
		case "j", "down":
			m.menuIdx = (m.menuIdx + 1) % len(modelOptions)
		case "k", "up":
			m.menuIdx = (m.menuIdx + len(modelOptions) - 1) % len(modelOptions)
		case "enter":
			m.selectedModel = m.menuIdx
			m.state = datasetSelect
			m.menuIdx = 0
			m.appendLog(fmt.Sprintf("Selected model: %s", modelOptions[m.selectedModel]))
		}
	case datasetSelect:
		switch msg.String() {
		case "esc":
			m.state = mainMenu
		case "j", "down":
			m.menuIdx = (m.menuIdx + 1) % len(datasetOptions)
		case "k", "up":
			m.menuIdx = (m.menuIdx + len(datasetOptions) - 1) % len(datasetOptions)
		case "enter":
			m.selectedDataset = m.menuIdx
			m.state = outputName
			m.outputName = ""
			m.appendLog(fmt.Sprintf("Selected dataset: %s", datasetOptions[m.selectedDataset]))
		}
	case outputName:
		switch msg.Type {
		case tea.KeyRunes:
			m.outputName += msg.String()
		case tea.KeyBackspace:
			if len(m.outputName) > 0 {
				m.outputName = m.outputName[:len(m.outputName)-1]
			}
		case tea.KeyEnter:
			if m.outputName != "" {
				m.state = confirmRun
				m.confirmMsg = ""
				m.appendLog(fmt.Sprintf("Set output name: %s", m.outputName))
			}
		case tea.KeyEsc:
			m.state = mainMenu
		}
	case confirmRun:
		switch msg.String() {
		case "y":
			m.confirmMsg = "[TODO] Launching fine-tune job..."
			m.appendLog("Confirmed fine-tune run")
		case "n", "esc":
			m.state = mainMenu
		}
	case help:
		if msg.String() == "esc" || msg.String() == "q" {
			m.state = mainMenu
		}
	case logs:
		if msg.String() == "esc" || msg.String() == "q" {
			m.state = mainMenu
			m.showLog = false
		}
	}
	return m, nil
}

// --- CLI Update ---
func (m model) updateCLI(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case mainMenu:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			m.menuIdx = (m.menuIdx + 1) % 4
		case "k", "up":
			m.menuIdx = (m.menuIdx + 3) % 4
		case "enter":
			m.loadingMenu = true
			m.loadingFrame = 0
			switch m.menuIdx {
			case 0:
				m.state = fineTune
				m.loading = true
				m.appendLog("Started fine-tune session (CLI)")
				return m, tea.Batch(checkBackendHealth, tickLoading(), tickMenuLoading())
			case 1:
				m.state = logs
				m.showLog = true
				m.logs = loadLogs()
				return m, tickMenuLoading()
			case 2:
				m.state = help
				return m, tickMenuLoading()
			case 3:
				m.state = tokenMenu
				m.tokenStatus = ""
				return m, tickMenuLoading()
			}
		}
	case fineTune:
		if msg.String() == "esc" || msg.String() == "q" {
			m.state = mainMenu
			m.loading = false
			return m, nil
		}
	case logs:
		if msg.String() == "esc" || msg.String() == "q" {
			m.state = mainMenu
			m.showLog = false
			return m, nil
		}
	case help:
		if msg.String() == "esc" || msg.String() == "q" {
			m.state = mainMenu
			return m, nil
		}
	case tokenMenu:
		switch msg.String() {
		case "esc", "q":
			m.state = mainMenu
			m.tokenInput = ""
			return m, nil
		case "1":
			return m, getToken
		case "2":
			m.tokenInput = ""
			m.tokenStatus = "Enter token: "
			return m, nil
		case "3":
			return m, clearToken
		}
		if m.tokenStatus == "Enter token: " {
			if msg.Type == tea.KeyRunes {
				m.tokenInput += msg.String()
			} else if msg.Type == tea.KeyEnter {
				return m, setToken(m.tokenInput)
			}
		}
	}
	return m, nil
}

// --- Message Types ---
type backendHealthMsg string
type tokenStatusMsg string
type tickMsg struct{}
type menuLoadedMsg struct{}
type tickMenuMsg struct{}

// --- Backend Health Check ---
func checkBackendHealth() tea.Msg {
	resp, err := http.Get("http://localhost:8000/health")
	if err != nil {
		return backendHealthMsg("Backend not available: " + err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return backendHealthMsg("Error reading response: " + err.Error())
	}
	return backendHealthMsg(string(body))
}

// --- Token Management ---
func getToken() tea.Msg {
	token := os.Getenv("HF_TOKEN")
	if token == "" {
		return tokenStatusMsg("No token found")
	}
	return tokenStatusMsg("Token: " + token[:4] + "..." + token[len(token)-4:])
}

func setToken(token string) tea.Cmd {
	return func() tea.Msg {
		err := os.Setenv("HF_TOKEN", token)
		if err != nil {
			return tokenStatusMsg("Failed to set token: " + err.Error())
		}
		return tokenStatusMsg("Token set successfully")
	}
}

func clearToken() tea.Msg {
	err := os.Unsetenv("HF_TOKEN")
	if err != nil {
		return tokenStatusMsg("Failed to clear token: " + err.Error())
	}
	return tokenStatusMsg("Token cleared")
}

// --- Loading Spinner ---
func tickLoading() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(120 * time.Millisecond)
		return tickMsg{}
	}
}

func tickMenuLoading() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(120 * time.Millisecond)
		return tickMenuMsg{}
	}
}

// --- Views ---
func (m model) splash() string {
	return CreateNexaSplash()
}

func (m model) View() string {
	if m.mode == 1 {
		return m.cliView()
	}
	return m.tuiView()
}

func (m model) tuiView() string {
	loadingIndicator := ""
	if m.loadingMenu {
		loadingIndicator = fmt.Sprintf(" %s Loading...", loadingSlash[m.loadingFrame])
	}
	switch m.state {
	case modeSelect:
		out := m.splash() + "\n\n"
		out += headerStyle.Render("Choose Interface Mode") + "\n\n"
		for i, item := range modeOptions {
			if i == m.menuIdx {
				out += selectedStyle.Render("> " + item) + "\n"
			} else {
				out += "  " + item + "\n"
			}
		}
		out += "\n[q] Quit"
		return out
	case mainMenu:
		menu := []string{"Fine-tune Model", "View Logs", "Help", "Token Management"}
		out := CreateNexaSplash() + "\n\n"
		out += headerStyle.Render("Main Menu") + "\n\n"
		for i, item := range menu {
			if i == m.menuIdx {
				out += selectedStyle.Render("> " + item) + "\n"
			} else {
				out += "  " + item + "\n"
			}
		}
		out += "\n[q] Quit " + loadingIndicator
		return out
	case tokenMenu:
		return boxStyle.Render("[Token Management] (ESC/q to return)\n" +
			"1. Get Token\n2. Set Token\n3. Clear Token\n\n" + m.tokenStatus + m.tokenInput)
	case fineTune:
		status := ""
		if m.loading {
			status = fmt.Sprintf("Checking backend... %s", loadingSlash[m.loadingFrame])
		} else if m.backendStatus != "" {
			status = m.backendStatus
		}
		return boxStyle.Render("[Fine-tune] (ESC/q to return)\n\n" + status + "\n\n" +
			"If prompted, enter your Hugging Face token.\n" +
			"TODO: Implement model/dataset/config selection and REST call to backend.")
	case logs:
		return boxStyle.Render("[Logs] (ESC/q to return)\n\n" +
			"Log file: Tune.log\n\n" +
			strings.Join(m.logs, "\n") +
			"\n\n[Open Tune.log in your file explorer to view or share logs.]")
	case help:
		return boxStyle.Render("[Help] (ESC/q to return)\n\n" +
			"Commands:\n" +
			"  ↑/↓ or j/k   - Navigate menu\n" +
			"  Enter        - Select\n" +
			"  q or ESC     - Back/Exit\n" +
			"  1/2/3        - Token actions in Token Management\n\n" +
			"Workflow:\n" +
			"  1. Fine-tune: Checks backend, prompts for HF token if needed, then launches session\n" +
			"  2. Token Management: Set, get, or clear your Hugging Face token\n" +
			"  3. Logs: View backend logs (future)\n" +
			"  4. Help: Show this help screen\n")
	case modelSelect:
		out := headerStyle.Render("Select Model") + "\n\n"
		for i, opt := range modelOptions {
			if i == m.menuIdx {
				out += selectedStyle.Render("> " + opt) + "\n"
			} else {
				out += "  " + opt + "\n"
			}
		}
		out += "\n[ESC to cancel]"
		return boxStyle.Render(out)
	case datasetSelect:
		out := headerStyle.Render("Select Dataset") + "\n\n"
		for i, opt := range datasetOptions {
			if i == m.menuIdx {
				out += selectedStyle.Render("> " + opt) + "\n"
			} else {
				out += "  " + opt + "\n"
			}
		}
		out += "\n[ESC to cancel]"
		return boxStyle.Render(out)
	case outputName:
		return boxStyle.Render(headerStyle.Render("Output Name") + "\n\n" + m.outputName + "_adapter" + "\n\n[Type name, Enter to confirm, ESC to cancel]")
	case confirmRun:
		return boxStyle.Render(headerStyle.Render("Confirm Fine-tune") +
			fmt.Sprintf("\n\nModel: %s\nDataset: %s\nOutput: %s_adapter\n\nProceed? (y/n)\n%s",
				modelOptions[m.selectedModel], datasetOptions[m.selectedDataset], m.outputName, m.confirmMsg))
	}
	return ""
}

func (m model) cliView() string {
	loadingIndicator := ""
	if m.loadingMenu {
		loadingIndicator = fmt.Sprintf("%s ", loadingSlash[m.loadingFrame])
	}
	switch m.state {
	case mainMenu:
		menu := []string{"Fine-tune Model", "View Logs", "Help", "Token Management"}
		out := "\n" + headerStyle.Render("NEXA CLI") + "\n\n"
		for i, item := range menu {
			if i == m.menuIdx {
				out += loadingIndicator + selectedStyle.Render(item) + "\n"
			} else {
				out += "  " + item + "\n"
			}
		}
		return m.cliStyle.Render(out + "\n[q] Quit")
	case fineTune:
		status := ""
		if m.loading {
			status = fmt.Sprintf("%sChecking backend...", loadingSlash[m.loadingFrame])
		} else if m.backendStatus != "" {
			status = m.backendStatus
		}
		return m.cliStyle.Render(fmt.Sprintf("[Fine-tune]\n\n%s\n", status))
	case logs:
		return m.cliStyle.Render("[Logs]\n\n" + strings.Join(m.logs, "\n"))
	}
	return ""
}

// --- Logging ---
func (m *model) appendLog(entry string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s", timestamp, entry)
	m.logs = append(m.logs, logEntry)
	appendLogFile(logEntry)
}

func appendLogFile(entry string) {
	f, err := os.OpenFile("Tune.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(entry + "\n")
}

func loadLogs() []string {
	data, err := ioutil.ReadFile("Tune.log")
	if err != nil {
		return []string{}
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// --- Main ---
func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}