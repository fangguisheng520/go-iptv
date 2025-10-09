package dto

import "encoding/xml"

type TV struct {
	XMLName       xml.Name     `xml:"tv"`
	GeneratorName string       `xml:"generator_info_name,attr"`
	GeneratorURL  string       `xml:"generator_info_url,attr"`
	Channels      []XmlChannel `xml:"channel"`
	Programmes    []Programme  `xml:"programme"`
}

type XmlChannel struct {
	ID          string      `xml:"id,attr"`
	DisplayName DisplayName `xml:"display-name"`
}

type Programme struct {
	Start   string `xml:"start,attr"`
	Stop    string `xml:"stop,attr"`
	Channel string `xml:"channel,attr"`
	Title   Title  `xml:"title"`
	Desc    Desc   `xml:"desc"`
}

type DisplayName struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

// Title 节目标题，支持 lang 属性
type Title struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

// Desc 节目描述，支持 lang 属性
type Desc struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}
