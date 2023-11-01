package main

import (
	"embed"
	"fmt"
	"strings"
)

var (
	Version string = strings.TrimSpace(version)
	//go:embed version/version.txt
	version string

	//go:embed src/static
	staticEmbed embed.FS

	//go:embed src/css/*
	cssEmbed embed.FS

	//go:embed src/tmpl/*.html
	tmplEmbed embed.FS

	// to continuo i must to run db and queries with sqlc
)

func main() {
	fmt.Println("hello")
}
