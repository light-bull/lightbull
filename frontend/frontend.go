package frontend

import (
	"embed"
	"io/fs"
)

//go:embed files/*
var tmp embed.FS
var Frontend, _ = fs.Sub(tmp, "files")
