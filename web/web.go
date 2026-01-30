package web

import "embed"

//go:embed all:dist
var Public embed.FS
