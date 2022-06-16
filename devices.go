package main

import (
	"sort"

	"github.com/moutend/go-wca/pkg/wca"
)

type Device struct {
	MMDevice *wca.IMMDevice
}

func (device *Device) Id() string {
	if device.MMDevice == nil {
		return ""
	}

	var ps *wca.IPropertyStore
	device.MMDevice.OpenPropertyStore(wca.STGM_READ, &ps)

	var pv wca.PROPVARIANT
	ps.GetValue(&wca.PKEY_AudioEndpoint_GUID, &pv)

	return pv.String()
}

func (device *Device) Name() string {
	if device.MMDevice == nil {
		return ""
	}

	var ps *wca.IPropertyStore
	device.MMDevice.OpenPropertyStore(wca.STGM_READ, &ps)

	var pv wca.PROPVARIANT
	ps.GetValue(&wca.PKEY_Device_FriendlyName, &pv)

	return pv.String()
}

func (device *Device) IsMuted() bool {
	if device.MMDevice == nil {
		return false
	}

	var aev *wca.IAudioEndpointVolume
	if err := device.MMDevice.Activate(wca.IID_IAudioEndpointVolume, wca.CLSCTX_ALL, nil, &aev); err != nil {
		return false
	}

	var mute bool
	aev.GetMute(&mute)

	return mute
}

func (device *Device) ToggleMute() bool {
	if device.MMDevice == nil {
		return false
	}

	var aev *wca.IAudioEndpointVolume
	if err := device.MMDevice.Activate(wca.IID_IAudioEndpointVolume, wca.CLSCTX_ALL, nil, &aev); err != nil {
		return false
	}

	currentState := device.IsMuted()
	newState := !currentState
	aev.SetMute(newState, nil)

	return newState
}

func GetMicrophones() []*Device {
	var devices []*Device

	var mmde *wca.IMMDeviceEnumerator
	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return nil
	}

	var mmdc *wca.IMMDeviceCollection
	mmde.EnumAudioEndpoints(wca.ECapture, wca.DEVICE_STATE_ACTIVE, &mmdc)

	var count uint32
	mmdc.GetCount(&count)

	var i uint32
	for i = 0; i < count; i++ {
		var item *wca.IMMDevice
		mmdc.Item(i, &item)
		device := Device{MMDevice: item}
		devices = append(devices, &device)
	}

	// order devices ascending by name
	sort.Slice(devices, func(i, j int) bool {
		return devices[i].Name() < devices[j].Name()
	})

	return devices
}

func GetCurrentMicrophone() *Device {
	var mmde *wca.IMMDeviceEnumerator
	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return nil
	}

	var mmd *wca.IMMDevice
	if err := mmde.GetDefaultAudioEndpoint(wca.ECapture, wca.DEVICE_STATE_ACTIVE, &mmd); err != nil {
		return nil
	}

	device := Device{MMDevice: mmd}

	return &device
}
