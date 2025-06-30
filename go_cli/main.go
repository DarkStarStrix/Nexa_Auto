package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
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
	clearLogs
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
    local           bool
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

// Add this helper to return to main menu after a delay
func returnToMainMenuAfterDelay() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		return menuLoadedMsg{}
	}
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
			if m.menuIdx == 0 { // TUI Mode
				m.state = mainMenu
				m.menuIdx = 0
				return m, nil
			} else { // CLI Mode
				appendLogFile("Launching CLI mode (cli.py)")
				go func() {
					cmd := exec.Command("python", "cli.py")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					cmd.Stdin = os.Stdin
					_ = cmd.Run()
					os.Exit(0)
				}()
				return m, tea.Quit
			}
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
				// Transition to fineTune state and immediately trigger backend ping
				m.state = fineTune
				m.loading = true
				m.backendStatus = "pending" // Mark ping as in progress
				m.appendLog("Started fine-tune session, pinging backend...")
				return m, checkBackendHealthCmd()
			case 1:
				m.state = logs
				m.showLog = true
				m.logs = loadLogs()
				return m, nil
			case 2:
				m.state = help
				return m, nil
			case 3:
				m.state = tokenMenu
				m.tokenStatus = ""
				return m, nil
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
		// Allow ESC to exit anytime
		if msg.String() == "esc" || msg.String() == "q" {
			m.state = mainMenu
			m.backendStatus = ""
			m.tokenStatus = ""
			m.tokenInput = ""
			m.appendLog("Exited fine-tune session")
			return m, nil
		}
		// If no ping has been initiated, immediately ping the backend.
		if m.backendStatus == "" {
			m.appendLog("Pinging backend...")
			m.backendStatus = "pending"
			return m, checkBackendHealthCmd()
		}
		// While backend ping is pending, keep ticking.
		if m.backendStatus == "pending" {
			return m, tickLoading()
		}
		// Once ping returns (backendStatus is not "pending") and no token prompt yet, set tokenStatus.
		if m.tokenStatus == "" {
			if strings.Contains(m.backendStatus, `"status":"ok"`) {
				m.tokenStatus = "Enter your Hugging Face token:"
			} else {
				m.tokenStatus = "Backend unavailable. Press ESC to return."
			}
			return m, nil
		}
		// Handle token input if healthy.
		if m.tokenStatus == "Enter your Hugging Face token:" {
			if msg.Type == tea.KeyRunes {
				m.tokenInput += msg.String()
				return m, nil
			}
			if msg.Type == tea.KeyEnter && m.tokenInput != "" {
				return m, setToken(m.tokenInput)
			}
		}
		// When token is set successfully, switch to model selection.
		if m.tokenStatus == "Token set successfully" {
			m.state = modelSelect
			m.menuIdx = 0
			m.backendStatus = ""
			m.tokenStatus = ""
			m.tokenInput = ""
			m.appendLog("Token set, proceeding to model selection")
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
			return m, sendTrainRequest(m)
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
		if msg.String() == "c" {
			m.state = clearLogs
			return m, clearLogFile()
		}
	case clearLogs:
		m.state = logs
		m.logs = loadLogs()
		return m, nil
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

type TrainRequest struct {
	Model   string `json:"model"`
	Dataset string `json:"dataset"`
	Output  string `json:"output"`
	Local   bool   `json:"local"`
}

type TrainResponse struct {
	JobID string `json:"job_id"`
}

func sendTrainRequest(m model) tea.Cmd {
	return func() tea.Msg {
		trainRequest := TrainRequest{
			Model:   modelOptions[m.selectedModel],
			Dataset: datasetOptions[m.selectedDataset],
			Output:  m.outputName,
			Local:   m.local,
		}
		jsonData, err := json.Marshal(trainRequest)
		if err != nil {
			return backendHealthMsg(fmt.Sprintf("Error marshaling JSON: %v", err))
		}

		resp, err := http.Post("http://localhost:8770/train", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return backendHealthMsg(fmt.Sprintf("Error sending request: %v", err))
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return backendHealthMsg(fmt.Sprintf("Error reading response: %v", err))
		}

		var trainResponse TrainResponse
		err = json.Unmarshal(body, &trainResponse)
		if err != nil {
			return backendHealthMsg(fmt.Sprintf("Error unmarshaling response: %v", err))
		}

		return backendHealthMsg(fmt.Sprintf("Training job started with job ID: %s", trainResponse.JobID))
	}
}

// --- Backend Health Check ---
// Simplify to use a single known endpoint for quick ping.
func checkBackendHealth() tea.Msg {
	const endpoint = "http://localhost:8770/health"
	resp, err := http.Get(endpoint)
	if err != nil {
		msg := fmt.Sprintf("Backend not available: %v", err)
		appendLogFile("Backend health checked: " + msg)
		return backendHealthMsg(msg)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("Unexpected status code %d from %s", resp.StatusCode, endpoint)
		appendLogFile("Backend health checked: " + msg)
		return backendHealthMsg(msg)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg := fmt.Sprintf("Error reading response: %v", err)
		appendLogFile("Backend health checked: " + msg)
		return backendHealthMsg(msg)
	}
	// This is what we expect:
	// Backend health checked: {"status":"ok","components":{"session_server":"ok","trainer":"ok"},"timestamp":"..."} (endpoint: http://localhost:8770/health)
	logMsg := fmt.Sprintf("Backend health checked: %s (endpoint: %s)", string(body), endpoint)
	appendLogFile(logMsg)
	return backendHealthMsg(string(body))
}

func checkBackendHealthCmd() tea.Cmd {
	return func() tea.Msg {
		return checkBackendHealth()
	}
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
	switch m.state {
	case modeSelect:
		out := m.splash() + "\n\n"
		out += headerStyle.Render("Choose Interface Mode") + "\n\n"
		for i, item := range modeOptions {
			indicator := ""
			if i == m.menuIdx && m.loadingMenu {
				indicator = fmt.Sprintf(" %s", loadingSlash[m.loadingFrame])
			}
			if i == m.menuIdx {
				out += selectedStyle.Render("> " + item + indicator) + "\n"
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
		out += "\n[q] Quit"
		return out
	case tokenMenu:
		return boxStyle.Render("[Token Management] (ESC/q to return)\n" +
			"1. Get Token\n2. Set Token\n3. Clear Token\n\n" + m.tokenStatus + m.tokenInput)
	case fineTune:
		var display string
		// While ping in progress, show loading message.
		if m.backendStatus == "pending" {
			display = "Checking backend..."
		} else if m.tokenStatus == "" {
			display = "Backend check complete."
		} else {
			if m.tokenStatus == "Enter your Hugging Face token:" {
				display = m.tokenStatus + " " + m.tokenInput
			} else {
				display = m.tokenStatus
			}
		}
		return boxStyle.Render("[Fine-tune] (ESC/q to return)\n\n" + display)
	case logs:
		logsContent := strings.Join(m.logs, "\n")
		return boxStyle.Render("[Logs] (ESC/q to return, c to clear)\n\n" +
			logsContent)
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
	case clearLogs:
		return boxStyle.Render("[Logs Cleared]")
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

// --- Clear Logs ---
func clearLogFile() tea.Cmd {
	return func() tea.Msg {
		err := os.Truncate("Tune.log", 0)
		if err != nil {
			return backendHealthMsg("Failed to clear log file: " + err.Error())
		}
		return backendHealthMsg("Log file cleared")
	}
}

// --- Main ---
func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
