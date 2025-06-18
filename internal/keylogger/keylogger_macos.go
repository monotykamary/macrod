//go:build darwin && !cgo
// +build darwin,!cgo

package keylogger

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/micmonay/keybd_event"
	"github.com/monotykamary/macrod/pkg/models"
)

// Keylogger uses keybd_event for playback and manual recording for now
type Keylogger struct {
	recording     bool
	currentKeys   []models.KeyAction
	onKeyPress    func(key models.KeyAction)
	hotkeys       map[string]func()
	kb            keybd_event.KeyBonding
	lastKeyTime   time.Time
	recordingChan chan models.KeyAction
	stopChan      chan bool
}

func New() *Keylogger {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		log.Printf("Failed to create key bonding: %v", err)
	}

	return &Keylogger{
		hotkeys:       make(map[string]func()),
		kb:            kb,
		recordingChan: make(chan models.KeyAction, 100),
		stopChan:      make(chan bool),
	}
}

func (k *Keylogger) StartRecording(onKeyPress func(key models.KeyAction)) error {
	if k.recording {
		return fmt.Errorf("already recording")
	}

	k.recording = true
	k.currentKeys = []models.KeyAction{}
	k.onKeyPress = onKeyPress
	k.lastKeyTime = time.Now()
	
	log.Println("Recording started - press keys in the terminal")
	
	// Start a goroutine to process recorded keys
	go k.processRecordedKeys()
	
	return nil
}

func (k *Keylogger) processRecordedKeys() {
	for {
		select {
		case keyAction := <-k.recordingChan:
			if k.recording {
				k.currentKeys = append(k.currentKeys, keyAction)
				if k.onKeyPress != nil {
					k.onKeyPress(keyAction)
				}
			}
		case <-k.stopChan:
			return
		}
	}
}

func (k *Keylogger) StopRecording() []models.KeyAction {
	if !k.recording {
		return nil
	}
	
	k.recording = false
	k.stopChan <- true
	keys := k.currentKeys
	k.currentKeys = nil
	
	log.Printf("Recording stopped - captured %d keys", len(keys))
	return keys
}

func (k *Keylogger) AddRecordedKey(key string, modifiers []string) {
	if !k.recording {
		return
	}
	
	currentTime := time.Now()
	delay := currentTime.Sub(k.lastKeyTime)
	k.lastKeyTime = currentTime
	
	keyAction := models.KeyAction{
		Key:       key,
		Delay:     delay,
		Modifiers: modifiers,
	}
	
	k.recordingChan <- keyAction
	log.Printf("Recorded: %s (delay: %v, modifiers: %v)", key, delay, modifiers)
}

func (k *Keylogger) PlaybackMacro(macro models.Macro) error {
	if !macro.Enabled {
		return fmt.Errorf("macro is disabled")
	}
	
	log.Printf("Playing back macro: %s (%d actions)", macro.Name, len(macro.Actions))
	
	for i, action := range macro.Actions {
		// Wait for the specified delay (except for the first action)
		if i > 0 && action.Delay > 0 {
			time.Sleep(action.Delay)
		}
		
		// Convert key string to keycode
		keyCode := k.getKeyCode(action.Key)
		if keyCode < 0 {
			log.Printf("Unknown key: %s", action.Key)
			continue
		}
		
		k.kb.Clear()
		
		// Set modifiers
		for _, mod := range action.Modifiers {
			switch strings.ToLower(mod) {
			case "ctrl", "control":
				k.kb.HasCTRL(true)
			case "alt", "option":
				k.kb.HasALT(true)
			case "shift":
				k.kb.HasSHIFT(true)
			case "cmd", "command", "super":
				k.kb.HasSuper(true)
			}
		}
		
		// Set the key
		k.kb.SetKeys(keyCode)
		
		// Press and release
		err := k.kb.Launching()
		if err != nil {
			log.Printf("Failed to send key %s: %v", action.Key, err)
		}
		
		// Small delay between press and release
		time.Sleep(10 * time.Millisecond)
	}
	
	return nil
}

func (k *Keylogger) RegisterHotkey(hotkey string, callback func()) error {
	k.hotkeys[hotkey] = callback
	log.Printf("Registered hotkey: %s (manual trigger required for now)", hotkey)
	
	// TODO: Implement actual hotkey registration using CGEventTap or similar
	// For now, hotkeys need to be triggered manually through the TUI
	
	return nil
}

// getKeyCode converts a key string to keybd_event keycode
func (k *Keylogger) getKeyCode(key string) int {
	switch strings.ToLower(key) {
	// Letters
	case "a": return VK_A
	case "b": return VK_B
	case "c": return VK_C
	case "d": return VK_D
	case "e": return VK_E
	case "f": return VK_F
	case "g": return VK_G
	case "h": return VK_H
	case "i": return VK_I
	case "j": return VK_J
	case "k": return VK_K
	case "l": return VK_L
	case "m": return VK_M
	case "n": return VK_N
	case "o": return VK_O
	case "p": return VK_P
	case "q": return VK_Q
	case "r": return VK_R
	case "s": return VK_S
	case "t": return VK_T
	case "u": return VK_U
	case "v": return VK_V
	case "w": return VK_W
	case "x": return VK_X
	case "y": return VK_Y
	case "z": return VK_Z
	
	// Numbers
	case "0": return VK_0
	case "1": return VK_1
	case "2": return VK_2
	case "3": return VK_3
	case "4": return VK_4
	case "5": return VK_5
	case "6": return VK_6
	case "7": return VK_7
	case "8": return VK_8
	case "9": return VK_9
	
	// Arrow keys
	case "up", "arrow_up": return VK_UP
	case "down", "arrow_down": return VK_DOWN
	case "left", "arrow_left": return VK_LEFT
	case "right", "arrow_right": return VK_RIGHT
	
	// Special keys
	case "space", " ": return VK_SPACE
	case "enter", "return": return VK_ENTER
	case "tab": return VK_TAB
	case "escape", "esc": return VK_ESC
	case "backspace": return VK_BACKSPACE
	case "delete": return VK_DELETE
	
	// Function keys
	case "f1": return VK_F1
	case "f2": return VK_F2
	case "f3": return VK_F3
	case "f4": return VK_F4
	case "f5": return VK_F5
	case "f6": return VK_F6
	case "f7": return VK_F7
	case "f8": return VK_F8
	case "f9": return VK_F9
	case "f10": return VK_F10
	case "f11": return VK_F11
	case "f12": return VK_F12
	
	// Punctuation
	case ".": return VK_DOT
	case ",": return VK_COMMA
	case ";": return VK_SEMICOLON
	case "/": return VK_SLASH
	case "[": return VK_LEFTBRACE
	case "]": return VK_RIGHTBRACE
	case "-": return VK_MINUS
	case "=": return VK_EQUAL
	
	default:
		return -1
	}
}