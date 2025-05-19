package controllers

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"github.com/m-oons/mike/config"
)

type VoicemeeterController struct {
	remoteDLL         *syscall.DLL
	login             *syscall.Proc
	logout            *syscall.Proc
	getParameterFloat *syscall.Proc
	setParameterFloat *syscall.Proc
	isParametersDirty *syscall.Proc
	bus               string
}

func (c *VoicemeeterController) Init() error {
	dll, err := syscall.LoadDLL(config.Current.Controller.Voicemeeter.RemoteDLLPath)
	if err != nil {
		return err
	}
	c.remoteDLL = dll

	login, err := c.remoteDLL.FindProc("VBVMR_Login")
	if err != nil {
		return err
	}
	c.login = login

	logout, err := c.remoteDLL.FindProc("VBVMR_Logout")
	if err != nil {
		return err
	}
	c.logout = logout

	getParameterFloat, err := c.remoteDLL.FindProc("VBVMR_GetParameterFloat")
	if err != nil {
		return err
	}
	c.getParameterFloat = getParameterFloat

	setParameterFloat, err := c.remoteDLL.FindProc("VBVMR_SetParameterFloat")
	if err != nil {
		return err
	}
	c.setParameterFloat = setParameterFloat

	isParametersDirty, err := c.remoteDLL.FindProc("VBVMR_IsParametersDirty")
	if err != nil {
		return err
	}
	c.isParametersDirty = isParametersDirty

	ret, _, _ := c.login.Call()
	if ret != 0 {
		return fmt.Errorf("error logging in to Voicemeeter, code: %d", ret)
	}

	switch config.Current.Controller.Voicemeeter.Output {
	case 1:
		c.bus = "Bus[3]"
	case 2:
		c.bus = "Bus[4]"
	default:
		c.bus = "Bus[3]"
	}

	return nil
}

func (c *VoicemeeterController) Mute() error {
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

func (c *VoicemeeterController) Unmute() error {
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

func (c *VoicemeeterController) ToggleMute() (bool, error) {
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

func (c *VoicemeeterController) IsMuted() (bool, error) {
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

func (c *VoicemeeterController) Close() error {
	if c.logout != nil {
		c.logout.Call()
	}

	if c.remoteDLL != nil {
		c.remoteDLL.Release()
	}

	return nil
}

func (c *VoicemeeterController) syncParameters() {
	for i := 0; i < 10; i++ {
		ret, _, _ := c.isParametersDirty.Call()
		if ret == 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}
