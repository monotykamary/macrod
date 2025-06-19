package models

import (
	"time"
)

type Macro struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Hotkey       string      `json:"hotkey"`
	Actions      []KeyAction `json:"actions"`
	Enabled      bool        `json:"enabled"`
	SpeedMultiplier float32  `json:"speed_multiplier,omitempty"` // 1.0 = normal, 2.0 = 2x speed, 0.5 = half speed
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type KeyAction struct {
	Key       string        `json:"key"`
	Delay     time.Duration `json:"delay"`
	Modifiers []string      `json:"modifiers"`
}

type MacroState struct {
	Macros   []Macro `json:"macros"`
	Recording bool    `json:"recording"`
	DaemonRunning bool `json:"daemon_running"`
}