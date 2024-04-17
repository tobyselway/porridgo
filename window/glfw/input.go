package glfw

import (
	"porridgo/window"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var keyMapping = map[glfw.Key]window.Key{
	glfw.KeyUnknown:      window.KeyUnknown,
	glfw.KeySpace:        window.KeySpace,
	glfw.KeyApostrophe:   window.KeyApostrophe,
	glfw.KeyComma:        window.KeyComma,
	glfw.KeyMinus:        window.KeyMinus,
	glfw.KeyPeriod:       window.KeyPeriod,
	glfw.KeySlash:        window.KeySlash,
	glfw.Key0:            window.Key0,
	glfw.Key1:            window.Key1,
	glfw.Key2:            window.Key2,
	glfw.Key3:            window.Key3,
	glfw.Key4:            window.Key4,
	glfw.Key5:            window.Key5,
	glfw.Key6:            window.Key6,
	glfw.Key7:            window.Key7,
	glfw.Key8:            window.Key8,
	glfw.Key9:            window.Key9,
	glfw.KeySemicolon:    window.KeySemicolon,
	glfw.KeyEqual:        window.KeyEqual,
	glfw.KeyA:            window.KeyA,
	glfw.KeyB:            window.KeyB,
	glfw.KeyC:            window.KeyC,
	glfw.KeyD:            window.KeyD,
	glfw.KeyE:            window.KeyE,
	glfw.KeyF:            window.KeyF,
	glfw.KeyG:            window.KeyG,
	glfw.KeyH:            window.KeyH,
	glfw.KeyI:            window.KeyI,
	glfw.KeyJ:            window.KeyJ,
	glfw.KeyK:            window.KeyK,
	glfw.KeyL:            window.KeyL,
	glfw.KeyM:            window.KeyM,
	glfw.KeyN:            window.KeyN,
	glfw.KeyO:            window.KeyO,
	glfw.KeyP:            window.KeyP,
	glfw.KeyQ:            window.KeyQ,
	glfw.KeyR:            window.KeyR,
	glfw.KeyS:            window.KeyS,
	glfw.KeyT:            window.KeyT,
	glfw.KeyU:            window.KeyU,
	glfw.KeyV:            window.KeyV,
	glfw.KeyW:            window.KeyW,
	glfw.KeyX:            window.KeyX,
	glfw.KeyY:            window.KeyY,
	glfw.KeyZ:            window.KeyZ,
	glfw.KeyLeftBracket:  window.KeyLeftBracket,
	glfw.KeyBackslash:    window.KeyBackslash,
	glfw.KeyRightBracket: window.KeyRightBracket,
	glfw.KeyGraveAccent:  window.KeyGraveAccent,
	glfw.KeyWorld1:       window.KeyWorld1,
	glfw.KeyWorld2:       window.KeyWorld2,
	glfw.KeyEscape:       window.KeyEscape,
	glfw.KeyEnter:        window.KeyEnter,
	glfw.KeyTab:          window.KeyTab,
	glfw.KeyBackspace:    window.KeyBackspace,
	glfw.KeyInsert:       window.KeyInsert,
	glfw.KeyDelete:       window.KeyDelete,
	glfw.KeyRight:        window.KeyRight,
	glfw.KeyLeft:         window.KeyLeft,
	glfw.KeyDown:         window.KeyDown,
	glfw.KeyUp:           window.KeyUp,
	glfw.KeyPageUp:       window.KeyPageUp,
	glfw.KeyPageDown:     window.KeyPageDown,
	glfw.KeyHome:         window.KeyHome,
	glfw.KeyEnd:          window.KeyEnd,
	glfw.KeyCapsLock:     window.KeyCapsLock,
	glfw.KeyScrollLock:   window.KeyScrollLock,
	glfw.KeyNumLock:      window.KeyNumLock,
	glfw.KeyPrintScreen:  window.KeyPrintScreen,
	glfw.KeyPause:        window.KeyPause,
	glfw.KeyF1:           window.KeyF1,
	glfw.KeyF2:           window.KeyF2,
	glfw.KeyF3:           window.KeyF3,
	glfw.KeyF4:           window.KeyF4,
	glfw.KeyF5:           window.KeyF5,
	glfw.KeyF6:           window.KeyF6,
	glfw.KeyF7:           window.KeyF7,
	glfw.KeyF8:           window.KeyF8,
	glfw.KeyF9:           window.KeyF9,
	glfw.KeyF10:          window.KeyF10,
	glfw.KeyF11:          window.KeyF11,
	glfw.KeyF12:          window.KeyF12,
	glfw.KeyF13:          window.KeyF13,
	glfw.KeyF14:          window.KeyF14,
	glfw.KeyF15:          window.KeyF15,
	glfw.KeyF16:          window.KeyF16,
	glfw.KeyF17:          window.KeyF17,
	glfw.KeyF18:          window.KeyF18,
	glfw.KeyF19:          window.KeyF19,
	glfw.KeyF20:          window.KeyF20,
	glfw.KeyF21:          window.KeyF21,
	glfw.KeyF22:          window.KeyF22,
	glfw.KeyF23:          window.KeyF23,
	glfw.KeyF24:          window.KeyF24,
	glfw.KeyF25:          window.KeyF25,
	glfw.KeyKP0:          window.KeyKP0,
	glfw.KeyKP1:          window.KeyKP1,
	glfw.KeyKP2:          window.KeyKP2,
	glfw.KeyKP3:          window.KeyKP3,
	glfw.KeyKP4:          window.KeyKP4,
	glfw.KeyKP5:          window.KeyKP5,
	glfw.KeyKP6:          window.KeyKP6,
	glfw.KeyKP7:          window.KeyKP7,
	glfw.KeyKP8:          window.KeyKP8,
	glfw.KeyKP9:          window.KeyKP9,
	glfw.KeyKPDecimal:    window.KeyKPDecimal,
	glfw.KeyKPDivide:     window.KeyKPDivide,
	glfw.KeyKPMultiply:   window.KeyKPMultiply,
	glfw.KeyKPSubtract:   window.KeyKPSubtract,
	glfw.KeyKPAdd:        window.KeyKPAdd,
	glfw.KeyKPEnter:      window.KeyKPEnter,
	glfw.KeyKPEqual:      window.KeyKPEqual,
	glfw.KeyLeftShift:    window.KeyLeftShift,
	glfw.KeyLeftControl:  window.KeyLeftControl,
	glfw.KeyLeftAlt:      window.KeyLeftAlt,
	glfw.KeyLeftSuper:    window.KeyLeftSuper,
	glfw.KeyRightShift:   window.KeyRightShift,
	glfw.KeyRightControl: window.KeyRightControl,
	glfw.KeyRightAlt:     window.KeyRightAlt,
	glfw.KeyRightSuper:   window.KeyRightSuper,
	glfw.KeyMenu:         window.KeyMenu,
}

var actionMapping = map[glfw.Action]window.Action{
	glfw.Release: window.Release,
	glfw.Press:   window.Press,
	glfw.Repeat:  window.Repeat,
}

var modifierKeyMapping = map[glfw.ModifierKey]window.ModifierKey{
	glfw.ModShift:    window.ModShift,
	glfw.ModControl:  window.ModControl,
	glfw.ModAlt:      window.ModAlt,
	glfw.ModSuper:    window.ModSuper,
	glfw.ModCapsLock: window.ModCapsLock,
	glfw.ModNumLock:  window.ModNumLock,
}