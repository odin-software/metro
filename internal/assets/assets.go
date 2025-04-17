package assets

import "embed"

//go:embed files/*
var Assets embed.FS
