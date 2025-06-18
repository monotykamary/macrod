// +build !darwin

package keylogger

import (
	"fmt"
	"log"
	"time"

	"github.com/monotykamary/macrod/pkg/models"
)

// Default implementation for non-macOS systems
type Keylogger struct {
	recording   bool
	currentKeys []models.KeyAction
	onKeyPress  func(key models.KeyAction)
	hotkeys     map[string]func()
}

func New() *Keylogger {
	log.Println("Using default keylogger (limited functionality)")
	return &Keylogger{
		hotkeys: make(map[string]func()),
	}
}

func (k *Keylogger) StartRecording(onKeyPress func(key models.KeyAction)) error {
	if k.recording {
		return fmt.Errorf("already recording")
	}

	k.recording = true
	k.currentKeys = []models.KeyAction{}
	k.onKeyPress = onKeyPress
	
	log.Println("Recording started (stub implementation)")
	return nil
}

func (k *Keylogger) StopRecording() []models.KeyAction {
	if !k.recording {
		return nil
	}
	
	k.recording = false
	keys := k.currentKeys
	k.currentKeys = nil
	
	log.Println("Recording stopped")
	return keys
}

func (k *Keylogger) AddRecordedKey(key string, modifiers []string) {
	if !k.recording {
		return
	}
	
	keyAction := models.KeyAction{
		Key:       key,
		Delay:     100 * time.Millisecond,
		Modifiers: modifiers,
	}
	
	k.currentKeys = append(k.currentKeys, keyAction)
	
	if k.onKeyPress != nil {
		k.onKeyPress(keyAction)
	}
}

func (k *Keylogger) PlaybackMacro(macro models.Macro) error {
	if !macro.Enabled {
		return fmt.Errorf("macro is disabled")
	}
	
	log.Printf("Playing back macro: %s (stub implementation)", macro.Name)
	return nil
}

func (k *Keylogger) RegisterHotkey(hotkey string, callback func()) error {
	k.hotkeys[hotkey] = callback
	log.Printf("Registered hotkey: %s (stub implementation)", hotkey)
	return nil
}