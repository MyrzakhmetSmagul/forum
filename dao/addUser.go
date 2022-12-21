package dao

import (
	"database/sql"
	"fmt"
	"forum/models"
	"log"
)

func AddUser(db *sql.DB, u *models.User) error {
	fmt.Println("######################################################\n\nMETHOD ADD USER THAT ADD USER IN DB \nFIRST CHECKPOINT\n\n######################################################")
	record := `INSERT INTO 	users(name, surname, gender, email, pwd) VALUES(?, ?, ?, ?, ?)`
	query, err := db.Prepare(record)
	if err != nil {
		fmt.Printf("##################################################\n%s\n##################################################\n", err)
		return err
	}
	fmt.Println("######################################################\n\nMETHOD ADD USER THAT ADD USER IN DB\nSECOND CHECKPOINT\n\n######################################################")

	_, err = query.Exec(u.Name, u.Surname, u.Gender, u.Email, u.Pwd)
	if err != nil {
		fmt.Printf("##################################################\n%s\n##################################################\n", err)
		return err
	}
	fmt.Println("######################################################\n\nMETHOD ADD USER THAT ADD USER IN DB\nTHIRD CHECKPOINT\n\n######################################################")
	log.Println("INSERT INTO OK")
	return nil
}
