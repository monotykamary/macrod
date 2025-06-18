package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/monotykamary/macrod/internal/keylogger"
	"github.com/monotykamary/macrod/internal/storage"
	"github.com/monotykamary/macrod/pkg/models"
)

const socketPath = "/tmp/macrod.sock"

type Daemon struct {
	keylogger *keylogger.Keylogger
	storage   *storage.Storage
	macros    map[string]models.Macro
	mu        sync.RWMutex
	recording bool
}

func NewDaemon() *Daemon {
	return &Daemon{
		keylogger: keylogger.New(),
		storage:   storage.New(),
		macros:    make(map[string]models.Macro),
	}
}

func (d *Daemon) Start() error {
	log.Println("Starting macro daemon...")
	
	// Check for accessibility permissions on macOS
	if runtime.GOOS == "darwin" {
		log.Println("ℹ️  Checking accessibility permissions...")
		// We'll check when recording starts, but give a heads up
		log.Println("Note: Global key recording requires accessibility permissions.")
		log.Println("You may be prompted when recording starts.")
	}

	// Load existing macros
	if err := d.loadMacros(); err != nil {
		log.Printf("Failed to load macros: %v", err)
	}

	// Register hotkeys for all enabled macros
	d.registerAllHotkeys()

	// Start IPC server
	go d.startIPCServer()

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down daemon...")
	return nil
}

func (d *Daemon) loadMacros() error {
	macros, err := d.storage.LoadMacros()
	if err != nil {
		return err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	// If no macros exist, create some examples
	if len(macros) == 0 {
		exampleMacros := []models.Macro{
			{
				ID:          "example1",
				Name:        "Hello World",
				Description: "Types 'Hello World!'",
				Hotkey:      "ctrl+shift+1",
				Actions: []models.KeyAction{
					{Key: "h", Delay: 50 * time.Millisecond},
					{Key: "e", Delay: 50 * time.Millisecond},
					{Key: "l", Delay: 50 * time.Millisecond},
					{Key: "l", Delay: 50 * time.Millisecond},
					{Key: "o", Delay: 50 * time.Millisecond},
					{Key: "space", Delay: 50 * time.Millisecond},
					{Key: "w", Delay: 50 * time.Millisecond, Modifiers: []string{"shift"}},
					{Key: "o", Delay: 50 * time.Millisecond},
					{Key: "r", Delay: 50 * time.Millisecond},
					{Key: "l", Delay: 50 * time.Millisecond},
					{Key: "d", Delay: 50 * time.Millisecond},
					{Key: "1", Delay: 50 * time.Millisecond, Modifiers: []string{"shift"}},
				},
				Enabled:   true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:          "example2",
				Name:        "Quick Save",
				Description: "Performs Cmd+S to save",
				Hotkey:      "ctrl+shift+2",
				Actions: []models.KeyAction{
					{Key: "s", Modifiers: []string{"cmd"}},
				},
				Enabled:   false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		
		for _, macro := range exampleMacros {
			d.macros[macro.ID] = macro
		}
		
		// Save the examples
		d.saveMacros()
		log.Println("Created example macros")
	} else {
		for _, macro := range macros {
			d.macros[macro.ID] = macro
		}
	}

	return nil
}

func (d *Daemon) registerAllHotkeys() {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for _, macro := range d.macros {
		if macro.Enabled && macro.Hotkey != "" {
			macroID := macro.ID
			d.keylogger.RegisterHotkey(macro.Hotkey, func() {
				d.playbackMacro(macroID)
			})
		}
	}
}

func (d *Daemon) playbackMacro(macroID string) {
	d.mu.RLock()
	macro, exists := d.macros[macroID]
	d.mu.RUnlock()

	if !exists {
		log.Printf("Macro %s not found", macroID)
		return
	}

	if err := d.keylogger.PlaybackMacro(macro); err != nil {
		log.Printf("Failed to playback macro: %v", err)
	}
}

func (d *Daemon) startIPCServer() {
	// Remove existing socket
	os.Remove(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatal("Failed to create socket:", err)
	}
	defer listener.Close()

	log.Printf("IPC server listening on %s", socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go d.handleConnection(conn)
	}
}

func (d *Daemon) handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	var request map[string]interface{}
	if err := decoder.Decode(&request); err != nil {
		log.Printf("Failed to decode request: %v", err)
		return
	}

	command, ok := request["command"].(string)
	if !ok {
		encoder.Encode(map[string]string{"error": "invalid command"})
		return
	}

	switch command {
	case "status":
		d.handleStatus(encoder)
	case "list":
		d.handleList(encoder)
	case "toggle":
		d.handleToggle(request, encoder)
	case "delete":
		d.handleDelete(request, encoder)
	case "startRecording":
		d.handleStartRecording(encoder)
	case "stopRecording":
		d.handleStopRecording(request, encoder)
	case "addKey":
		d.handleAddKey(request, encoder)
	case "play":
		d.handlePlay(request, encoder)
	case "update":
		d.handleUpdate(request, encoder)
	case "getRecordingStatus":
		d.handleGetRecordingStatus(encoder)
	default:
		encoder.Encode(map[string]string{"error": "unknown command"})
	}
}

func (d *Daemon) handleStatus(encoder *json.Encoder) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	response := map[string]interface{}{
		"running":   true,
		"recording": d.recording,
		"macros":    len(d.macros),
	}
	encoder.Encode(response)
}

func (d *Daemon) handleList(encoder *json.Encoder) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	macroList := make([]models.Macro, 0, len(d.macros))
	for _, macro := range d.macros {
		macroList = append(macroList, macro)
	}

	encoder.Encode(map[string]interface{}{
		"macros": macroList,
	})
}

func (d *Daemon) handleToggle(request map[string]interface{}, encoder *json.Encoder) {
	macroID, ok := request["id"].(string)
	if !ok {
		encoder.Encode(map[string]string{"error": "missing macro id"})
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	macro, exists := d.macros[macroID]
	if !exists {
		encoder.Encode(map[string]string{"error": "macro not found"})
		return
	}

	// Unregister old hotkey if it was enabled
	if macro.Enabled && macro.Hotkey != "" {
		d.keylogger.UnregisterHotkey(macro.Hotkey)
	}
	
	macro.Enabled = !macro.Enabled
	d.macros[macroID] = macro

	// Register new hotkey if now enabled
	if macro.Enabled && macro.Hotkey != "" {
		macroIDCopy := macroID // Capture for closure
		d.keylogger.RegisterHotkey(macro.Hotkey, func() {
			d.playbackMacro(macroIDCopy)
		})
	}

	// Save to storage
	d.saveMacros()

	encoder.Encode(map[string]bool{"success": true})
}

func (d *Daemon) handleDelete(request map[string]interface{}, encoder *json.Encoder) {
	macroID, ok := request["id"].(string)
	if !ok {
		encoder.Encode(map[string]string{"error": "missing macro id"})
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	// Unregister hotkey if macro had one
	if macro, exists := d.macros[macroID]; exists {
		if macro.Enabled && macro.Hotkey != "" {
			d.keylogger.UnregisterHotkey(macro.Hotkey)
		}
	}

	delete(d.macros, macroID)

	// Save to storage
	d.saveMacros()

	encoder.Encode(map[string]bool{"success": true})
}

func (d *Daemon) handleStartRecording(encoder *json.Encoder) {
	if d.recording {
		encoder.Encode(map[string]string{"error": "already recording"})
		return
	}

	d.recording = true
	err := d.keylogger.StartRecording(func(key models.KeyAction) {
		// Callback for each key press during recording
		log.Printf("Recorded key: %s (modifiers: %v)", key.Key, key.Modifiers)
	})
	
	if err != nil {
		d.recording = false
		if err.Error() == "accessibility permissions required" || err.Error() == "accessibility permissions not granted" {
			log.Println("⚠️  Accessibility permissions required!")
			log.Println("Please grant accessibility permissions to macrod-daemon in:")
			log.Println("System Preferences → Security & Privacy → Privacy → Accessibility")
			encoder.Encode(map[string]string{"error": "Accessibility permissions required. Check System Preferences → Security & Privacy → Privacy → Accessibility"})
		} else {
			encoder.Encode(map[string]string{"error": err.Error()})
		}
		return
	}

	log.Println("✅ Started global key recording")
	encoder.Encode(map[string]bool{"success": true})
}

func (d *Daemon) handleStopRecording(request map[string]interface{}, encoder *json.Encoder) {
	if !d.recording {
		encoder.Encode(map[string]string{"error": "not recording"})
		return
	}

	keys := d.keylogger.StopRecording()
	d.recording = false

	// Extract macro details from request
	name, _ := request["name"].(string)
	description, _ := request["description"].(string)
	hotkey, _ := request["hotkey"].(string)

	// Create new macro
	macro := models.Macro{
		ID:          fmt.Sprintf("macro_%d", len(d.macros)+1),
		Name:        name,
		Description: description,
		Hotkey:      hotkey,
		Actions:     keys,
		Enabled:     true,
	}

	d.mu.Lock()
	d.macros[macro.ID] = macro
	d.mu.Unlock()

	// Save to storage
	d.saveMacros()

	// Register hotkey if provided
	if hotkey != "" {
		d.keylogger.RegisterHotkey(hotkey, func() {
			d.playbackMacro(macro.ID)
		})
	}

	encoder.Encode(map[string]interface{}{
		"success": true,
		"macro":   macro,
	})
}

func (d *Daemon) handleAddKey(request map[string]interface{}, encoder *json.Encoder) {
	if !d.recording {
		encoder.Encode(map[string]string{"error": "not recording"})
		return
	}
	
	key, ok := request["key"].(string)
	if !ok {
		encoder.Encode(map[string]string{"error": "missing key"})
		return
	}
	
	// Extract modifiers
	var modifiers []string
	if mods, ok := request["modifiers"].([]interface{}); ok {
		for _, mod := range mods {
			if modStr, ok := mod.(string); ok {
				modifiers = append(modifiers, modStr)
			}
		}
	}
	
	// Add the key to the keylogger
	d.keylogger.AddRecordedKey(key, modifiers)
	
	encoder.Encode(map[string]bool{"success": true})
}

func (d *Daemon) handlePlay(request map[string]interface{}, encoder *json.Encoder) {
	macroID, ok := request["id"].(string)
	if !ok {
		encoder.Encode(map[string]string{"error": "missing macro id"})
		return
	}
	
	// Play the macro in a goroutine to not block
	go d.playbackMacro(macroID)
	
	encoder.Encode(map[string]bool{"success": true})
}

func (d *Daemon) handleGetRecordingStatus(encoder *json.Encoder) {
	if !d.recording {
		encoder.Encode(map[string]interface{}{
			"recording": false,
			"keys":      []models.KeyAction{},
		})
		return
	}
	
	// Get current recorded keys from keylogger (without stopping)
	keys := d.keylogger.GetCurrentRecordedKeys()
	
	encoder.Encode(map[string]interface{}{
		"recording": true,
		"keys":      keys,
	})
}

func (d *Daemon) handleUpdate(request map[string]interface{}, encoder *json.Encoder) {
	macroData, ok := request["macro"].(map[string]interface{})
	if !ok {
		encoder.Encode(map[string]string{"error": "missing macro data"})
		return
	}
	
	// Convert the macro data
	data, err := json.Marshal(macroData)
	if err != nil {
		encoder.Encode(map[string]string{"error": "invalid macro data"})
		return
	}
	
	var macro models.Macro
	if err := json.Unmarshal(data, &macro); err != nil {
		encoder.Encode(map[string]string{"error": "failed to parse macro"})
		return
	}
	
	d.mu.Lock()
	defer d.mu.Unlock()
	
	// Check if macro exists
	oldMacro, exists := d.macros[macro.ID]
	if !exists {
		encoder.Encode(map[string]string{"error": "macro not found"})
		return
	}
	
	// Unregister old hotkey if it was enabled
	if oldMacro.Enabled && oldMacro.Hotkey != "" {
		d.keylogger.UnregisterHotkey(oldMacro.Hotkey)
	}
	
	// Update the macro
	d.macros[macro.ID] = macro
	
	// Register new hotkey if enabled
	if macro.Enabled && macro.Hotkey != "" {
		macroIDCopy := macro.ID // Capture for closure
		d.keylogger.RegisterHotkey(macro.Hotkey, func() {
			d.playbackMacro(macroIDCopy)
		})
	}
	
	// Save to storage
	d.saveMacros()
	
	encoder.Encode(map[string]bool{"success": true})
}

func (d *Daemon) saveMacros() {
	macroList := make([]models.Macro, 0, len(d.macros))
	for _, macro := range d.macros {
		macroList = append(macroList, macro)
	}

	if err := d.storage.SaveMacros(macroList); err != nil {
		log.Printf("Failed to save macros: %v", err)
	}
}

func main() {
	daemon := NewDaemon()
	if err := daemon.Start(); err != nil {
		log.Fatal(err)
	}
}