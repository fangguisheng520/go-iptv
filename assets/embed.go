package assets

import "embed"

//go:embed static/favicon.ico
var Favicon []byte

//go:embed templates/*
var EmbeddedFS embed.FS

//go:embed static/*
var StaticFS embed.FS
