package main

import (
	"log"
	"os"

	"github.com/kataras/golog"
)

const logfile = "infolog.txt"

// funciion init establece el nivel de registro
func init() {
	golog.SetLevel("debug")
	configureLogger()
}

func main() {
	golog.Println("This is a raw message, no lavels, no colors")
	golog.Info("This is a info message, with colors (if the output is terminal)")
	golog.Warn("This is a warning message")
	golog.Error("This is a error message")
	golog.Debug("This is a debug message")
	golog.Fatal(`Fatal will exist no matter what`)
}

// organiza el registro en diferentes archivos
// de acuerdo a su nivel
func configureLogger() {
	infof, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	golog.SetLevelOutput("info", infof)

	errf, err := os.OpenFile("infoerr.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	golog.SetLevelOutput("error", errf)
}

// El resto de mensajes a nivel de registro.
// se escriben en stdout, 'terminal'
// que está configurado de forma predeterminada por la biblioteca.
