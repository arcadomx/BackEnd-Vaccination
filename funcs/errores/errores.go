package errores

/*
Todos los errores que se generen dentro de la aplicacion se manejan desde este paquete para poder escribirlos dentro del archivo errores.log
para identificar de manera rapida el error y la linea donde se genero.
*/

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

func CheckErr(err error) {
	if err != nil {
		file, Ferror := os.OpenFile("errores.log", os.O_APPEND|os.O_WRONLY, 0644)
		if Ferror != nil {
			panic(Ferror)
		}
		_, parte, line, _ := runtime.Caller(1)
		text := time.Now().Format("2006-01-02 15:04:05") + parte + ": " + strconv.Itoa(line) + " | Error : " + err.Error() + "\n"
		fmt.Fprintf(file, "%v", text)
		file.Close()
		panic(err)

	}
}