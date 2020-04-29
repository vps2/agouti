package appium

import "github.com/vps2/agouti"

func NewTestDevice(session mobileSession) *Device {
	return &Device{
		Page:    &agouti.Page{},
		session: session,
	}
}
