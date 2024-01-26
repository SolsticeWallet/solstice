package solstice

// This file is generated !
// DO NOT EDIT

import "fyne.io/fyne/v2"

func GetMetadata() fyne.AppMetadata {
	md := App.Metadata()
	if md.ID == "" {
		md = fyne.AppMetadata {
			ID: "",
			Name: "",
			Version: "",
			Build: 0,
			Icon: nil,
			Release: false,
			Custom: map[string]string{},
		}
	}
	return md
}
