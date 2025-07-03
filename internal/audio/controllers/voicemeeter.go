package controllers

import (
	"fmt"
	"slices"
	"syscall"
	"time"
	"unsafe"

	"github.com/m-oons/mike/internal/config"
)

type voicemeeterController struct {
	config             config.ConfigControllerVoicemeeter
	remoteDLL          *syscall.DLL
	login              *syscall.Proc
	logout             *syscall.Proc
	getVoicemeeterType *syscall.Proc
	getParameterFloat  *syscall.Proc
	setParameterFloat  *syscall.Proc
	isParametersDirty  *syscall.Proc
	parameter          string
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

	parameter, err := c.getTargetParameter()
	if err != nil {
		return fmt.Errorf("error getting Voicemeeter target parameter: %w", err)
	}
	c.parameter = parameter

	return nil
}

func (c *voicemeeterController) Mute() error {
	c.syncParameters()

	param := fmt.Appendf(nil, "%s.Mute\x00", c.parameter)
	ret, _, _ := c.setParameterFloat.Call(
		uintptr(unsafe.Pointer(&param[0])),
		uintptr(1.0),
	)
	if ret != 0 {
		return fmt.Errorf("error muting Voicemeeter parameter %s, code: %d", c.parameter, ret)
	}

	return nil
}

func (c *voicemeeterController) Unmute() error {
	c.syncParameters()

	param := fmt.Appendf(nil, "%s.Mute\x00", c.parameter)
	ret, _, _ := c.setParameterFloat.Call(
		uintptr(unsafe.Pointer(&param[0])),
		uintptr(0.0),
	)
	if ret != 0 {
		return fmt.Errorf("error unmuting Voicemeeter parameter %s, code: %d", c.parameter, ret)
	}

	return nil
}

func (c *voicemeeterController) ToggleMute() (bool, error) {
	c.syncParameters()

	param := fmt.Appendf(nil, "%s.Mute\x00", c.parameter)
	var value float32
	ret, _, _ := c.getParameterFloat.Call(
		uintptr(unsafe.Pointer(&param[0])),
		uintptr(unsafe.Pointer(&value)),
	)
	if ret != 0 {
		return false, fmt.Errorf("error getting mute state for Voicemeeter parameter %s, code: %d", c.parameter, ret)
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
		return false, fmt.Errorf("error toggling mute state for Voicemeeter parameter %s, code: %d", c.parameter, ret)
	}

	return newValue == 1.0, nil
}

func (c *voicemeeterController) IsMuted() (bool, error) {
	c.syncParameters()

	param := fmt.Appendf(nil, "%s.Mute\x00", c.parameter)
	var value float32
	ret, _, _ := c.getParameterFloat.Call(
		uintptr(unsafe.Pointer(&param[0])),
		uintptr(unsafe.Pointer(&value)),
	)
	if ret != 0 {
		return false, fmt.Errorf("error getting mute state for Voicemeeter parameter %s, code: %d", c.parameter, ret)
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

	c.getVoicemeeterType, err = c.remoteDLL.FindProc("VBVMR_GetVoicemeeterType")
	if err != nil {
		return fmt.Errorf("error finding procedure 'VBVMR_GetVoicemeeterType': %w", err)
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

func (c *voicemeeterController) getVersion() (int32, error) {
	var value int32
	ret, _, _ := c.getVoicemeeterType.Call(uintptr(unsafe.Pointer(&value)))
	if ret != 0 {
		return 0, fmt.Errorf("error getting Voicemeeter version, code: %d", ret)
	}
	if value < 1 || value > 3 {
		return value, fmt.Errorf("invalid Voicemeeter version: %d", value)
	}

	return value, nil
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

func (c *voicemeeterController) getTargetParameter() (string, error) {
	parameters, defaultParameter, err := c.getAvailableParameters()
	if err != nil {
		return "", fmt.Errorf("error getting available Voicemeeter parameters: %w", err)
	}

	if slices.Contains(parameters, c.config.Parameter) {
		return c.config.Parameter, nil
	}

	return defaultParameter, nil
}

func (c *voicemeeterController) getAvailableParameters() ([]string, string, error) {
	version, err := c.getVersion()
	if err != nil {
		return nil, "", err
	}

	var parameters []string
	var defaultParameter string

	switch version {
	case 1: // Standard
		parameters = []string{
			// physical inputs
			"Strip[0]",
			"Strip[1]",

			// virtual outputs
			"Bus[2]",
		}
		defaultParameter = "Bus[2]"

	case 2: // Banana
		parameters = []string{
			// physical inputs
			"Strip[0]",
			"Strip[1]",
			"Strip[2]",

			// virtual outputs
			"Bus[3]",
			"Bus[4]",
		}
		defaultParameter = "Bus[3]"

	case 3: // Potato
		parameters = []string{
			// physical inputs
			"Strip[0]",
			"Strip[1]",
			"Strip[2]",
			"Strip[3]",
			"Strip[4]",

			// virtual outputs
			"Bus[5]",
			"Bus[6]",
			"Bus[7]",
		}
		defaultParameter = "Bus[5]"
	}

	return parameters, defaultParameter, nil
}
