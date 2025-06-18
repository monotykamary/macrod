package keylogger

// Virtual key code constants
// These are based on Windows virtual key codes but work cross-platform with keybd_event
const (
	// Letters
	VK_A = 30
	VK_B = 48
	VK_C = 46
	VK_D = 32
	VK_E = 18
	VK_F = 33
	VK_G = 34
	VK_H = 35
	VK_I = 23
	VK_J = 36
	VK_K = 37
	VK_L = 38
	VK_M = 50
	VK_N = 49
	VK_O = 24
	VK_P = 25
	VK_Q = 16
	VK_R = 19
	VK_S = 31
	VK_T = 20
	VK_U = 22
	VK_V = 47
	VK_W = 17
	VK_X = 45
	VK_Y = 21
	VK_Z = 44
	
	// Numbers
	VK_0 = 11
	VK_1 = 2
	VK_2 = 3
	VK_3 = 4
	VK_4 = 5
	VK_5 = 6
	VK_6 = 7
	VK_7 = 8
	VK_8 = 9
	VK_9 = 10
	
	// Function keys
	VK_F1  = 59
	VK_F2  = 60
	VK_F3  = 61
	VK_F4  = 62
	VK_F5  = 63
	VK_F6  = 64
	VK_F7  = 65
	VK_F8  = 66
	VK_F9  = 67
	VK_F10 = 68
	VK_F11 = 87
	VK_F12 = 88
	
	// Special keys
	VK_ESC       = 1
	VK_BACKSPACE = 14
	VK_TAB       = 15
	VK_ENTER     = 28
	VK_SPACE     = 57
	VK_DELETE    = 0x2E + 0xFFF
	
	// Arrow keys (virtual keys)
	VK_LEFT  = 0x25 + 0xFFF
	VK_UP    = 0x26 + 0xFFF
	VK_RIGHT = 0x27 + 0xFFF
	VK_DOWN  = 0x28 + 0xFFF
	
	// Punctuation
	VK_MINUS      = 12
	VK_EQUAL      = 13
	VK_LEFTBRACE  = 26
	VK_RIGHTBRACE = 27
	VK_SEMICOLON  = 39
	VK_COMMA      = 51
	VK_DOT        = 52
	VK_SLASH      = 53
)