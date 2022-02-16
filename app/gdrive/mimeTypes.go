package gdrive

import "errors"

type DriveMimeType string

// Source: https://developers.google.com/drive/api/v3/mime-types
const (
	Audio              = "application/vnd.google-apps.audio"
	Document           = "application/vnd.google-apps.document"
	ThirdPartyShortcut = "application/vnd.google-apps.drive-sdk"
	Drawing            = "application/vnd.google-apps.drawing"
	File               = "application/vnd.google-apps.file"
	Folder             = "application/vnd.google-apps.folder"
	Form               = "application/vnd.google-apps.form"
	Fusiontable        = "application/vnd.google-apps.fusiontable"
	Map                = "application/vnd.google-apps.map"
	Photo              = "application/vnd.google-apps.photo"
	Presentation       = "application/vnd.google-apps.presentation"
	Script             = "application/vnd.google-apps.script"
	Shortcut           = "application/vnd.google-apps.shortcut"
	Site               = "application/vnd.google-apps.site"
	Spreadsheet        = "application/vnd.google-apps.spreadsheet"
	Unknown            = "application/vnd.google-apps.unknown"
	Video              = "application/vnd.google-apps.video"
)

func (mimeType DriveMimeType) IsValid() error {
	switch mimeType {
	case Audio, Document, ThirdPartyShortcut, Drawing, File, Folder, Form, Fusiontable, Map, Photo, Presentation, Script, Shortcut, Site, Spreadsheet, Unknown, Video:
		return nil
	}

	return errors.New("invalid mime type")
}
