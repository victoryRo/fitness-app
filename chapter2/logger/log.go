package logger

import (
	"bytes"
	"log"
	"net/http"
	"os"

	glog "github.com/kataras/golog"
)

// Logger nueva instacia de golog pkg
var Logger = glog.New()

type remote struct{}

// Write implementación para enviar al servidor remoto
func (r remote) Write(data []byte) (n int, err error) {
	go func() {
		req, err := http.NewRequest("POST", "http://localhost:8010/log", bytes.NewBuffer(data))

		if err == nil {
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			// envia la solicitud http y retorna la respuesta http
			resp, er := client.Do(req)
			if er != nil {
				log.Fatal(er)
			}
			// retorna la respuesta del body y cierra la solicitud
			defer resp.Body.Close()
		}
	}()
	return len(data), nil
}

// SetLoggingOutput establecemos el registro local o remoto
func SetLoggingOutput(localStdout bool) {
	if localStdout {
		configureLocal()
		return
	}
	configureRemote()
}

// configureLocal para implementación local
func configureLocal() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// escribe el registro a la terminal
	Logger.SetOutput(os.Stdout)
	// establece el nivel del registro
	Logger.SetLevel("debug")
	// establece el destino del registro 'file' en el nivel info
	Logger.SetLevelOutput("info", file)
}

// configureRemote para la configuración del registrador remoto
func configureRemote() {
	r := remote{}
	// establece el nivel de registro en info y lo envia con el formato json
	Logger.SetLevelFormat("info", "json")
	// especifica el nivel de registro y la salida del log
	Logger.SetLevelOutput("info", r)
}
