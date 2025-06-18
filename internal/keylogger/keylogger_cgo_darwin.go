//go:build darwin && cgo
// +build darwin,cgo

package keylogger

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Carbon -framework ApplicationServices -framework CoreFoundation

#ifndef MACROD_KEYLOGGER_H
#define MACROD_KEYLOGGER_H

#import <CoreFoundation/CoreFoundation.h>
#import <ApplicationServices/ApplicationServices.h>
#import <Carbon/Carbon.h>
#import <dispatch/dispatch.h>
#include <string.h>

// Define event types if not already defined
#ifndef kCGEventKeyDown
#define kCGEventKeyDown 10
#endif
#ifndef kCGEventKeyUp
#define kCGEventKeyUp 11
#endif
#ifndef kCGEventFlagsChanged
#define kCGEventFlagsChanged 12
#endif

// Define modifier flags
#ifndef kCGEventFlagMaskShift
#define kCGEventFlagMaskShift 0x00020000
#endif
#ifndef kCGEventFlagMaskControl
#define kCGEventFlagMaskControl 0x00040000
#endif
#ifndef kCGEventFlagMaskAlternate
#define kCGEventFlagMaskAlternate 0x00080000
#endif
#ifndef kCGEventFlagMaskCommand
#define kCGEventFlagMaskCommand 0x00100000
#endif

// Define keyboard event field
#ifndef kCGKeyboardEventAutorepeat
#define kCGKeyboardEventAutorepeat 3
#endif

// Callback function declaration
extern void goKeyCallback(int keyCode, int flags, int eventType, int isRepeat);

// Global variables for the event tap
static CFMachPortRef eventTap = NULL;
static CFRunLoopSourceRef runLoopSource = NULL;
static dispatch_queue_t eventQueue = NULL;

// Key event callback
static CGEventRef keyEventCallback(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *refcon) {
    if (type == kCGEventTapDisabledByTimeout || type == kCGEventTapDisabledByUserInput) {
        // Re-enable the event tap
        if (eventTap) {
            CGEventTapEnable(eventTap, true);
        }
        return event;
    }
    
    if (type != kCGEventKeyDown) {
        return event;
    }
    
    CGKeyCode keyCode = (CGKeyCode)CGEventGetIntegerValueField(event, kCGKeyboardEventKeycode);
    CGEventFlags flags = CGEventGetFlags(event);
    
    // Check if this is a key repeat
    int64_t keyRepeat = CGEventGetIntegerValueField(event, kCGKeyboardEventAutorepeat);
    
    // Call Go callback with repeat flag
    goKeyCallback((int)keyCode, (int)flags, (int)type, (int)keyRepeat);
    
    // Pass through the event
    return event;
}

// Start capturing keys
static int startKeyCapture() {
    // Check for accessibility permissions
    if (!AXIsProcessTrustedWithOptions(NULL)) {
        // Request permissions
        NSDictionary *options = @{(__bridge id)kAXTrustedCheckOptionPrompt: @YES};
        AXIsProcessTrustedWithOptions((__bridge CFDictionaryRef)options);
        return -1; // Not authorized
    }
    
    // Create event tap - only capture key down events for recording
    CGEventMask eventMask = (1 << kCGEventKeyDown);
    eventTap = CGEventTapCreate(kCGSessionEventTap,
                                kCGHeadInsertEventTap,
                                kCGEventTapOptionDefault,
                                eventMask,
                                keyEventCallback,
                                NULL);
    
    if (!eventTap) {
        return -2; // Failed to create event tap
    }
    
    // Create run loop source
    runLoopSource = CFMachPortCreateRunLoopSource(kCFAllocatorDefault, eventTap, 0);
    
    // Add to current run loop
    CFRunLoopAddSource(CFRunLoopGetCurrent(), runLoopSource, kCFRunLoopCommonModes);
    
    // Enable the event tap
    CGEventTapEnable(eventTap, true);
    
    // Create dispatch queue for event processing
    eventQueue = dispatch_queue_create("com.macrod.eventqueue", DISPATCH_QUEUE_SERIAL);
    
    return 0; // Success
}

// Stop capturing keys
static void stopKeyCapture() {
    if (eventTap) {
        CGEventTapEnable(eventTap, false);
        CFRelease(eventTap);
        eventTap = NULL;
    }
    
    if (runLoopSource) {
        CFRunLoopRemoveSource(CFRunLoopGetCurrent(), runLoopSource, kCFRunLoopCommonModes);
        CFRelease(runLoopSource);
        runLoopSource = NULL;
    }
    
    if (eventQueue) {
        dispatch_release(eventQueue);
        eventQueue = NULL;
    }
    
    // Stop the run loop
    CFRunLoopStop(CFRunLoopGetCurrent());
}

// Play a key event
static void playKeyEvent(int keyCode, int flags, int isKeyDown) {
    CGEventRef event;
    
    if (isKeyDown) {
        event = CGEventCreateKeyboardEvent(NULL, (CGKeyCode)keyCode, true);
    } else {
        event = CGEventCreateKeyboardEvent(NULL, (CGKeyCode)keyCode, false);
    }
    
    if (event) {
        // Set modifier flags
        CGEventSetFlags(event, (CGEventFlags)flags);
        
        // Post the event
        CGEventPost(kCGHIDEventTap, event);
        
        // Release the event
        CFRelease(event);
    }
}

// Get keycode from string
static int getKeycodeFromString(const char* key) {
    // Single character keys
    if (strlen(key) == 1) {
        char c = key[0];
        switch(c) {
            case 'a': return 0;
            case 'b': return 11;
            case 'c': return 8;
            case 'd': return 2;
            case 'e': return 14;
            case 'f': return 3;
            case 'g': return 5;
            case 'h': return 4;
            case 'i': return 34;
            case 'j': return 38;
            case 'k': return 40;
            case 'l': return 37;
            case 'm': return 46;
            case 'n': return 45;
            case 'o': return 31;
            case 'p': return 35;
            case 'q': return 12;
            case 'r': return 15;
            case 's': return 1;
            case 't': return 17;
            case 'u': return 32;
            case 'v': return 9;
            case 'w': return 13;
            case 'x': return 7;
            case 'y': return 16;
            case 'z': return 6;
            
            case '0': return 29;
            case '1': return 18;
            case '2': return 19;
            case '3': return 20;
            case '4': return 21;
            case '5': return 23;
            case '6': return 22;
            case '7': return 26;
            case '8': return 28;
            case '9': return 25;
        }
    }
    
    // Punctuation and symbols
    if (strcmp(key, "-") == 0) return 27;
    if (strcmp(key, "=") == 0) return 24;
    if (strcmp(key, "[") == 0) return 33;
    if (strcmp(key, "]") == 0) return 30;
    if (strcmp(key, "\\") == 0) return 42;
    if (strcmp(key, ";") == 0) return 41;
    if (strcmp(key, "'") == 0) return 39;
    if (strcmp(key, ",") == 0) return 43;
    if (strcmp(key, ".") == 0) return 47;
    if (strcmp(key, "/") == 0) return 44;
    if (strcmp(key, "`") == 0) return 50;
    
    // Special keys
    if (strcmp(key, "space") == 0) return 49;
    if (strcmp(key, "enter") == 0) return 36;
    if (strcmp(key, "tab") == 0) return 48;
    if (strcmp(key, "escape") == 0 || strcmp(key, "esc") == 0) return 53;
    if (strcmp(key, "backspace") == 0) return 51;
    if (strcmp(key, "delete") == 0) return 117;
    
    // Arrow keys
    if (strcmp(key, "up") == 0) return 126;
    if (strcmp(key, "down") == 0) return 125;
    if (strcmp(key, "left") == 0) return 123;
    if (strcmp(key, "right") == 0) return 124;
    
    // Function keys
    if (strcmp(key, "f1") == 0) return 122;
    if (strcmp(key, "f2") == 0) return 120;
    if (strcmp(key, "f3") == 0) return 99;
    if (strcmp(key, "f4") == 0) return 118;
    if (strcmp(key, "f5") == 0) return 96;
    if (strcmp(key, "f6") == 0) return 97;
    if (strcmp(key, "f7") == 0) return 98;
    if (strcmp(key, "f8") == 0) return 100;
    if (strcmp(key, "f9") == 0) return 101;
    if (strcmp(key, "f10") == 0) return 109;
    if (strcmp(key, "f11") == 0) return 103;
    if (strcmp(key, "f12") == 0) return 111;
    
    return -1; // Unknown key
}

// Check if we have accessibility permissions
static int hasAccessibilityPermissions() {
    return AXIsProcessTrustedWithOptions(NULL) ? 1 : 0;
}

// Convert keycode to string (basic mapping)
// Using numeric values for macOS keycodes
static const char* keycodeToString(int keyCode) {
    switch(keyCode) {
        // Letters (always return lowercase, shift is handled separately)
        case 0: return "a";
        case 11: return "b";
        case 8: return "c";
        case 2: return "d";
        case 14: return "e";
        case 3: return "f";
        case 5: return "g";
        case 4: return "h";
        case 34: return "i";
        case 38: return "j";
        case 40: return "k";
        case 37: return "l";
        case 46: return "m";
        case 45: return "n";
        case 31: return "o";
        case 35: return "p";
        case 12: return "q";
        case 15: return "r";
        case 1: return "s";
        case 17: return "t";
        case 32: return "u";
        case 9: return "v";
        case 13: return "w";
        case 7: return "x";
        case 16: return "y";
        case 6: return "z";
        
        // Numbers
        case 29: return "0";
        case 18: return "1";
        case 19: return "2";
        case 20: return "3";
        case 21: return "4";
        case 23: return "5";
        case 22: return "6";
        case 26: return "7";
        case 28: return "8";
        case 25: return "9";
        
        // Special keys
        case 49: return "space";
        case 36: return "enter";
        case 48: return "tab";
        case 53: return "escape";
        case 51: return "backspace";
        
        // Arrow keys
        case 126: return "up";
        case 125: return "down";
        case 123: return "left";
        case 124: return "right";
        
        // Function keys
        case 122: return "f1";
        case 120: return "f2";
        case 99: return "f3";
        case 118: return "f4";
        case 96: return "f5";
        case 97: return "f6";
        case 98: return "f7";
        case 100: return "f8";
        case 101: return "f9";
        case 109: return "f10";
        case 103: return "f11";
        case 111: return "f12";
        
        default: return "unknown";
    }
}

#endif // MACROD_KEYLOGGER_H
*/
import "C"
import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/micmonay/keybd_event"
	"github.com/monotykamary/macrod/pkg/models"
)

// Global keylogger instance for C callbacks
var globalKeylogger *Keylogger
var globalMutex sync.Mutex

// Keylogger implementation using CGEventTap
type Keylogger struct {
	recording     bool
	currentKeys   []models.KeyAction
	onKeyPress    func(key models.KeyAction)
	hotkeys       map[string]func()
	kb            keybd_event.KeyBonding
	lastKeyTime   time.Time
	lastKeyCode   int     // Track last key to detect repeats
	recordingChan chan models.KeyAction
	stopChan      chan bool
	runLoop       bool
	monitoring    bool  // True when monitoring for hotkeys
	mu            sync.Mutex
}

func New() *Keylogger {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		log.Printf("Failed to create key bonding: %v", err)
	}

	kl := &Keylogger{
		hotkeys:       make(map[string]func()),
		kb:            kb,
		recordingChan: make(chan models.KeyAction, 100),
		stopChan:      make(chan bool),
	}

	// Set global instance
	globalMutex.Lock()
	globalKeylogger = kl
	globalMutex.Unlock()

	return kl
}

func (k *Keylogger) StartRecording(onKeyPress func(key models.KeyAction)) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.recording {
		return fmt.Errorf("already recording")
	}

	// Check permissions
	if C.hasAccessibilityPermissions() == 0 {
		return fmt.Errorf("accessibility permissions not granted")
	}

	k.recording = true
	k.currentKeys = []models.KeyAction{}
	k.onKeyPress = onKeyPress
	k.lastKeyTime = time.Now()
	
	// Only start event tap if not already monitoring
	if !k.monitoring && !k.runLoop {
		// Start the event tap
		result := C.startKeyCapture()
		if result != 0 {
			k.recording = false
			if result == -1 {
				return fmt.Errorf("accessibility permissions required")
			}
			return fmt.Errorf("failed to start key capture")
		}
		
		// Start the run loop in a separate goroutine
		k.runLoop = true
		go k.runEventLoop()
	}
	
	// Start processing recorded keys
	go k.processRecordedKeys()
	
	log.Printf("Started recording (event tap already running: %v)", k.monitoring || k.runLoop)
	return nil
}

func (k *Keylogger) PauseRecording() {
	k.mu.Lock()
	defer k.mu.Unlock()
	
	if !k.recording {
		return
	}
	
	// Mark recording as paused
	k.recording = false
	
	// Don't stop the event tap if we're monitoring for hotkeys
	if !k.monitoring {
		k.runLoop = false
		C.stopKeyCapture()
	}
	
	log.Println("Paused key recording")
}

func (k *Keylogger) StopRecording() []models.KeyAction {
	k.mu.Lock()
	defer k.mu.Unlock()

	// Even if recording is paused, we still want to return the captured keys
	wasRecording := k.recording || len(k.currentKeys) > 0
	if !wasRecording {
		return nil
	}
	
	k.recording = false
	
	// Only stop the event tap if we're not monitoring for hotkeys
	if !k.monitoring && k.runLoop {
		k.runLoop = false
		C.stopKeyCapture()
	}
	
	// Signal stop
	select {
	case k.stopChan <- true:
	default:
	}
	
	// Filter out the escape key if it's the last key
	keys := make([]models.KeyAction, 0, len(k.currentKeys))
	for i, key := range k.currentKeys {
		// Skip the last key if it's escape (the stop trigger)
		if i == len(k.currentKeys)-1 && key.Key == "escape" {
			continue
		}
		keys = append(keys, key)
	}
	k.currentKeys = nil
	
	log.Printf("Stopped global key capture - captured %d keys", len(keys))
	return keys
}

func (k *Keylogger) GetCurrentRecordedKeys() []models.KeyAction {
	k.mu.Lock()
	defer k.mu.Unlock()
	
	// Return a copy of current keys (excluding the last one if it's escape)
	keys := make([]models.KeyAction, 0, len(k.currentKeys))
	for i, key := range k.currentKeys {
		// Skip the last key if it's escape (the stop trigger)
		if i == len(k.currentKeys)-1 && key.Key == "escape" {
			continue
		}
		keys = append(keys, key)
	}
	return keys
}

func (k *Keylogger) runEventLoop() {
	k.runLoop = true
	for k.runLoop {
		// Run the event loop
		C.CFRunLoopRun()
		if !k.runLoop {
			break
		}
	}
}

func (k *Keylogger) processRecordedKeys() {
	for {
		select {
		case keyAction := <-k.recordingChan:
			k.mu.Lock()
			if k.recording {
				k.currentKeys = append(k.currentKeys, keyAction)
				if k.onKeyPress != nil {
					k.onKeyPress(keyAction)
				}
			}
			k.mu.Unlock()
		case <-k.stopChan:
			return
		}
	}
}

// Export for C callback
//export goKeyCallback
func goKeyCallback(keyCode C.int, flags C.int, eventType C.int, isRepeat C.int) {
	globalMutex.Lock()
	kl := globalKeylogger
	globalMutex.Unlock()

	if kl == nil {
		return
	}

	// Only process key down events (ignore key up and modifier changes)
	if eventType != C.kCGEventKeyDown {
		return
	}
	
	// Skip key repeats when recording
	if kl.recording && isRepeat != 0 {
		log.Printf("Skipping key repeat: keyCode=%d", keyCode)
		return
	}

	// Convert keycode to string
	keyStr := C.GoString(C.keycodeToString(keyCode))
	if keyStr == "unknown" {
		return
	}

	// Extract modifiers from flags
	modifiers := []string{}
	if flags&C.kCGEventFlagMaskShift != 0 {
		modifiers = append(modifiers, "shift")
	}
	if flags&C.kCGEventFlagMaskControl != 0 {
		modifiers = append(modifiers, "ctrl")
	}
	if flags&C.kCGEventFlagMaskAlternate != 0 {
		modifiers = append(modifiers, "alt")
	}
	if flags&C.kCGEventFlagMaskCommand != 0 {
		modifiers = append(modifiers, "cmd")
	}

	// Check for hotkey matches when monitoring
	if kl.monitoring && !kl.recording {
		hotkeyStr := buildHotkeyString(keyStr, modifiers)
		kl.mu.Lock()
		if callback, exists := kl.hotkeys[hotkeyStr]; exists {
			kl.mu.Unlock()
			// Execute callback in a goroutine to avoid blocking
			go callback()
			return
		}
		kl.mu.Unlock()
	}

	// Handle recording
	if kl.recording {
		// Calculate delay
		currentTime := time.Now()
		delay := currentTime.Sub(kl.lastKeyTime)
		kl.lastKeyTime = currentTime

		keyAction := models.KeyAction{
			Key:       keyStr,
			Delay:     delay,
			Modifiers: modifiers,
		}

		select {
		case kl.recordingChan <- keyAction:
		default:
			// Channel full, skip
		}
	}
}

// buildHotkeyString creates a normalized hotkey string from key and modifiers
func buildHotkeyString(key string, modifiers []string) string {
	// Sort modifiers for consistent ordering
	sortedMods := make([]string, len(modifiers))
	copy(sortedMods, modifiers)
	
	// Custom sort order: ctrl, alt, shift, cmd
	modOrder := map[string]int{"ctrl": 0, "alt": 1, "shift": 2, "cmd": 3}
	for i := 0; i < len(sortedMods)-1; i++ {
		for j := i + 1; j < len(sortedMods); j++ {
			if modOrder[sortedMods[i]] > modOrder[sortedMods[j]] {
				sortedMods[i], sortedMods[j] = sortedMods[j], sortedMods[i]
			}
		}
	}
	
	// Build the hotkey string
	if len(sortedMods) > 0 {
		return strings.Join(sortedMods, "+") + "+" + key
	}
	return key
}

// Keep the existing AddRecordedKey for manual recording via TUI
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
	log.Printf("Manually recorded: %s (delay: %v, modifiers: %v)", key, delay, modifiers)
}

// PlaybackMacro using native CGEventPost
func (k *Keylogger) PlaybackMacro(macro models.Macro) error {
	if !macro.Enabled {
		return fmt.Errorf("macro is disabled")
	}
	
	// Default speed multiplier is 1.0
	speedMultiplier := macro.SpeedMultiplier
	if speedMultiplier == 0 {
		speedMultiplier = 1.0
	}
	
	log.Printf("Playing back macro: %s (%d actions, speed: %.1fx)", macro.Name, len(macro.Actions), speedMultiplier)
	
	for i, action := range macro.Actions {
		// Wait for the specified delay (except for the first action)
		if i > 0 && action.Delay > 0 {
			// Apply speed multiplier to delay
			delay := time.Duration(float32(action.Delay) / speedMultiplier)
			
			// For navigation keys, ensure minimum responsiveness
			if isNavigationKey(action.Key) && delay > 20*time.Millisecond {
				delay = 20 * time.Millisecond
			} else if delay < 1*time.Millisecond {
				// Ensure minimum 1ms delay to prevent issues
				delay = 1 * time.Millisecond
			}
			
			time.Sleep(delay)
		}
		
		// Get the keycode
		keyCode := C.getKeycodeFromString(C.CString(action.Key))
		if keyCode < 0 {
			log.Printf("Unknown key: %s", action.Key)
			continue
		}
		
		// Build modifier flags
		flags := 0
		for _, mod := range action.Modifiers {
			switch mod {
			case "ctrl", "control":
				flags |= int(C.kCGEventFlagMaskControl)
			case "alt", "option":
				flags |= int(C.kCGEventFlagMaskAlternate)
			case "shift":
				flags |= int(C.kCGEventFlagMaskShift)
			case "cmd", "command", "super":
				flags |= int(C.kCGEventFlagMaskCommand)
			}
		}
		
		// Send key down event
		C.playKeyEvent(keyCode, C.int(flags), 1)
		
		// Minimal delay between press and release (1ms)
		time.Sleep(1 * time.Millisecond)
		
		// Send key up event
		C.playKeyEvent(keyCode, C.int(flags), 0)
		
		// Only add post-release delay for non-navigation keys
		if !isNavigationKey(action.Key) {
			time.Sleep(2 * time.Millisecond)
		}
	}
	
	return nil
}

// isNavigationKey checks if a key is a navigation key that should have reduced delays
func isNavigationKey(key string) bool {
	switch key {
	case "up", "down", "left", "right", 
	     "home", "end", "pageup", "pagedown",
	     "tab", "escape", "esc":
		return true
	default:
		return false
	}
}

func (k *Keylogger) RegisterHotkey(hotkey string, callback func()) error {
	k.mu.Lock()
	defer k.mu.Unlock()
	
	k.hotkeys[hotkey] = callback
	log.Printf("Registered hotkey: %s", hotkey)
	
	// Start monitoring if not already started
	if !k.monitoring && len(k.hotkeys) > 0 {
		go k.StartHotkeyMonitoring()
	}
	
	return nil
}

func (k *Keylogger) UnregisterHotkey(hotkey string) {
	k.mu.Lock()
	delete(k.hotkeys, hotkey)
	log.Printf("Unregistered hotkey: %s", hotkey)
	
	// Check if we should stop monitoring
	shouldStop := k.monitoring && len(k.hotkeys) == 0
	k.mu.Unlock()
	
	// Stop monitoring if no hotkeys left (without holding the lock)
	if shouldStop {
		k.StopHotkeyMonitoring()
	}
}

func (k *Keylogger) StartHotkeyMonitoring() error {
	k.mu.Lock()
	if k.monitoring {
		k.mu.Unlock()
		return fmt.Errorf("already monitoring")
	}
	k.monitoring = true
	k.mu.Unlock()
	
	// Check permissions
	if C.hasAccessibilityPermissions() == 0 {
		k.monitoring = false
		return fmt.Errorf("accessibility permissions not granted")
	}
	
	// Start the event tap
	result := C.startKeyCapture()
	if result != 0 {
		k.monitoring = false
		if result == -1 {
			return fmt.Errorf("accessibility permissions required")
		}
		return fmt.Errorf("failed to start key capture")
	}
	
	// Start the run loop
	go k.runEventLoop()
	
	log.Println("Started hotkey monitoring")
	return nil
}

func (k *Keylogger) StopHotkeyMonitoring() {
	k.mu.Lock()
	defer k.mu.Unlock()
	
	if !k.monitoring {
		return
	}
	
	k.monitoring = false
	k.runLoop = false
	
	// Stop the event tap
	C.stopKeyCapture()
	
	log.Println("Stopped hotkey monitoring")
}

// getKeyCode converts a key string to keybd_event keycode
func (k *Keylogger) getKeyCode(key string) int {
	// Use the same mapping from keycode_constants.go
	switch key {
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
	case "up": return VK_UP
	case "down": return VK_DOWN
	case "left": return VK_LEFT
	case "right": return VK_RIGHT
	
	// Special keys
	case "space": return VK_SPACE
	case "enter": return VK_ENTER
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
	
	default:
		return -1
	}
}