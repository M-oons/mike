package main

import (
	"fmt"

	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
)

func main() {
	if err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		return
	}

	var mmde *wca.IMMDeviceEnumerator
	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return
	}

	microphones := getMicrophones()
	for _, mic := range microphones {
		name := getDeviceName(mic)
		fmt.Printf("microphone: %v\n", name)
	}

	// var mmd *wca.IMMDevice
	// if err := mmde.GetDefaultAudioEndpoint(wca.ERender, wca.EConsole, &mmd); err != nil {
	// 	return
	// }

	// var aev *wca.IAudioEndpointVolume
	// if err := mmd.Activate(wca.IID_IAudioEndpointVolume, wca.CLSCTX_ALL, nil, &aev); err != nil {
	// 	return
	// }
}

func getMicrophones() []*wca.IMMDevice {
	var devices []*wca.IMMDevice

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
		devices = append(devices, item)
	}

	return devices
}

func getDeviceName(device *wca.IMMDevice) string {
	if device == nil {
		return ""
	}

	var ps *wca.IPropertyStore
	device.OpenPropertyStore(wca.STGM_READ, &ps)

	var pv wca.PROPVARIANT
	ps.GetValue(&wca.PKEY_Device_FriendlyName, &pv)

	return pv.String()
}
