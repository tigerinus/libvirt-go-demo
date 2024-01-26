package schema

import "encoding/xml"

type StoragePool struct {
	Source StoragePoolSource `xml:"source"`
	Target StoragePoolTarget `xml:"target"`

	XMLName xml.Name        `xml:"pool"`
	Name    string          `xml:"name"`
	Type    StoragePoolType `xml:"type,attr"`
}

type StoragePoolSource struct {
	Directory StoragePoolSourceDirectory `xml:"dir"`
}

type StoragePoolSourceDirectory struct {
	Path string `xml:"path,attr"`
}

type StoragePoolTarget struct {
	Path        string                       `xml:"path"`
	Permissions StoragePoolTargetPermissions `xml:"permissions"`
}

type StoragePoolTargetPermissions struct {
	Owner string `xml:"owner"`
	Group string `xml:"group"`
	Mode  string `xml:"mode"`
}

type StoragePoolType string

const StoragePoolTypeDir StoragePoolType = "dir"
