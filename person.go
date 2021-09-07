package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type person struct {
	id         int
	first_name string
	last_name  string
	email      string
	ip_address string
}

func deletePerson(db *sql.DB, idToDelete string) int64 {

	stmt, err := db.Prepare("DELETE FROM people where id = ?")
	checkErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(idToDelete)
	checkErr(err)

	affected, err := res.RowsAffected()
	checkErr(err)

	return affected

}

func updatePerson(db *sql.DB, ourPerson person) int64 {

	stmt, err := db.Prepare("UPDATE people set first_name = ?, last_name = ?, email = ?, ip_address = ? where id = ?")
	checkErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(ourPerson.first_name, ourPerson.last_name, ourPerson.email, ourPerson.ip_address, ourPerson.id)
	checkErr(err)

	affected, err := res.RowsAffected()
	checkErr(err)

	return affected
}

func getPersonById(db *sql.DB, ourID string) person {

	fmt.Printf("Value is %v", ourID)

	rows, _ := db.Query("SELECT id, first_name, last_name, email, ip_address FROM people WHERE id = '" + ourID + "'")
	defer rows.Close()

	ourPerson := person{}

	for rows.Next() {
		rows.Scan(&ourPerson.id, &ourPerson.first_name, &ourPerson.last_name, &ourPerson.email, &ourPerson.ip_address)
	}

	fmt.Printf("Person: %v\n", ourPerson)
	return ourPerson
}

func searchForPerson(db *sql.DB, searchString string) []person {

	rows, err := db.Query("SELECT id, first_name, last_name, email, ip_address FROM people WHERE first_name like '%" + searchString + "%' OR last_name like '%" + searchString + "%'")

	defer rows.Close()

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	people := make([]person, 0)

	for rows.Next() {
		ourPerson := person{}
		err = rows.Scan(&ourPerson.id, &ourPerson.first_name, &ourPerson.last_name, &ourPerson.email, &ourPerson.ip_address)
		if err != nil {
			log.Fatal(err)
		}

		people = append(people, ourPerson)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return people
}

func addPerson(db *sql.DB, newPerson person) {

	stmt, _ := db.Prepare("INSERT INTO people (id, first_name, last_name, email, ip_address) VALUES (?, ?, ?, ?, ?)")
	stmt.Exec(nil, newPerson.first_name, newPerson.last_name, newPerson.email, newPerson.ip_address)
	defer stmt.Close()

	fmt.Printf("Added %v %v \n", newPerson.first_name, newPerson.last_name)
}

func showRecent(db *sql.DB, numresults int) (personlist []person) {

	rows, err := db.Query("SELECT id, first_name, last_name, email, ip_address FROM people ORDER by id DESC LIMIT (" + strconv.Itoa(numresults) + ")")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	people := make([]person, 0)

	for rows.Next() {
		ourPerson := person{}
		err = rows.Scan(&ourPerson.id, &ourPerson.first_name, &ourPerson.last_name, &ourPerson.email, &ourPerson.ip_address)
		if err != nil {
			log.Fatal(err)
		}

		people = append(people, ourPerson)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return people

}

func ShowPerson() {
	fmt.Println("Show Person was called!")
}
