package utils

import (
	"errors"
	"reflect"
	"strings"
	"syscall"
	"unsafe"
)

var (
	user32 = syscall.MustLoadDLL("user32.dll")
	modadvapi32 = syscall.MustLoadDLL("advapi32.dll")

	procEnumWindows              = user32.MustFindProc("EnumWindows")
	procEnumChildWindows         = user32.MustFindProc("EnumChildWindows")
	procIsWindow                 = user32.MustFindProc("IsWindow")
	procGetWindow                = user32.MustFindProc("GetWindow")
	procFindWindow               = user32.MustFindProc("FindWindowW")
	procGetWindowTextW           = user32.MustFindProc("GetWindowTextW")
	procSetWindowsHookEx         = user32.MustFindProc("SetWindowsHookExW")
	procSetWinEventHook          = user32.MustFindProc("SetWinEventHook")
	procUnhookWinEvent           = user32.MustFindProc("UnhookWinEvent")
	procGetWindowThreadProcessID = user32.MustFindProc("GetWindowThreadProcessId")
	procSetFocus                 = user32.MustFindProc("SetFocus")
	procSetForegroundWindow      = user32.MustFindProc("SetForegroundWindow")
	procSetWindowPos             = user32.MustFindProc("SetWindowPos")
	procEnableWindow             = user32.MustFindProc("EnableWindow")
	procShowWindow               = user32.MustFindProc("ShowWindow")
	procMoveWindow               = user32.MustFindProc("MoveWindow")
	procGetWindowRect            = user32.MustFindProc("GetWindowRect")
	procScreenToClient           = user32.MustFindProc("ScreenToClient")
	procGetSystemMetrics         = user32.MustFindProc("GetSystemMetrics")

	procRegQueryValueEx = modadvapi32.MustFindProc("RegQueryValueExW")
)

// ...
const (
	WINEVENT_OUTOFCONTEXT = 0x0000
	WM_SYSCOMMAND         = 0x0112
	WM_SETREDRAW          = 11
	SC_MOVE               = 61456
	HTCAPTION             = 2

	EVENT_MIN                   = 0x00000001
	EVENT_MAX                   = 0x7FFFFFFF
	EVENT_SYSTEM_FOREGROUND     = 0x0003
	EVENT_SYSTEM_MOVESIZESTART  = 0x000A
	EVENT_SYSTEM_MOVESIZEEND    = 0x000B
	EVENT_SYSTEM_SCROLLINGEND   = 0x0013
	EVENT_SYSTEM_MINIMIZESTART  = 0x0016
	EVENT_SYSTEM_MINIMIZEEND    = 0x0017
	EVENT_OBJECT_SHOW           = 0x8002
	EVENT_OBJECT_HIDE           = 0x8003
	EVENT_OBJECT_LOCATIONCHANGE = 0x800B
	EVENT_OBJECT_DRAGSTART      = 0x8021
	EVENT_OBJECT_DRAGCANCEL     = 0x8022
	EVENT_OBJECT_DRAGCOMPLETE   = 0x8023
	EVENT_OBJECT_DRAGENTER      = 0x8024
	EVENT_OBJECT_DRAGLEAVE      = 0x8025
	EVENT_OBJECT_DRAGDROPPED    = 0x8026
	EVENT_OBJECT_IME_CHANGE     = 0x8029

	GW_HWNDFIRST    = 0
	GW_HWNDLAST     = 1
	GW_HWNDNEXT     = 2
	GW_HWNDPREV     = 3
	GW_OWNER        = 4
	GW_CHILD        = 5
	GW_ENABLEDPOPUP = 6

	SW_HIDE            = 0
	SW_NORMAL          = 1
	SW_SHOWNORMAL      = 1
	SW_SHOWMINIMIZED   = 2
	SW_MAXIMIZE        = 3
	SW_SHOWMAXIMIZED   = 3
	SW_SHOWNOACTIVATE  = 4
	SW_SHOW            = 5
	SW_MINIMIZE        = 6
	SW_SHOWMINNOACTIVE = 7
	SW_SHOWNA          = 8
	SW_RESTORE         = 9
	SW_SHOWDEFAULT     = 10
	SW_FORCEMINIMIZE   = 11

	SM_CXFULLSCREEN = 16
	SM_CYFULLSCREEN = 17

	HWND_BOTTOM    = 1
	HWND_NOTOPMOST = -2
	HWND_TOP       = 0
	HWND_TOPMOST   = -1

	SWP_ASYNCWINDOWPOS = 0x4000
	SWP_DEFERERASE     = 0x2000
	SWP_DRAWFRAME      = 0x0020
	SWP_SHOWWINDOW     = 0x0040
	SWP_FRAMECHANGED   = 0x0020
	SWP_HIDEWINDOW     = 0x0080
	SWP_NOACTIVATE     = 0x0010
	SWP_NOCOPYBITS     = 0x0100
	SWP_NOSIZE         = 0x0001
	SWP_NOMOVE         = 0x0002
	SWP_NOZORDER       = 0x0004
)

const (
	NO_ERROR      = 0
	ERROR_SUCCESS = 0
)

// BoolToBOOL ...
func BoolToBOOL(value bool) int {
	if value {
		return 1
	}

	return 0
}

func GetSystemMetrics(index int) int {
	ret, _, _ := procGetSystemMetrics.Call(
		uintptr(index))

	return int(ret)
}

// EnumChildWindows ...
func EnumChildWindows(hwnd syscall.Handle, enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := procEnumChildWindows.Call(uintptr(hwnd), uintptr(enumFunc), uintptr(lparam))
	// r1, _, e1 := syscall.Syscall(procEnumChildWindows.Addr(), 3, uintptr(hwnd), uintptr(enumFunc), uintptr(lparam))
	if r1 == 0 {
		if e1 != nil && strings.Contains(e1.Error(), "Success") {
			return e1
		}
		// if e1 != 0 {
		// 	err = error(e1)
		// } else {
		// 	err = syscall.EINVAL
		// }
	}
	return
}

// EnumWindows ...
func EnumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, uintptr(enumFunc), uintptr(lparam), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

// GetWindowText ...
func GetWindowText(hwnd syscall.Handle) (str string, err error) {
	b := make([]uint16, 200)
	maxCount := int32(len(b))
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(&b[0])), uintptr(maxCount))
	len := int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
		return
	}
	str = syscall.UTF16ToString(b)
	return
}

// GetThread ...
func GetThread(hwnd syscall.Handle) (syscall.Handle, uint) {
	var id uint
	ret, _, _ := procGetWindowThreadProcessID.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&id)))

	return syscall.Handle(ret), id
}

// SetForegroundWindow ...
func SetForegroundWindow(hwnd syscall.Handle) bool {
	ret, _, _ := procSetForegroundWindow.Call(uintptr(hwnd))
	return ret != 0
}

// SetFocus ...
func SetFocus(hwnd syscall.Handle) syscall.Handle {
	ret, _, _ := procSetFocus.Call(uintptr(hwnd))
	return syscall.Handle(ret)
}

// EnableWindow ...
func EnableWindow(hwnd syscall.Handle, b bool) bool {
	ret, _, _ := procEnableWindow.Call(uintptr(hwnd), uintptr(BoolToBOOL(b)))
	return ret != 0
}

// ShowWindow ...
func ShowWindow(hwnd syscall.Handle, flag int) bool {
	ret, _, _ := procShowWindow.Call(uintptr(hwnd), uintptr(flag))
	return ret != 0
}

// MoveWindow ...
func MoveWindow(hwnd syscall.Handle, x, y, width, height int, repaint bool) bool {
	ret, _, _ := procMoveWindow.Call(
		uintptr(hwnd),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(BoolToBOOL(repaint)))

	return ret != 0
}

// SetWindowPos ...
func SetWindowPos(hwnd, hwndInsertAfter syscall.Handle, x, y, width, height int, flag uint) bool {
	ret, _, _ := procSetWindowPos.Call(
		uintptr(hwnd),
		uintptr(hwndInsertAfter),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(flag),
	)

	return ret != 0
}

// SetWindowZ ...
func SetWindowZ(hwnd syscall.Handle, order int) bool {
	return SetWindowPos(hwnd, syscall.Handle(order), 0, 0, 0, 0, SWP_NOMOVE|SWP_NOSIZE)
}

// RECT ...
type RECT struct {
	Left, Top, Right, Bottom int32
}

// GetWindowRect ...
func GetWindowRect(hwnd syscall.Handle) *RECT {
	var rect RECT
	procGetWindowRect.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&rect)))

	return &rect
}

// GetWindow ...
func GetWindow(hwnd syscall.Handle, rl int) (oh syscall.Handle, ok bool) {
	r1, _, e0 := syscall.Syscall(procGetWindow.Addr(), 2, uintptr(hwnd), uintptr(rl), 0)
	if r1 == 0 {
		if e0 != 0 {
			return
		}
	}
	oh = syscall.Handle(r1)
	ok = true
	return
}

// IsWindow ...
func IsWindow(hwnd syscall.Handle) bool {
	is, _, _ := procIsWindow.Call(uintptr(hwnd))
	if is == 0 {
		return false
	}
	return true
}

// Window ...
type Window struct {
	hwnd      syscall.Handle
	title     string
	wevthooks map[interface{}]syscall.Handle
	thread    *Thread
}

// Title ...
func (w *Window) Title() string {
	return w.title
}

// Handle ...
func (w *Window) Handle() syscall.Handle {
	return w.hwnd
}

// Thread ...
type Thread struct {
	hwnd syscall.Handle
	id   uint
}

// NewThread ...
func NewThread(h syscall.Handle, i uint) *Thread {
	return &Thread{
		h, i,
	}
}

// NewWindow ...
func NewWindow(hwnd syscall.Handle) *Window {
	title, _ := GetWindowText(hwnd)
	return &Window{
		hwnd:      hwnd,
		title:     title,
		thread:    NewThread(GetThread(hwnd)), // get the thread
		wevthooks: make(map[interface{}]syscall.Handle),
	}
}

// ParentWindow ...
func (w *Window) ParentWindow() *Window {
	oh, ok := GetWindow(w.hwnd, GW_OWNER)
	if !ok {
		return nil
	}
	return NewWindow(oh)
}

// ChildWindow ...
func (w *Window) ChildWindow() *Window {
	oh, ok := GetWindow(w.hwnd, GW_CHILD)
	if !ok {
		return nil
	}
	return NewWindow(oh)
}

// Move ...
func (w *Window) Move(x, y, width, height int, repaint bool) bool {
	return MoveWindow(w.hwnd, x, y, width, height, repaint)
}

// Show ...
func (w *Window) Show() bool {
	return ShowWindow(w.hwnd, SW_SHOWDEFAULT)
}

// Hidden ...
func (w *Window) Hidden() bool {
	return ShowWindow(w.hwnd, SW_HIDE)
}

// Enable ...
func (w *Window) Enable() bool {
	return EnableWindow(w.hwnd, true)
}

// Disable ...
func (w *Window) Disable() bool {
	return EnableWindow(w.hwnd, false)
}

// SetForeground ...
func (w *Window) SetForeground() bool {
	return SetForegroundWindow(w.hwnd)
}

// SetZ ...
func (w *Window) SetZ(order int) bool {
	return SetWindowZ(w.hwnd, order)
}

// GetRect ...
func (w *Window) GetRect() *RECT {
	return GetWindowRect(w.hwnd)
}

// GetPosition ...
func (w *Window) GetPosition() (x, y int) {
	rect := w.GetRect()
	return int(rect.Left), int(rect.Top)
}

// GetSize ...
func (w *Window) GetSize() (x, y int) {
	rect := w.GetRect()
	x = int(rect.Right - rect.Left)
	y = int(rect.Bottom - rect.Top)
	return
}

// Resize()

// SetPosition ...
func (w *Window) SetPosition(x, y int) {
	wd, hd := w.GetSize()
	if wd == 0 {
		wd = 100
	}
	if hd == 0 {
		hd = 25
	}
	w.Move(x, y, wd, hd, true)
}

// SetToCenter ...
func (w *Window) SetToCenter() {
	sWidth := GetSystemMetrics(SM_CXFULLSCREEN)
	sHeight := GetSystemMetrics(SM_CYFULLSCREEN)
	if sWidth != 0 && sHeight != 0 {
		ww, hh := w.GetSize()
		w.SetPosition((sWidth/2)-(ww/2), (sHeight/2)-(hh/2))
	}
}

// ListChildWindows ...
func (w *Window) ListChildWindows() ([]*Window, error) {
	ws := []*Window{}
	cb := syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
		ws = append(ws, NewWindow(hwnd))
		return 1 // continue enumeration
	})
	return ws, EnumChildWindows(w.hwnd, cb, 0)
}

// UnhookWinEvent ...
func (w *Window) UnhookWinEvent(fn interface{}) bool {
	fnkey := reflect.ValueOf(fn)
	v, ok := w.wevthooks[fnkey]
	if !ok {
		return ok
	}
	ret, _, _ := procUnhookWinEvent.Call(uintptr(v))
	return ret != 0
}

// Event ...
type Event struct {
	Hook     syscall.Handle
	HWND     syscall.Handle
	Type     int
	ObjectID int32
	ChildID  int32
	CreateAt uint32
}

// SetWinEventHook ... TODO: use options
func (w *Window) SetWinEventHook(fn func(evt *Event) error, evts ...int) (syscall.Handle, error) {

	fnkey := reflect.ValueOf(fn)
	if v, ok := w.wevthooks[fnkey]; ok {
		return v, nil
	}

	// create a new fn
	ofn := func(hook syscall.Handle, evt uint32, hwnd syscall.Handle, idObject int32, idChild int32, dwEventThread uint32, dwmsEventTime uint32) syscall.Handle {
		_ = fn(&Event{
			Hook:     hook,
			HWND:     hwnd,
			Type:     int(evt),
			ObjectID: idObject,
			ChildID:  idChild,
			CreateAt: dwEventThread,
		})

		return 1
	}

	evtMin := EVENT_MIN
	evtMax := EVENT_MAX

	switch len(evts) {
	case 0:
	case 1:
		evtMin = evts[0]
		evtMax = evts[0]
	default:
		evtMin = evts[0]
		evtMax = evts[1]
	}

	ret, _, e0 := syscall.Syscall9(
		procSetWinEventHook.Addr(), 7,
		uintptr(evtMin), uintptr(evtMax),
		uintptr(0), // ??? dll without
		syscall.NewCallback(ofn),
		uintptr(w.thread.id), // ??? process id - w.thread.id
		uintptr(0),           // ??? thread id
		uintptr(WINEVENT_OUTOFCONTEXT),
		0, 0,
	)

	if ret == 0 {
		if e0 != 0 {
			return 0, error(e0)
		}
		return 0, syscall.EINVAL
	}

	// store
	w.wevthooks[fnkey] = syscall.Handle(ret)

	return syscall.Handle(ret), nil
}

// FindWindow ...
func FindWindow(title string) (*Window, error) {

	ret, _, _ := syscall.Syscall(procFindWindow.Addr(), 2,
		uintptr(unsafe.Pointer(nil)),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		0)

	if ret == 0 {
		return nil, errors.New("can't found")
	}

	return NewWindow(syscall.Handle(ret)), nil
}

// ListWindows ...
func ListWindows() ([]*Window, error) {
	ws := []*Window{}
	cb := syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
		is, _, _ := procIsWindow.Call(uintptr(hwnd))
		if is == 0 {
			return 1
		}
		ws = append(ws, NewWindow(hwnd))
		return 1 // continue enumeration
	})
	return ws, EnumWindows(cb, 0)
}

func RegQueryValueEx(hKey syscall.Handle, subKey string) string {
	var bufLen uint32
	ret, _, _ := procRegQueryValueEx.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		uintptr(0),
		uintptr(0),
		uintptr(0),
		uintptr(unsafe.Pointer(&bufLen)),
	)

	if bufLen == 0 {
		return ""
	}

	buf := make([]uint16, bufLen)
	ret, _, _ = procRegQueryValueEx.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		uintptr(0),
		uintptr(0),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bufLen)),
	)

	if ret != ERROR_SUCCESS {
		return ""
	}

	return syscall.UTF16ToString(buf)
}