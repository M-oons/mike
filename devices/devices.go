package devices

import (
	"sort"

	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
)

func InitializeCOM() {
	ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
}

func GetMicrophones() []*Device {
	var devices []*Device

	var mmde *wca.IMMDeviceEnumerator
	wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde)

	var mmdc *wca.IMMDeviceCollection
	mmde.EnumAudioEndpoints(wca.ECapture, wca.DEVICE_STATE_ACTIVE, &mmdc)

	var count uint32
	mmdc.GetCount(&count)

	var i uint32
	for i = 0; i < count; i++ {
		var mmd *wca.IMMDevice
		mmdc.Item(i, &mmd)

		var aev *wca.IAudioEndpointVolume
		mmd.Activate(wca.IID_IAudioEndpointVolume, wca.CLSCTX_ALL, nil, &aev)

		device := Device{
			MMDevice: mmd,
			Volume:   aev,
		}
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
	mmde.GetDefaultAudioEndpoint(wca.ECapture, wca.DEVICE_STATE_ACTIVE, &mmd)

	var aev *wca.IAudioEndpointVolume
	mmd.Activate(wca.IID_IAudioEndpointVolume, wca.CLSCTX_ALL, nil, &aev)

	device := Device{
		MMDevice: mmd,
		Volume:   aev,
	}

	return &device
}
