// +build darwin,cgo

package keylogger

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Carbon -framework ApplicationServices -framework CoreFoundation

#import <CoreFoundation/CoreFoundation.h>
#import <ApplicationServices/ApplicationServices.h>
#import <Carbon/Carbon.h>
#import <dispatch/dispatch.h>

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

// Callback function declaration
extern void goKeyCallback(int keyCode, int flags, int eventType);

// Global variables for the event tap
static CFMachPortRef eventTap = NULL;
static CFRunLoopSourceRef runLoopSource = NULL;
static dispatch_queue_t eventQueue = NULL;

// Key event callback
CGEventRef keyEventCallback(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *refcon) {
    if (type == kCGEventTapDisabledByTimeout || type == kCGEventTapDisabledByUserInput) {
        // Re-enable the event tap
        if (eventTap) {
            CGEventTapEnable(eventTap, true);
        }
        return event;
    }
    
    if (type != kCGEventKeyDown && type != kCGEventKeyUp && type != kCGEventFlagsChanged) {
        return event;
    }
    
    CGKeyCode keyCode = (CGKeyCode)CGEventGetIntegerValueField(event, kCGKeyboardEventKeycode);
    CGEventFlags flags = CGEventGetFlags(event);
    
    // Call Go callback
    goKeyCallback((int)keyCode, (int)flags, (int)type);
    
    // Pass through the event
    return event;
}

// Start capturing keys
int startKeyCapture() {
    // Check for accessibility permissions
    if (!AXIsProcessTrustedWithOptions(NULL)) {
        // Request permissions
        NSDictionary *options = @{(__bridge id)kAXTrustedCheckOptionPrompt: @YES};
        AXIsProcessTrustedWithOptions((__bridge CFDictionaryRef)options);
        return -1; // Not authorized
    }
    
    // Create event tap
    CGEventMask eventMask = (1 << kCGEventKeyDown) | (1 << kCGEventKeyUp) | (1 << kCGEventFlagsChanged);
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
void stopKeyCapture() {
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

// Check if we have accessibility permissions
int hasAccessibilityPermissions() {
    return AXIsProcessTrustedWithOptions(NULL) ? 1 : 0;
}

// Convert keycode to string (basic mapping)
// Using numeric values for macOS keycodes
const char* keycodeToString(int keyCode) {
    switch(keyCode) {
        // Letters
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
*/
import "C"
import (
	"fmt"
	"log"
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
	recordingChan chan models.KeyAction
	stopChan      chan bool
	runLoop       bool
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
	
	// Start the event tap
	result := C.startKeyCapture()
	if result != 0 {
		k.recording = false
		if result == -1 {
			return fmt.Errorf("accessibility permissions required")
		}
		return fmt.Errorf("failed to start key capture")
	}
	
	// Start processing recorded keys
	go k.processRecordedKeys()
	
	// Start the run loop in a separate goroutine
	go k.runEventLoop()
	
	log.Println("Started global key capture")
	return nil
}

func (k *Keylogger) StopRecording() []models.KeyAction {
	k.mu.Lock()
	defer k.mu.Unlock()

	if !k.recording {
		return nil
	}
	
	k.recording = false
	k.runLoop = false
	
	// Stop the event tap
	C.stopKeyCapture()
	
	// Signal stop
	select {
	case k.stopChan <- true:
	default:
	}
	
	keys := k.currentKeys
	k.currentKeys = nil
	
	log.Printf("Stopped global key capture - captured %d keys", len(keys))
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
func goKeyCallback(keyCode C.int, flags C.int, eventType C.int) {
	globalMutex.Lock()
	kl := globalKeylogger
	globalMutex.Unlock()

	if kl == nil || !kl.recording {
		return
	}

	// Only process key down events
	if eventType != C.kCGEventKeyDown {
		return
	}

	// Convert keycode to string
	keyStr := C.GoString(C.keycodeToString(keyCode))
	if keyStr == "unknown" {
		return
	}

	// Calculate delay
	currentTime := time.Now()
	delay := currentTime.Sub(kl.lastKeyTime)
	kl.lastKeyTime = currentTime

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

// PlaybackMacro remains the same
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
			switch mod {
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
	log.Printf("Registered hotkey: %s", hotkey)
	// TODO: Implement actual hotkey registration
	return nil
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