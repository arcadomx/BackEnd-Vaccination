package structs

/*
En este paquete estan definidas las estructuras iniciales y que se ocupan dentro de toda la aplicacion
Datos iniciales, son los que se pueden cambiar si mover nada dentro de la aplicacion con el acrhivo config.json

User, Drug, Vaccination son las estructuras que se ocupan para el manejo de datos dentro de la aplicacion
*/

import (
	"app/funcs/errores"
	"github.com/spf13/viper"
	"strconv"
)

type DatosIniciales struct {
	Port string
	NombrePrincipal string
	Anio string
	Extra string
	JsonTokenTime int
}

var DatosInit DatosIniciales

func CargarConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	errores.CheckErr(err)

port := viper.Get("port").(string)
name := viper.Get("nombre").(string)
anio := viper.Get("anio").(string)
extra := viper.Get("extra").(string)
jsonTokenTime := viper.Get("jsonTokenTime").(string)
jwt,err := strconv.Atoi(jsonTokenTime)
errores.CheckErr(err)

	DatosInit = DatosIniciales{
		Port: port,
		NombrePrincipal: name,
		Anio: anio,
		Extra: extra,
		JsonTokenTime: jwt,
	}
}

type User struct {
	Ide int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type Drug struct {
	Ide int `json:"id"`
	Name string `json:"name"`
	Approved bool `json:"approved"`
	Min_dose int `json:"min_dose"`
	Max_dose int `json:"max_dose"`
	AvailableAt string `json:"availableAt"`
}

type Vaccination struct {
	Ide int `json:"id"`
	Name string `json:"name"`
	Drug_id int `json:"drug_id"`
	Dose int `json:"dose"`
	Date string `json:"date"`
}