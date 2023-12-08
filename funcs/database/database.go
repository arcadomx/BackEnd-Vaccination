package database
/*
Paquete para el manejo de la base de datos y las consultas que se realizan a la misma
Este poaquete se puede actualizar para manejar diferentes tipos de bases de datos

Esta configurada para trabajar con postgresql

*/
import (
	"strconv"
	"database/sql"
	"app/funcs/structs"
	"app/funcs/errores"
	_ "github.com/lib/pq"
	"log"
)


var sql_tipo = "postgres"
var sql_db = "user=postgres password=postgres dbname=postgres host=db sslmode=disable"



func DB() SetDB {
	return &setDB{}
}
type setDB struct {
}

type SetDB interface {
	GetUser(user structs.User) (bool, string)
	PostUser(user structs.User) (bool, string)

	GetDrug() (bool, []structs.Drug)
	GetDrugId(drug int) (bool, structs.Drug)
	PostDrug(drug structs.Drug) (bool, string)

	PutDrug(drug structs.Drug) (bool, string)
	DeleteDrug(drug int) (bool, string)

	GetVaccination() (bool, []structs.Vaccination)
	PostVaccination(vaccination structs.Vaccination) (bool, string)

	PutVaccination(vaccination structs.Vaccination) (bool, string)
	DeleteVaccination(vaccination int) (bool, string)
}


func (h *setDB) GetUser(user structs.User) (bool, string) {
	log.Println("DB login usuario")
	db, err := sql.Open(sql_tipo, sql_db)
	errores.CheckErr(err)
	err = db.Ping()
	errores.CheckErr(err)
	var id int
  var email string
	err = db.QueryRow("SELECT id,email FROM users WHERE email = $1", user.Email).Scan(&id, &email)
	if err != nil {
		if err == sql.ErrNoRows {
				return false, "Usuario no encontrado"
		} else {
errores.CheckErr(err)
		}
}

return true, strconv.Itoa(id)+"-"+email
	
}

func (h *setDB) PostUser(user structs.User) (bool, string) {
	log.Println("DB registro de usuario")
	db, err := sql.Open(sql_tipo, sql_db)
	errores.CheckErr(err)
	err = db.Ping()
	errores.CheckErr(err)
	var emailExists string
	err = db.QueryRow("SELECT email FROM users WHERE email = $1", user.Email).Scan(&emailExists)

	if err != nil && err != sql.ErrNoRows {
	errores.CheckErr(err)
	}
	if emailExists != "" {
			return false, "El correo ya esta registrado"
	}
	result, err := db.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", user.Name, user.Email, user.Password)
	errores.CheckErr(err)
	rowsAffected, err := result.RowsAffected()
	errores.CheckErr(err)

	if rowsAffected > 0 {
			return true, "Usuario registrado correctamente"
	} else {
			return false, "Error al registrar usuario"
	}
}

func (h *setDB) GetDrug() (bool, []structs.Drug) {
	log.Println("DB consulta de drugs")
	db, err := sql.Open(sql_tipo, sql_db)
	errores.CheckErr(err)
	err = db.Ping()
	errores.CheckErr(err)
	rows, err := db.Query("SELECT id,name,approved,min_dose,max_dose,available_at FROM drug")
	errores.CheckErr(err)
	defer rows.Close()

	var drugs []structs.Drug
	for rows.Next() {
			var drug structs.Drug
			err = rows.Scan(&drug.Ide, &drug.Name, &drug.Approved, &drug.Min_dose, &drug.Max_dose, &drug.AvailableAt)
			errores.CheckErr(err)
			drugs = append(drugs, drug)
	}
	err = rows.Err()
	errores.CheckErr(err)

	return true, drugs
}

func (h *setDB) PostDrug(drug structs.Drug) (bool, string) {
	log.Println("DB registro de drugs")
	db, err := sql.Open(sql_tipo, sql_db)
	errores.CheckErr(err)
	err = db.Ping()
	errores.CheckErr(err)

	var nameExists string
	err = db.QueryRow("SELECT name FROM drug WHERE name = $1", drug.Name).Scan(&nameExists)

	if err != nil && err != sql.ErrNoRows {
errores.CheckErr(err)
	}

	if nameExists != "" {
			return false, "DB-Error: El drug ya existe"
	}


	result, err := db.Exec("INSERT INTO drug (name, approved, min_dose, max_dose, available_at) VALUES ($1, $2, $3, $4, $5)", drug.Name, drug.Approved, drug.Min_dose, drug.Max_dose, drug.AvailableAt)
	errores.CheckErr(err)
	rowsAffected, err := result.RowsAffected()
	errores.CheckErr(err)

	if rowsAffected > 0 {
			return true, "Drug registrado correctamente"
	} else {
			return false, "Error al registrar drug"
	}
}

func (h *setDB) PutDrug(drug structs.Drug) (bool, string) {
	log.Println("DB actualizacion de drugs")
	db, err := sql.Open(sql_tipo, sql_db)
	errores.CheckErr(err)
	err = db.Ping()
	errores.CheckErr(err)

	result, err := db.Exec("UPDATE drug SET name = $1, approved = $2, min_dose = $3, max_dose = $4, available_at = $5 WHERE id = $6", drug.Name, drug.Approved, drug.Min_dose, drug.Max_dose, drug.AvailableAt, drug.Ide)
	errores.CheckErr(err)
	rowsAffected, err := result.RowsAffected()
	errores.CheckErr(err)

	if rowsAffected > 0 {
			return true, "Drug actualizado correctamente"
	} else {
			return false, "Error al actualizar drug"
	}
}

func (h *setDB) DeleteDrug(drug int) (bool, string) {
	log.Println("DB eliminacion de drugs")
	db, err := sql.Open(sql_tipo, sql_db)
	errores.CheckErr(err)
	err = db.Ping()
	errores.CheckErr(err)

	result, err := db.Exec("DELETE FROM drug WHERE id = $1", drug)
	errores.CheckErr(err)
	rowsAffected, err := result.RowsAffected()
	errores.CheckErr(err)

	if rowsAffected > 0 {
			return true, "Drug eliminado correctamente"
	} else {
			return false, "Error al eliminar drug"
	}
}

func (h *setDB) GetVaccination() (bool, []structs.Vaccination) {
	log.Println("DB consulta de vaccination")
	db, err := sql.Open(sql_tipo, sql_db)
	errores.CheckErr(err)
	err = db.Ping()
	errores.CheckErr(err)
	rows, err := db.Query("SELECT id,name,drug_id,dose,date FROM vaccination")
	errores.CheckErr(err)
	defer rows.Close()

	var vaccinations []structs.Vaccination
	for rows.Next() {
			var vaccination structs.Vaccination
			err = rows.Scan(&vaccination.Ide, &vaccination.Name, &vaccination.Drug_id, &vaccination.Dose, &vaccination.Date)
			errores.CheckErr(err)
			vaccinations = append(vaccinations, vaccination)
	}
	err = rows.Err()
	errores.CheckErr(err)

	return true, vaccinations
}

func (h *setDB) GetDrugId(drug int) (bool, structs.Drug) {
	log.Println("DB consulta de drug por id")
	db, err := sql.Open(sql_tipo, sql_db)
	errores.CheckErr(err)
	err = db.Ping()
	errores.CheckErr(err)
	var drugs structs.Drug
	err = db.QueryRow("SELECT id,name,approved,min_dose,max_dose,available_at FROM drug WHERE id = $1", drug).Scan(&drugs.Ide, &drugs.Name, &drugs.Approved, &drugs.Min_dose, &drugs.Max_dose, &drugs.AvailableAt)
	errores.CheckErr(err)
	if err == sql.ErrNoRows {
			return false, drugs
	} else if err != nil {
    errores.CheckErr(err)
}
return true, drugs
}

func (h *setDB) PostVaccination(vaccination structs.Vaccination) (bool, string) {
	log.Println("DB registro de vaccination")
	db, err := sql.Open(sql_tipo, sql_db)
	errores.CheckErr(err)
	err = db.Ping()
	errores.CheckErr(err)
	result, err := db.Exec("INSERT INTO vaccination (name, drug_id, dose, date) VALUES ($1, $2, $3, $4)", vaccination.Name, vaccination.Drug_id, vaccination.Dose, vaccination.Date)
	errores.CheckErr(err)

	rowsAffected, err := result.RowsAffected()
	errores.CheckErr(err)

	if rowsAffected > 0 {
			return true, "Vacunación insertada con éxito"
	} else {
			return false, "Error al insertar la vacunación"
	}
}
	
func (h *setDB) PutVaccination(vaccination structs.Vaccination) (bool, string) {
	log.Println("DB actualizacion de vaccination")
	db, err := sql.Open(sql_tipo, sql_db)
	errores.CheckErr(err)
	err = db.Ping()
	errores.CheckErr(err)

	result, err := db.Exec("UPDATE vaccination SET name = $1, drug_id = $2, dose = $3, date = $4 WHERE id = $5", vaccination.Name, vaccination.Drug_id, vaccination.Dose, vaccination.Date, vaccination.Ide)
	errores.CheckErr(err)
	rowsAffected, err := result.RowsAffected()
	errores.CheckErr(err)

	if rowsAffected > 0 {
			return true, "Vacunación actualizada correctamente"
	} else {
			return false, "Error al actualizar la vacunación"
	}
}

func (h *setDB) DeleteVaccination(vaccination int) (bool, string) {
	log.Println("DB eliminacion de vaccination")
	db, err := sql.Open(sql_tipo, sql_db)
	errores.CheckErr(err)
	err = db.Ping()
	errores.CheckErr(err)

	result, err := db.Exec("DELETE FROM vaccination WHERE id = $1", vaccination)
	errores.CheckErr(err)
	rowsAffected, err := result.RowsAffected()
	errores.CheckErr(err)

	if rowsAffected > 0 {
			return true, "Vacunación eliminada correctamente"
	} else {
			return false, "Error al eliminar la vacunación"
	}
}