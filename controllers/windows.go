package controllers

import (
	"errors"

	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
)

type WindowsController struct{}

func (c *WindowsController) Init() error {
	return ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
}

func (*WindowsController) Mute() error {
	device, err := getCurrentDevice()
	if err != nil {
		return err
	}
	if device == nil {
		return errors.New("device is null")
	}
	defer device.Release()

	return device.SetMute(true, nil)
}

func (*WindowsController) Unmute() error {
	device, err := getCurrentDevice()
	if err != nil {
		return err
	}
	if device == nil {
		return errors.New("device is null")
	}
	defer device.Release()

	return device.SetMute(false, nil)
}

func (*WindowsController) ToggleMute() (bool, error) {
	device, err := getCurrentDevice()
	if err != nil {
		return false, err
	}
	if device == nil {
		return false, errors.New("device is null")
	}
	defer device.Release()

	var mute bool
	if err := device.GetMute(&mute); err != nil {
		return false, err
	}

	newState := !mute
	if err := device.SetMute(newState, nil); err != nil {
		return false, err
	}

	return newState, nil
}

func (*WindowsController) IsMuted() (bool, error) {
	device, err := getCurrentDevice()
	if err != nil {
		return false, err
	}
	if device == nil {
		return false, errors.New("device is null")
	}
	defer device.Release()

	var mute bool
	if err := device.GetMute(&mute); err != nil {
		return false, err
	}

	return mute, nil
}

func (*WindowsController) Close() error {
	ole.CoUninitialize()
	return nil
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
