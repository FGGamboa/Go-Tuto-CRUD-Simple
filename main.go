package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var templates = template.Must(template.ParseGlob("templates/*"))

func connectionDb() (connection *sql.DB) {
	driver := "mysql"
	user := "root"
	password := "123456"
	database := "go_tutorial1"

	connection, err := sql.Open(driver, user+":"+password+"@tcp(127.0.0.1)/"+database)
	if err != nil {
		panic(err.Error())
	}

	return connection
}

type Employee struct {
	Id    int
	Name  string
	Email string
}

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/delete", Delete)

	log.Println("Iniciando servidor...")
	http.ListenAndServe(":8080", nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	connection := connectionDb()
	rows, err := connection.Query("SELECT * FROM employee")
	if err != nil {
		panic(err.Error())
	}

	employee := Employee{}
	employees := []Employee{}

	for rows.Next() {
		var id int
		var name, email string

		err = rows.Scan(&id, &name, &email)

		if err != nil {
			panic(err.Error())
		}

		employee.Id = id
		employee.Name = name
		employee.Email = email

		employees = append(employees, employee)

	}

	templates.ExecuteTemplate(w, "index", employees)
}

func Create(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(w, "create", nil)
	case "POST":
		name := r.FormValue("name")
		email := r.FormValue("email")

		insertEmployee(name, email)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			panic(err.Error())
		}
		employee := getEmployee(id)
		templates.ExecuteTemplate(w, "edit", employee)
	case "POST":
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			panic(err.Error())
		}
		name := r.FormValue("name")
		email := r.FormValue("email")

		updateEmployee(id, name, email)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		panic(err.Error())
	}
	deleteEmployee(id)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func insertEmployee(name, email string) {
	connection := connectionDb()
	insertEmployee, err := connection.Prepare("INSERT INTO employee(name, email) VALUES(?, ?)")
	if err != nil {
		panic(err.Error())
	}
	insertEmployee.Exec(name, email)
}

func getEmployee(id int) Employee {
	connection := connectionDb()
	row, err := connection.Query("SELECT * FROM employee WHERE id = ?", id)
	var employee Employee

	if err != nil {
		panic(err.Error())
	}
	for row.Next() {
		var id int
		var name, email string
		if err := row.Scan(&id, &name, &email); err != nil {
			panic(err.Error())
		}
		employee = Employee{Id: id, Name: name, Email: email}
	}
	return employee
}

func updateEmployee(id int, name, email string) {
	connection := connectionDb()
	updateEmployee, err := connection.Prepare("UPDATE employee SET name = ?, email = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	updateEmployee.Exec(name, email, id)
}

func deleteEmployee(id int) {
	connection := connectionDb()
	deleteEmployee, err := connection.Prepare("Delete FROM employee WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	deleteEmployee.Exec(id)
}
