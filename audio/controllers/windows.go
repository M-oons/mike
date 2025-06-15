package controllers

import (
	"errors"
	"fmt"

	"github.com/go-ole/go-ole"
	"github.com/m-oons/mike/config"
	"github.com/moutend/go-wca/pkg/wca"
)

type windowsController struct {
	config config.ConfigControllerWindows
}

func NewWindowsController(config config.ConfigControllerWindows) *windowsController {
	return &windowsController{
		config: config,
	}
}

func (c *windowsController) Setup() error {
	return ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
}

func (*windowsController) Mute() error {
	device, err := getCurrentDevice()
	if err != nil {
		return fmt.Errorf("error getting audio device: %w", err)
	}
	if device == nil {
		return errors.New("error muting audio device - audio device is nil")
	}
	defer device.Release()

	if err := device.SetMute(true, nil); err != nil {
		return fmt.Errorf("error muting audio device: %w", err)
	}

	return nil
}

func (*windowsController) Unmute() error {
	device, err := getCurrentDevice()
	if err != nil {
		return fmt.Errorf("error getting audio device: %w", err)
	}
	if device == nil {
		return errors.New("error unmuting audio device - audio device is nil")
	}
	defer device.Release()

	if err := device.SetMute(false, nil); err != nil {
		return fmt.Errorf("error unmuting audio device: %w", err)
	}

	return nil
}

func (*windowsController) ToggleMute() (bool, error) {
	device, err := getCurrentDevice()
	if err != nil {
		return false, fmt.Errorf("error getting audio device: %w", err)
	}
	if device == nil {
		return false, errors.New("error toggling mute state for audio device - audio device is nil")
	}
	defer device.Release()

	var mute bool
	if err := device.GetMute(&mute); err != nil {
		return false, fmt.Errorf("error getting mute state for audio device: %w", err)
	}

	newState := !mute
	if err := device.SetMute(newState, nil); err != nil {
		return false, fmt.Errorf("error toggling mute state for audio device: %w", err)
	}

	return newState, nil
}

func (*windowsController) IsMuted() (bool, error) {
	device, err := getCurrentDevice()
	if err != nil {
		return false, fmt.Errorf("error getting audio device: %w", err)
	}
	if device == nil {
		return false, errors.New("error getting mute state for audio device - audio device is nil")
	}
	defer device.Release()

	var mute bool
	if err := device.GetMute(&mute); err != nil {
		return false, fmt.Errorf("error getting mute state for audio device: %w", err)
	}

	return mute, nil
}

func (*windowsController) Close() {
	ole.CoUninitialize()
}

func getCurrentDevice() (*wca.IAudioEndpointVolume, error) {
	var mmde *wca.IMMDeviceEnumerator
	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return nil, err
	}
	defer mmde.Release()

	var mmd *wca.IMMDevice
	if err := mmde.GetDefaultAudioEndpoint(wca.ECapture, wca.DEVICE_STATE_ACTIVE, &mmd); err != nil {
		return nil, err
	}
	defer mmd.Release()

	var aev *wca.IAudioEndpointVolume
	if err := mmd.Activate(wca.IID_IAudioEndpointVolume, wca.CLSCTX_ALL, nil, &aev); err != nil {
		return nil, err
	}

	return aev, nil
}
