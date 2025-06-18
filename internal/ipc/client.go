package ipc

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/monotykamary/macrod/pkg/models"
)

const socketPath = "/tmp/macrod.sock"

type Client struct {
	// We'll create a new connection for each request to avoid broken pipe issues
}

func NewClient() (*Client, error) {
	// Just check if the daemon is reachable
	conn, err := net.DialTimeout("unix", socketPath, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to daemon: %w", err)
	}
	conn.Close()
	
	return &Client{}, nil
}

func (c *Client) Close() error {
	// Nothing to close anymore since we create connections per request
	return nil
}

func (c *Client) sendRequest(request map[string]interface{}) (map[string]interface{}, error) {
	// Create a new connection for each request
	conn, err := net.DialTimeout("unix", socketPath, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to daemon: %w", err)
	}
	defer conn.Close()
	
	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)
	
	if err := encoder.Encode(request); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	
	var response map[string]interface{}
	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if errMsg, ok := response["error"].(string); ok {
		return nil, fmt.Errorf("daemon error: %s", errMsg)
	}
	
	return response, nil
}

func (c *Client) GetStatus() (bool, bool, int, error) {
	resp, err := c.sendRequest(map[string]interface{}{
		"command": "status",
	})
	if err != nil {
		return false, false, 0, err
	}
	
	running, _ := resp["running"].(bool)
	recording, _ := resp["recording"].(bool)
	macrosFloat, _ := resp["macros"].(float64)
	macros := int(macrosFloat)
	
	return running, recording, macros, nil
}

func (c *Client) ListMacros() ([]models.Macro, error) {
	resp, err := c.sendRequest(map[string]interface{}{
		"command": "list",
	})
	if err != nil {
		return nil, err
	}
	
	// Convert response to macros
	macrosData, ok := resp["macros"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}
	
	// Re-encode and decode to get proper type
	data, err := json.Marshal(macrosData)
	if err != nil {
		return nil, err
	}
	
	var macros []models.Macro
	if err := json.Unmarshal(data, &macros); err != nil {
		return nil, err
	}
	
	return macros, nil
}

func (c *Client) ToggleMacro(macroID string) error {
	_, err := c.sendRequest(map[string]interface{}{
		"command": "toggle",
		"id":      macroID,
	})
	return err
}

func (c *Client) DeleteMacro(macroID string) error {
	_, err := c.sendRequest(map[string]interface{}{
		"command": "delete",
		"id":      macroID,
	})
	return err
}

func (c *Client) UpdateMacro(macro *models.Macro) error {
	_, err := c.sendRequest(map[string]interface{}{
		"command": "update",
		"macro":   macro,
	})
	return err
}

func (c *Client) StartRecording() error {
	_, err := c.sendRequest(map[string]interface{}{
		"command": "startRecording",
	})
	return err
}

func (c *Client) PauseRecording() error {
	_, err := c.sendRequest(map[string]interface{}{
		"command": "pauseRecording",
	})
	return err
}

func (c *Client) GetRecordingStatus() ([]models.KeyAction, error) {
	resp, err := c.sendRequest(map[string]interface{}{
		"command": "getRecordingStatus",
	})
	if err != nil {
		return nil, err
	}
	
	// Extract keys from response
	keysData, ok := resp["keys"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}
	
	// Re-encode and decode to get proper type
	data, err := json.Marshal(keysData)
	if err != nil {
		return nil, err
	}
	
	var keys []models.KeyAction
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, err
	}
	
	return keys, nil
}

func (c *Client) AddRecordedKey(key string, modifiers []string) error {
	_, err := c.sendRequest(map[string]interface{}{
		"command":   "addKey",
		"key":       key,
		"modifiers": modifiers,
	})
	return err
}

func (c *Client) StopRecording(name, description, hotkey string) (*models.Macro, error) {
	resp, err := c.sendRequest(map[string]interface{}{
		"command":     "stopRecording",
		"name":        name,
		"description": description,
		"hotkey":      hotkey,
	})
	if err != nil {
		return nil, err
	}
	
	// Convert macro from response
	macroData, ok := resp["macro"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}
	
	// Re-encode and decode to get proper type
	data, err := json.Marshal(macroData)
	if err != nil {
		return nil, err
	}
	
	var macro models.Macro
	if err := json.Unmarshal(data, &macro); err != nil {
		return nil, err
	}
	
	return &macro, nil
}

func (c *Client) PlayMacro(macroID string) error {
	_, err := c.sendRequest(map[string]interface{}{
		"command": "play",
		"id":      macroID,
	})
	return err
}

// Convenience function to check if daemon is running
func IsDaemonRunning() bool {
	client, err := NewClient()
	if err != nil {
		return false
	}
	defer client.Close()
	
	running, _, _, err := client.GetStatus()
	return err == nil && running
}