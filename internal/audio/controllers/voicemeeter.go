package controllers

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"github.com/m-oons/mike/internal/config"
)

type voicemeeterController struct {
	config            config.ConfigControllerVoicemeeter
	remoteDLL         *syscall.DLL
	login             *syscall.Proc
	logout            *syscall.Proc
	getParameterFloat *syscall.Proc
	setParameterFloat *syscall.Proc
	isParametersDirty *syscall.Proc
	bus               string
}

func NewVoicemeeterController(config config.ConfigControllerVoicemeeter) *voicemeeterController {
	return &voicemeeterController{
		config: config,
	}
}

func (c *voicemeeterController) Setup() error {
	if err := c.loadAPI(); err != nil {
		return fmt.Errorf("error loading Voicemeeter Remote API: %w", err)
	}

	// ensure Voicemeeter is running and login is successful
	for {
		ret, _, _ := c.login.Call()
		if ret == 0 {
			break
		}
		c.logout.Call()
		time.Sleep(2 * time.Second)
	}

	switch c.config.Output {
	case 1:
		c.bus = "Bus[3]"
	case 2:
		c.bus = "Bus[4]"
	default:
		c.bus = "Bus[3]"
	}

	return nil
}

func (c *voicemeeterController) Mute() error {
	c.syncParameters()

	param := fmt.Appendf(nil, "%s.Mute\x00", c.bus)
	ret, _, _ := c.setParameterFloat.Call(
		uintptr(unsafe.Pointer(&param[0])),
		uintptr(1.0),
	)
	if ret != 0 {
		return fmt.Errorf("error muting Voicemeeter bus %s, code: %d", c.bus, ret)
	}

	return nil
}

func (c *voicemeeterController) Unmute() error {
	c.syncParameters()

	param := fmt.Appendf(nil, "%s.Mute\x00", c.bus)
	ret, _, _ := c.setParameterFloat.Call(
		uintptr(unsafe.Pointer(&param[0])),
		uintptr(0.0),
	)
	if ret != 0 {
		return fmt.Errorf("error unmuting Voicemeeter bus %s, code: %d", c.bus, ret)
	}

	return nil
}

func (c *voicemeeterController) ToggleMute() (bool, error) {
	c.syncParameters()

	param := fmt.Appendf(nil, "%s.Mute\x00", c.bus)
	var value float32
	ret, _, _ := c.getParameterFloat.Call(
		uintptr(unsafe.Pointer(&param[0])),
		uintptr(unsafe.Pointer(&value)),
	)
	if ret != 0 {
		return false, fmt.Errorf("error getting mute state for Voicemeeter bus %s, code: %d", c.bus, ret)
	}

	var newValue float32
	if value == 1.0 {
		newValue = 0.0
	} else {
		newValue = 1.0
	}

	ret, _, _ = c.setParameterFloat.Call(
		uintptr(unsafe.Pointer(&param[0])),
		uintptr(newValue),
	)
	if ret != 0 {
		return false, fmt.Errorf("error toggling mute state for Voicemeeter bus %s, code: %d", c.bus, ret)
	}

	return newValue == 1.0, nil
}

func (c *voicemeeterController) IsMuted() (bool, error) {
	c.syncParameters()

	param := fmt.Appendf(nil, "%s.Mute\x00", c.bus)
	var value float32
	ret, _, _ := c.getParameterFloat.Call(
		uintptr(unsafe.Pointer(&param[0])),
		uintptr(unsafe.Pointer(&value)),
	)
	if ret != 0 {
		return false, fmt.Errorf("error getting mute state for Voicemeeter bus %s, code: %d", c.bus, ret)
	}

	return value == 1.0, nil
}

func (c *voicemeeterController) Close() {
	if c.logout != nil {
		c.logout.Call()
	}

	if c.remoteDLL != nil {
		c.remoteDLL.Release()
	}
}

func (c *voicemeeterController) loadAPI() error {
	var err error
	c.remoteDLL, err = syscall.LoadDLL(c.config.RemoteDLLPath)
	if err != nil {
		return fmt.Errorf("error loading Remote API DLL: %w", err)
	}

	c.login, err = c.remoteDLL.FindProc("VBVMR_Login")
	if err != nil {
		return fmt.Errorf("error finding procedure 'VBVMR_Login': %w", err)
	}

	c.logout, err = c.remoteDLL.FindProc("VBVMR_Logout")
	if err != nil {
		return fmt.Errorf("error finding procedure 'VBVMR_Logout': %w", err)
	}

	c.getParameterFloat, err = c.remoteDLL.FindProc("VBVMR_GetParameterFloat")
	if err != nil {
		return fmt.Errorf("error finding procedure 'VBVMR_GetParameterFloat': %w", err)
	}

	c.setParameterFloat, err = c.remoteDLL.FindProc("VBVMR_SetParameterFloat")
	if err != nil {
		return fmt.Errorf("error finding procedure 'VBVMR_SetParameterFloat': %w", err)
	}

	c.isParametersDirty, err = c.remoteDLL.FindProc("VBVMR_IsParametersDirty")
	if err != nil {
		return fmt.Errorf("error finding procedure 'VBVMR_IsParametersDirty': %w", err)
	}

	return nil
}

func (c *voicemeeterController) syncParameters() {
	for range 10 {
		ret, _, _ := c.isParametersDirty.Call()
		if ret == 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}
