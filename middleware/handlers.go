package middleware

import (
	"RestfulApi/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	return db
}

func CreateCourse(w http.ResponseWriter, r *http.Request) {

	var course models.Course
	err := json.NewDecoder(r.Body).Decode(&course)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}
	insertID := insertCourse(course)
	fmt.Println(course)
	res := response{
		ID:      insertID,
		Message: "Course created successfully",
	}
	json.NewEncoder(w).Encode(res)
}

func GetCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}
	user, err := getCourse(int64(id))

	if err != nil {
		log.Fatalf("Unable to get course. %v", err)
	}
	json.NewEncoder(w).Encode(user)
}

func GetAllCourses(w http.ResponseWriter, r *http.Request) {
	users, err := getAllCourses()
	if err != nil {
		log.Fatalf("Unable to get all courses. %v", err)
	}
	json.NewEncoder(w).Encode(users)

}

func UpdateCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}
	var course models.Course
	err = json.NewDecoder(r.Body).Decode(&course)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}
	updatedRows := updateCourse(int64(id), course)

	msg := fmt.Sprintf("Course updated successfully. Total rows/record affected %v", updatedRows)
	res := response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

func DeleteCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}
	deletedRows := deleteCourse(int64(id))

	// format the message string
	msg := fmt.Sprintf("Course updated successfully. Total rows/record affected %v", deletedRows)
	res := response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

// ------------------------- handler functions ----------------
// insert one course in the DB
func insertCourse(course models.Course) int64 {

	// create connection
	db := createConnection()
	// close the db connection
	defer db.Close()
	var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	//var crs_id int = 1011
	err := db.QueryRow(`INSERT INTO courses ( cname, cduration, cfee) VALUES ($1, $2, $3 ) RETURNING courseid`, course.CName, course.CDuration, course.CFee).Scan(&id)
	fmt.Println(course)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)
	return id
}

// get one course from the DB by its userid
func getCourse(id int64) (models.Course, error) {
	db := createConnection()
	defer db.Close()
	var course models.Course
	sqlStatement := `SELECT * FROM courses WHERE courseid=$1`
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&course.CID, &course.CName, &course.CDuration, &course.CFee)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return course, nil
	case nil:
		return course, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}
	return course, err
}

func getAllCourses() ([]models.Course, error) {
	// create connection
	db := createConnection()
	defer db.Close()

	var courses []models.Course
	sqlStatement := `SELECT * FROM courses`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var course models.Course
		err = rows.Scan(&course.CID, &course.CName, &course.CDuration, &course.CFee)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the user in the users slice
		courses = append(courses, course)

	}

	// return empty user on error
	return courses, err
}

func updateCourse(id int64, course models.Course) int64 {

	db := createConnection()
	defer db.Close()

	sqlStatement := `UPDATE courses SET cname=$2, cduration=$3, cfee=$4 WHERE courseid=$1`
	res, err := db.Exec(sqlStatement, id, course.CName, course.CDuration, course.CFee)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
func deleteCourse(id int64) int64 {

	db := createConnection()
	defer db.Close()
	sqlStatement := `DELETE FROM courses WHERE courseid=$1`
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
