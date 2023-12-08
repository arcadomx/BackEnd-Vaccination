package main

/*
	App API REST para coordinar registros de vacunacion
	Aqui los principal es cargar la configuracion inicial y la parte del servidor con los puertos y rutas
*/


import (
	"log"
 	"net/http"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"

	"app/funcs"
	"app/funcs/structs"
	"app/funcs/errores"
)	




func main() {
	log.Println("Inicia levantamiento de Servidor")

	structs.CargarConfig()
	
	log.Println("Carga de Configuracion Completa")


	router := mux.NewRouter()
	app := handler.NewApi()

	
	router.HandleFunc("/signup", app.Signup)
	router.HandleFunc("/login", app.Login)

	router.HandleFunc("/drug", app.Drug)
	router.HandleFunc("/drug/{id}", app.DrugAct)

	router.HandleFunc("/vaccination", app.Vaccination)
	router.HandleFunc("/vaccination/{id}", app.VaccinationAct)

	
	server := http.Server{
		Addr:         structs.DatosInit.Port,
		Handler:      gziphandler.GzipHandler(router),
		ReadTimeout:  30 * time.Second,
		IdleTimeout:  1 * time.Minute,
		WriteTimeout: 30 * time.Second,
	}

	err := server.ListenAndServe()
	//err := server.ListenAndServeTLS("./cert.crt", "./cert.key")
	errores.CheckErr(err)
	
}
