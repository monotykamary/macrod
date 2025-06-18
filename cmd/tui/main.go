package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/monotykamary/macrod/internal/ipc"
	"github.com/monotykamary/macrod/pkg/models"
)

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
	
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Padding(0, 1)
	
	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(0, 1)
)

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Help     key.Binding
	Quit     key.Binding
	Toggle   key.Binding
	Delete   key.Binding
	New      key.Binding
	Record   key.Binding
	Play     key.Binding
	Edit     key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Toggle: key.NewBinding(
		key.WithKeys(" ", "enter"),
		key.WithHelp("space/enter", "toggle enable/disable"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d", "delete"),
		key.WithHelp("d/del", "delete macro"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new macro"),
	),
	Record: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "record macro"),
	),
	Play: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "play macro"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit macro"),
	),
}

type state int

const (
	stateTable state = iota
	stateRecording
	stateEditMacro
	stateConfirmDelete
)

type model struct {
	state         state
	table         table.Model
	macros        []models.Macro
	help          help.Model
	showHelp      bool
	width         int
	height        int
	daemonRunning bool
	recording     bool
	ipcClient     *ipc.Client
	err           error
	
	// Recording state
	recordedKeys  []string
	nameInput     textinput.Model
	descInput     textinput.Model
	hotkeyInput   textinput.Model
	activeInput   int
	isRecordingKeys bool  // true when recording keys, false when filling form
	
	// Delete confirmation
	deleteTarget  *models.Macro
	
	// Edit state
	editTarget    *models.Macro
}

func newModel() model {
	columns := []table.Column{
		{Title: "Status", Width: 8},
		{Title: "Name", Width: 20},
		{Title: "Description", Width: 30},
		{Title: "Hotkey", Width: 15},
		{Title: "Actions", Width: 10},
		{Title: "Created", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	// Initialize text inputs
	nameInput := textinput.New()
	nameInput.Placeholder = "Macro name"
	nameInput.Focus()
	nameInput.CharLimit = 30
	nameInput.Width = 30

	descInput := textinput.New()
	descInput.Placeholder = "Description"
	descInput.CharLimit = 50
	descInput.Width = 50

	hotkeyInput := textinput.New()
	hotkeyInput.Placeholder = "Hotkey (e.g., ctrl+1)"
	hotkeyInput.CharLimit = 20
	hotkeyInput.Width = 20

	m := model{
		state:       stateTable,
		table:       t,
		help:        help.New(),
		showHelp:    false,
		nameInput:   nameInput,
		descInput:   descInput,
		hotkeyInput: hotkeyInput,
	}

	// Try to connect to daemon
	client, err := ipc.NewClient()
	if err != nil {
		m.err = err
		m.macros = getMockMacros() // Use mock data if daemon is not running
	} else {
		m.ipcClient = client
		m.daemonRunning = true
		// Load macros from daemon
		if macros, err := client.ListMacros(); err == nil {
			m.macros = macros
		} else {
			m.macros = getMockMacros()
		}
	}

	m.updateTable()
	return m
}

func (m *model) updateTable() {
	rows := []table.Row{}
	for _, macro := range m.macros {
		status := "❌"
		if macro.Enabled {
			status = "✅"
		}
		rows = append(rows, table.Row{
			status,
			macro.Name,
			macro.Description,
			macro.Hotkey,
			fmt.Sprintf("%d", len(macro.Actions)),
			macro.CreatedAt.Format("2006-01-02 15:04"),
		})
	}
	m.table.SetRows(rows)
}

func getMockMacros() []models.Macro {
	return []models.Macro{
		{
			ID:          "1",
			Name:        "Combo 1",
			Description: "Basic attack combo",
			Hotkey:      "Ctrl+1",
			Actions: []models.KeyAction{
				{Key: "a", Delay: 100 * time.Millisecond},
				{Key: "b", Delay: 100 * time.Millisecond},
				{Key: "c", Delay: 100 * time.Millisecond},
			},
			Enabled:   true,
			CreatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          "2",
			Name:        "Special Move",
			Description: "Quarter circle forward + punch",
			Hotkey:      "Ctrl+2",
			Actions: []models.KeyAction{
				{Key: "down", Delay: 50 * time.Millisecond},
				{Key: "right", Delay: 50 * time.Millisecond},
				{Key: "x", Delay: 50 * time.Millisecond},
			},
			Enabled:   false,
			CreatedAt: time.Now().Add(-48 * time.Hour),
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetHeight(msg.Height - 10)
		m.help.Width = msg.Width

	case tea.KeyMsg:
		// Handle different states
		switch m.state {
		case stateRecording:
			return m.updateRecording(msg)
		case stateConfirmDelete:
			return m.updateConfirmDelete(msg)
		case stateEditMacro:
			return m.updateEditMacro(msg)
		}
		if m.showHelp {
			switch {
			case key.Matches(msg, keys.Help):
				m.showHelp = false
			case key.Matches(msg, keys.Quit):
				return m, tea.Quit
			}
			return m, nil
		}

		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Help):
			m.showHelp = true
		case key.Matches(msg, keys.Toggle):
			if len(m.macros) > 0 && m.table.Cursor() < len(m.macros) {
				macro := m.macros[m.table.Cursor()]
				if m.ipcClient != nil {
					if err := m.ipcClient.ToggleMacro(macro.ID); err != nil {
						m.err = err
					} else {
						// Reload macros from daemon
						if macros, err := m.ipcClient.ListMacros(); err == nil {
							m.macros = macros
							m.updateTable()
						}
					}
				} else {
					// Offline mode
					m.macros[m.table.Cursor()].Enabled = !m.macros[m.table.Cursor()].Enabled
					m.updateTable()
				}
			}
		case key.Matches(msg, keys.Delete):
			if len(m.macros) > 0 && m.table.Cursor() < len(m.macros) {
				// Show confirmation dialog
				m.deleteTarget = &m.macros[m.table.Cursor()]
				m.state = stateConfirmDelete
			}
		case key.Matches(msg, keys.New):
			// TODO: Open macro creation view
			return m, tea.Printf("Creating new macro...")
		case key.Matches(msg, keys.Record):
			if m.ipcClient != nil {
				// Start recording mode
				m.state = stateRecording
				m.recordedKeys = []string{}
				m.activeInput = 0
				m.nameInput.Reset()
				m.descInput.Reset()
				m.hotkeyInput.Reset()
				m.isRecordingKeys = true  // Start in key recording mode
				m.nameInput.Blur()  // Don't focus input yet
				if err := m.ipcClient.StartRecording(); err != nil {
					m.err = err
					m.state = stateTable
				} else {
					m.recording = true
				}
			} else {
				return m, tea.Printf("Daemon not running")
			}
		case key.Matches(msg, keys.Play):
			if len(m.macros) > 0 && m.table.Cursor() < len(m.macros) {
				macro := m.macros[m.table.Cursor()]
				if m.ipcClient != nil {
					go func() {
						if err := m.ipcClient.PlayMacro(macro.ID); err != nil {
							// Can't update error from goroutine, just log it
							fmt.Printf("Error playing macro: %v\n", err)
						}
					}()
					return m, tea.Printf("Playing macro: %s", macro.Name)
				} else {
					return m, tea.Printf("Daemon not running")
				}
			}
		case key.Matches(msg, keys.Edit):
			if len(m.macros) > 0 && m.table.Cursor() < len(m.macros) {
				// Start editing
				macro := m.macros[m.table.Cursor()]
				m.editTarget = &macro
				m.state = stateEditMacro
				m.activeInput = 0
				
				// Pre-fill the inputs with current values
				m.nameInput.SetValue(macro.Name)
				m.descInput.SetValue(macro.Description)
				m.hotkeyInput.SetValue(macro.Hotkey)
				m.nameInput.Focus()
				
				// Convert actions back to key strings for display
				m.recordedKeys = []string{}
				for _, action := range macro.Actions {
					m.recordedKeys = append(m.recordedKeys, action.Key)
				}
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) updateEditMacro(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg.String() {
	case "ctrl+c", "esc":
		// Cancel editing
		m.state = stateTable
		m.editTarget = nil
		return m, nil
		
	case "tab", "shift+tab":
		// Cycle through inputs
		if msg.String() == "tab" {
			m.activeInput = (m.activeInput + 1) % 3
		} else {
			m.activeInput = (m.activeInput + 2) % 3 // Go backwards
		}
		
		// Update focus
		m.nameInput.Blur()
		m.descInput.Blur()
		m.hotkeyInput.Blur()
		
		switch m.activeInput {
		case 0:
			m.nameInput.Focus()
		case 1:
			m.descInput.Focus()
		case 2:
			m.hotkeyInput.Focus()
		}
		
	case "enter":
		if m.activeInput == 2 && m.nameInput.Value() != "" && m.editTarget != nil {
			// Save the edited macro
			if m.ipcClient != nil {
				// Update the macro
				updatedMacro := *m.editTarget
				updatedMacro.Name = m.nameInput.Value()
				updatedMacro.Description = m.descInput.Value()
				updatedMacro.Hotkey = m.hotkeyInput.Value()
				updatedMacro.UpdatedAt = time.Now()
				
				if err := m.ipcClient.UpdateMacro(&updatedMacro); err != nil {
					m.err = err
				} else {
					// Reload macros
					if macros, err := m.ipcClient.ListMacros(); err == nil {
						m.macros = macros
						m.updateTable()
					}
					m.state = stateTable
					m.editTarget = nil
					return m, tea.Printf("Macro updated: %s", updatedMacro.Name)
				}
			}
		}
		
	default:
		// Update the active input
		switch m.activeInput {
		case 0:
			m.nameInput, cmd = m.nameInput.Update(msg)
		case 1:
			m.descInput, cmd = m.descInput.Update(msg)
		case 2:
			m.hotkeyInput, cmd = m.hotkeyInput.Update(msg)
		}
	}
	
	return m, cmd
}

func (m model) updateConfirmDelete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		// Perform deletion
		if m.deleteTarget != nil && m.ipcClient != nil {
			if err := m.ipcClient.DeleteMacro(m.deleteTarget.ID); err != nil {
				m.err = err
			} else {
				// Reload macros from daemon
				if macros, err := m.ipcClient.ListMacros(); err == nil {
					m.macros = macros
					m.updateTable()
				}
			}
		}
		m.state = stateTable
		m.deleteTarget = nil
		return m, nil
		
	case "n", "N", "esc":
		// Cancel deletion
		m.state = stateTable
		m.deleteTarget = nil
		return m, nil
	}
	
	return m, nil
}

func (m model) updateRecording(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	// If we're recording keys
	if m.isRecordingKeys {
		switch msg.String() {
		case "ctrl+c", "esc":
			// Stop recording and move to form
			m.isRecordingKeys = false
			m.nameInput.Focus()
			m.activeInput = 0
			
			// Pause the daemon recording
			if m.ipcClient != nil {
				m.ipcClient.PauseRecording()
			}
			
			// Get the recorded keys from daemon
			if m.ipcClient != nil {
				if keys, err := m.ipcClient.GetRecordingStatus(); err == nil {
					// Convert KeyActions to display strings
					m.recordedKeys = []string{}
					for _, action := range keys {
						m.recordedKeys = append(m.recordedKeys, action.Key)
					}
				}
			}
			return m, nil
			
		default:
			// While in recording mode, we just wait for the user to press Esc
			// The daemon is capturing all keys globally
			return m, nil
		}
	}
	
	// Otherwise we're in form mode
	switch msg.String() {
	case "ctrl+c", "esc":
		// Cancel everything
		if m.ipcClient != nil {
			m.ipcClient.StopRecording("", "", "")
		}
		m.recording = false
		m.state = stateTable
		m.isRecordingKeys = false
		return m, nil
		
	case "tab", "shift+tab":
		// Cycle through inputs
		if msg.String() == "tab" {
			m.activeInput = (m.activeInput + 1) % 3
		} else {
			m.activeInput = (m.activeInput + 2) % 3 // Go backwards
		}
		
		// Update focus
		m.nameInput.Blur()
		m.descInput.Blur()
		m.hotkeyInput.Blur()
		
		switch m.activeInput {
		case 0:
			m.nameInput.Focus()
		case 1:
			m.descInput.Focus()
		case 2:
			m.hotkeyInput.Focus()
		}
		
	case "enter":
		if m.activeInput == 2 && m.nameInput.Value() != "" {
			// Save the macro
			if m.ipcClient != nil {
				if macro, err := m.ipcClient.StopRecording(
					m.nameInput.Value(),
					m.descInput.Value(),
					m.hotkeyInput.Value(),
				); err != nil {
					m.err = err
				} else {
					m.recording = false
					m.state = stateTable
					// Reload macros
					if macros, err := m.ipcClient.ListMacros(); err == nil {
						m.macros = macros
						m.updateTable()
					}
					return m, tea.Printf("Macro recorded: %s", macro.Name)
				}
			}
		}
		
	default:
		
		// Update the active input
		switch m.activeInput {
		case 0:
			m.nameInput, cmd = m.nameInput.Update(msg)
		case 1:
			m.descInput, cmd = m.descInput.Update(msg)
		case 2:
			m.hotkeyInput, cmd = m.hotkeyInput.Update(msg)
		}
	}
	
	return m, cmd
}

func isSpecialKey(key string) bool {
	specialKeys := []string{
		"up", "down", "left", "right",
		"space", "enter", "backspace", "delete", "tab", "escape", "esc",
		"f1", "f2", "f3", "f4", "f5", "f6",
		"f7", "f8", "f9", "f10", "f11", "f12",
	}
	
	for _, sk := range specialKeys {
		if key == sk {
			return true
		}
	}
	return false
}

// formatKeyDisplay formats a key for display in the UI
func formatKeyDisplay(key string) string {
	// Handle shift indicator
	if strings.HasSuffix(key, " (shift)") {
		baseKey := strings.TrimSuffix(key, " (shift)")
		// For single letters, show as uppercase
		if len(baseKey) == 1 && baseKey >= "a" && baseKey <= "z" {
			return strings.ToUpper(baseKey)
		}
		// For other keys, show with shift indicator
		return formatKeyDisplay(baseKey) + "⇧"
	}
	
	switch key {
	case "space":
		return "␣"
	case "enter":
		return "⏎"
	case "tab":
		return "⇥"
	case "escape", "esc":
		return "⎋"
	case "backspace":
		return "⌫"
	case "delete":
		return "⌦"
	case "up":
		return "↑"
	case "down":
		return "↓"
	case "left":
		return "←"
	case "right":
		return "→"
	default:
		return key
	}
}

func (m model) View() string {
	if m.showHelp {
		return m.viewHelp()
	}
	
	switch m.state {
	case stateRecording:
		return m.viewRecording()
	case stateConfirmDelete:
		return m.viewConfirmDelete()
	case stateEditMacro:
		return m.viewEditMacro()
	}

	s := titleStyle.Render("🎮 Macro Daemon TUI") + "\n\n"
	
	// Status bar
	daemonStatus := "🔴 Daemon: Offline"
	if m.daemonRunning {
		daemonStatus = "🟢 Daemon: Running"
	}
	
	recordStatus := ""
	if m.recording {
		recordStatus = " | 🔴 Recording..."
	}
	
	s += statusStyle.Render(daemonStatus + recordStatus) + "\n\n"
	
	// Table
	s += baseStyle.Render(m.table.View()) + "\n\n"
	
	// Error display
	if m.err != nil {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
		errMsg := m.err.Error()
		
		// Special formatting for permission errors
		if strings.Contains(errMsg, "Accessibility permissions") {
			errStyle = errStyle.Bold(true)
			s += errStyle.Render("⚠️  "+errMsg) + "\n\n"
		} else {
			s += errStyle.Render("Error: "+errMsg) + "\n\n"
		}
	}
	
	// Help hint with key shortcuts
	helpText := "space: toggle • e: edit • p: play • r: record • d: delete • ?: help • q: quit"
	s += statusStyle.Render(helpText) + "\n"
	
	return s
}

func (m model) viewEditMacro() string {
	s := titleStyle.Render("✏️  Edit Macro") + "\n\n"
	
	// Show current keys (read-only for now)
	s += "Current keys: "
	if len(m.recordedKeys) == 0 {
		s += statusStyle.Render("(none)")
	} else {
		keyStr := ""
		for i, k := range m.recordedKeys {
			if i > 0 {
				keyStr += " · "
			}
			keyStr += formatKeyDisplay(k)
		}
		s += lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Render(keyStr)
	}
	s += "\n\n"
	
	// Show inputs
	s += "Edit Macro Details:\n\n"
	
	inputStyle := lipgloss.NewStyle().PaddingLeft(2)
	activeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
	
	// Name input
	label := "Name: "
	if m.activeInput == 0 {
		label = activeStyle.Render(label)
	}
	s += inputStyle.Render(label + m.nameInput.View()) + "\n\n"
	
	// Description input
	label = "Description: "
	if m.activeInput == 1 {
		label = activeStyle.Render(label)
	}
	s += inputStyle.Render(label + m.descInput.View()) + "\n\n"
	
	// Hotkey input
	label = "Hotkey: "
	if m.activeInput == 2 {
		label = activeStyle.Render(label)
	}
	s += inputStyle.Render(label + m.hotkeyInput.View()) + "\n\n"
	
	// Instructions
	s += statusStyle.Render("Tab/Shift+Tab: Navigate • Enter: Save (when on hotkey) • Esc: Cancel") + "\n"
	s += statusStyle.Render("Note: Key sequence editing not yet supported") + "\n"
	
	return s
}

func (m model) viewConfirmDelete() string {
	s := titleStyle.Render("⚠️  Confirm Deletion") + "\n\n"
	
	if m.deleteTarget != nil {
		confirmStyle := lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("196"))
		
		content := fmt.Sprintf("Are you sure you want to delete this macro?\n\n"+
			"Name: %s\n"+
			"Description: %s\n"+
			"Hotkey: %s\n\n"+
			"Press Y to confirm, N or Esc to cancel",
			m.deleteTarget.Name,
			m.deleteTarget.Description,
			m.deleteTarget.Hotkey)
		
		s += confirmStyle.Render(content)
	}
	
	return s
}

func (m model) viewRecording() string {
	var s string
	
	if m.isRecordingKeys {
		// Key recording mode
		s = titleStyle.Render("🔴 Recording Keys") + "\n\n"
		s += "Press keys to record them. Press Esc or Ctrl+C when done.\n\n"
		
		// Show recording status
		s += "Status: " + lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Render("Recording globally...") + "\n\n"
		s += "The daemon is capturing all your keystrokes.\n"
		s += "Type your key sequence naturally.\n\n"
		
		s += statusStyle.Render("Press Esc or Ctrl+C to finish recording and fill in macro details")
		return s
	}
	
	// Form mode
	s = titleStyle.Render("✏️  Macro Details") + "\n\n"
	
	// Show recorded keys
	s += "Recorded keys: "
	if len(m.recordedKeys) == 0 {
		s += statusStyle.Render("(none)")
	} else {
		keyStr := ""
		for i, k := range m.recordedKeys {
			if i > 0 {
				keyStr += " · "
			}
			keyStr += formatKeyDisplay(k)
		}
		s += lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Render(keyStr)
	}
	s += "\n\n"
	
	// Show inputs
	s += "Macro Details:\n\n"
	
	inputStyle := lipgloss.NewStyle().PaddingLeft(2)
	activeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
	
	// Name input
	label := "Name: "
	if m.activeInput == 0 {
		label = activeStyle.Render(label)
	}
	s += inputStyle.Render(label + m.nameInput.View()) + "\n\n"
	
	// Description input
	label = "Description: "
	if m.activeInput == 1 {
		label = activeStyle.Render(label)
	}
	s += inputStyle.Render(label + m.descInput.View()) + "\n\n"
	
	// Hotkey input
	label = "Hotkey: "
	if m.activeInput == 2 {
		label = activeStyle.Render(label)
	}
	s += inputStyle.Render(label + m.hotkeyInput.View()) + "\n\n"
	
	// Instructions
	s += statusStyle.Render("Tab/Shift+Tab: Navigate • Enter: Save (when on hotkey) • Esc: Cancel") + "\n"
	s += statusStyle.Render("Type any keys to record them") + "\n"
	
	return s
}

func (m model) viewHelp() string {
	s := titleStyle.Render("🎮 Macro Daemon TUI - Help") + "\n\n"
	s += m.help.View(keys) + "\n"
	s += "\nPress ? to return\n"
	return s
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Toggle, k.Edit, k.Delete},
		{k.Play, k.Record, k.New},
		{k.Help, k.Quit},
	}
}

func main() {
	m := newModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	
	// Run the program
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	
	// Clean up IPC connection
	if fm, ok := finalModel.(model); ok && fm.ipcClient != nil {
		fm.ipcClient.Close()
	}
}