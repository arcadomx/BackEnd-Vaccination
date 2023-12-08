package handler
/*
Este paquete maneja toda la informacion que se recibe desde el cliente y se envia al cliente desde formato json, ademas de verificar que la informacion recibida sea confiable y segura
*/

import (
	"app/funcs/database"
	"app/funcs/errores"
	"app/funcs/structs"
	"app/funcs/jsonWebToken"
	"log"
	"net/http"
	"encoding/json"
	"strings"
	"strconv"
	"time"
)

var db = database.DB()

func NewApi() SetApi {
	return &setApi{}
}

type setApi struct {
}

type SetApi interface {
	Signup(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)

	Drug(w http.ResponseWriter, r *http.Request)
	DrugAct(w http.ResponseWriter, r *http.Request)

	Vaccination(w http.ResponseWriter, r *http.Request)
	VaccinationAct(w http.ResponseWriter, r *http.Request)

}

func encabezado(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type","json")
}

func (h *setApi) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
}
log.Println("HF registro de usuario")
var user structs.User
err := json.NewDecoder(r.Body).Decode(&user)
if err != nil {
		errores.CheckErr(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
}

if user.Ide != 0 ||  user.Name == "" || user.Email == "" || user.Password == "" {
http.Error(w, "Datos con errores", http.StatusBadRequest)
} else {
	
check, text := db.PostUser(user)
encabezado(w)
if check == false {
	w.WriteHeader(http.StatusBadRequest)
} else {
	w.WriteHeader(http.StatusOK)
}
response := map[string]string{"mensaje": text}
            err := json.NewEncoder(w).Encode(response)
            errores.CheckErr(err)
}
}

func (h *setApi) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
}
var user structs.User
err := json.NewDecoder(r.Body).Decode(&user)
if err != nil {
	errores.CheckErr(err)
	http.Error(w, err.Error(), http.StatusBadRequest)
	return
}
log.Println("HF login >> ",user.Email)

if user.Ide != 0 ||  user.Name != "" || user.Email == "" || user.Password == "" {
	http.Error(w, "Datos con errores", http.StatusBadRequest)
	} else {
		check, text := db.GetUser(user)
		encabezado(w)
		if check == false {
				w.WriteHeader(http.StatusBadRequest)
		} else {
				w.WriteHeader(http.StatusOK)
				token,err := jsonWebToken.GenerateToken(text)
				errores.CheckErr(err)
				response := map[string]string{"token": token}
				err = json.NewEncoder(w).Encode(response)
				errores.CheckErr(err)
		}
	}
}

func (h *setApi) Drug(w http.ResponseWriter, r *http.Request) {

	TrueToken, err := jsonWebToken.ValidateToken(r.Header.Get("Authorization"))
	errores.CheckErr(err)

	if TrueToken == false {
		http.Error(w, "Token invalido", http.StatusBadRequest)
		return
	}

	encabezado(w)
	switch r.Method {
		
	case http.MethodGet:
		log.Println("HF consulta de medicamentos")		
		check, drugs := db.GetDrug()
		if check == false {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(drugs)
			errores.CheckErr(err)
		}

	case http.MethodPost:
		log.Println("HF registro de medicamento")
		var drugs structs.Drug
		err := json.NewDecoder(r.Body).Decode(&drugs)
		errores.CheckErr(err)

		if drugs.Ide != 0 || drugs.Name == "" ||  drugs.Min_dose == 0 || drugs.Max_dose == 0 || drugs.AvailableAt == "" {
			http.Error(w, "Datos con errores", http.StatusBadRequest)
		} else {
			encabezado(w)
			check, text := db.PostDrug(drugs)
			if check == false {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusOK)
			}
			response := map[string]string{"mensaje": text}
			err := json.NewEncoder(w).Encode(response)
			errores.CheckErr(err)
		}
	
	default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (h *setApi) DrugAct(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) != 3 {
			http.Error(w, "Ruta inválida", http.StatusBadRequest)
			return
	}
	id := parts[2]
	drugID, err := strconv.Atoi(id)
	if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			errores.CheckErr(err)
			return
	}
	
	TrueToken, err := jsonWebToken.ValidateToken(r.Header.Get("Authorization"))
	errores.CheckErr(err)

	if TrueToken == false {
		http.Error(w, "Token invalido", http.StatusBadRequest)
		return
	}

encabezado(w)
var check bool
var text string
	switch r.Method {
		case http.MethodPut:
		
		log.Println("HF actualizacion de drug")
		var drugs structs.Drug
		err = json.NewDecoder(r.Body).Decode(&drugs)
		errores.CheckErr(err)
		drugs.Ide = drugID

		if drugs.Ide == 0 || drugs.Name == "" ||  drugs.Min_dose == 0 || drugs.Max_dose == 0 || drugs.AvailableAt == "" {
			http.Error(w, "Datos con errores", http.StatusBadRequest)
		} else {
			
			check, text = db.PutDrug(drugs)
			if check == false {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusOK)
			}
			response := map[string]string{"mensaje": text}
			err := json.NewEncoder(w).Encode(response)
			errores.CheckErr(err)
		}

		case http.MethodDelete:
		log.Println("HF eliminacion de drug")
		check, text = db.DeleteDrug(drugID)
		if check == false {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}

			
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		
}
response := map[string]string{"mensaje": text}
err = json.NewEncoder(w).Encode(response)
errores.CheckErr(err)

}

func (h *setApi) Vaccination(w http.ResponseWriter, r *http.Request) {
	
	TrueToken, err := jsonWebToken.ValidateToken(r.Header.Get("Authorization"))
	errores.CheckErr(err)

	if TrueToken == false {
		http.Error(w, "Token invalido", http.StatusBadRequest)
		return
	}

	encabezado(w)
	switch r.Method {
		
	case http.MethodGet:
		log.Println("HF consulta de vacunacion")		
		check, vaccination := db.GetVaccination()
		if check == false {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(vaccination)
			errores.CheckErr(err)
		}

	case http.MethodPost:
		log.Println("HF registro de vacunacion")
		var vaccination structs.Vaccination
		err := json.NewDecoder(r.Body).Decode(&vaccination)
		errores.CheckErr(err)

		if vaccination.Ide != 0 || vaccination.Name == "" ||  vaccination.Drug_id == 0 || vaccination.Dose == 0 || vaccination.Date == "" {
			http.Error(w, "Datos con errores", http.StatusBadRequest)
		} else {
			encabezado(w)

			checkDrug, textDrug := db.GetDrugId(vaccination.Drug_id)
			if checkDrug == false {
				w.WriteHeader(http.StatusBadRequest)
				response := map[string]string{"mensaje": "NO cumple los requisitos de vacunacion"}
				err := json.NewEncoder(w).Encode(response)
				errores.CheckErr(err)
				return
			}

			vaccinationDate, err := time.Parse("2006-01-02T15:04:05Z", vaccination.Date)
			if err != nil {
					http.Error(w, "Error en fecha de Vacunacion", http.StatusBadRequest)
					return
			}
			
			drugAvailableAt, err := time.Parse("2006-01-02T15:04:05Z", textDrug.AvailableAt)
			if err != nil {
					http.Error(w, "Error en fecha de drug", http.StatusBadRequest)
					return
			}

			if vaccinationDate.Before(drugAvailableAt) {
				http.Error(w, "La vacunacion esta antes de la disponibilidad", http.StatusBadRequest)
				return
		}
		
		if vaccination.Dose < textDrug.Min_dose || vaccination.Dose > textDrug.Max_dose {
				http.Error(w, "Vaccination dose esta fuera de rango de dosis", http.StatusBadRequest)
				return
		}


			check, text := db.PostVaccination(vaccination)
			if check == false {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusOK)
			}
			response := map[string]string{"mensaje": text}
			err = json.NewEncoder(w).Encode(response)
			errores.CheckErr(err)
		}
	
	default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (h *setApi) VaccinationAct(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) != 3 {
			http.Error(w, "Ruta inválida", http.StatusBadRequest)
			return
	}
	id := parts[2]
	vaccinationID, err := strconv.Atoi(id)
	if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			errores.CheckErr(err)
			return
	}
	
	TrueToken, err := jsonWebToken.ValidateToken(r.Header.Get("Authorization"))
	errores.CheckErr(err)

	if TrueToken == false {
		http.Error(w, "Token invalido", http.StatusBadRequest)
		return
	}

encabezado(w)
var check bool
var text string
var Vaccination structs.Vaccination
		
	switch r.Method {
		case http.MethodPut:
			log.Println("HF actualizacion de Vaccination")
			err = json.NewDecoder(r.Body).Decode(&Vaccination)
			errores.CheckErr(err)
			Vaccination.Ide = vaccinationID

		if Vaccination.Ide	== 0 || Vaccination.Name == "" ||  Vaccination.Drug_id == 0 || Vaccination.Dose == 0 || Vaccination.Date == "" {
			http.Error(w, "Datos con errores", http.StatusBadRequest)
		} else {
			

			check, text = db.PutVaccination(Vaccination)
			if check == false {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusOK)
			}
			response := map[string]string{"mensaje": text}
			err := json.NewEncoder(w).Encode(response)
			errores.CheckErr(err)
		}

		case http.MethodDelete:
		log.Println("HF eliminacion de Vaccination")
		check, text = db.DeleteDrug(Vaccination.Ide)
		if check == false {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}

			
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
response := map[string]string{"mensaje": text}
err = json.NewEncoder(w).Encode(response)
errores.CheckErr(err)

}
		 