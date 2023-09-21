package main

import "github.com/kataras/golog"

func main() {
	// golog.SetLevel("info")
	golog.SetLevel("debug")
	// golog.SetLevel("warn")
	// golog.SetLevel("error")
	// golog.SetLevel("fatal")

	golog.Println("This is a raw message, no lavels, no colors")
	golog.Info("This is a info message, with colors (if the output is terminal)")
	golog.Warn("This is a warning message")
	golog.Error("This is a error message")
	golog.Debug("This is a debug message")
	golog.Fatal(`Fatal will exist no matter what`)
}
