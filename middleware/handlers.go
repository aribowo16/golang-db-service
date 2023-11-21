package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"go-postgres/models" // models package where User schema is defined
	"log"
	"net/http" // used to access the request and response object of the api
	"os"       // used to read the environment variable
	"strconv"  // package used to covert string into int type

	"github.com/gorilla/mux" // used to get the params from the route

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
	Data []models.User `json:"data"`
}

// create connection with postgres db
func createConnection() *sql.DB {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Panicf("Error loading .env file")
	}

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to DB!")
	// return the connection
	return db
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// CreateUser create a user in the postgres db
func CreateUser(w http.ResponseWriter, r *http.Request) {
	message := "ERROR"

	defer func() {
        if r := recover(); r != nil {
            // format a response object
		res := response{
			Message: message,
		}

		// send the response
		respondWithJSON(w, http.StatusOK, res)
        }
    }()

	// create an empty user of type models.User
	var user models.User

	// decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Panicf("Unable to decode the request body.  %v", err)
	}

	// call insert user function and pass the user
	insertID := insertUser(user)

	user.ID = insertID

	if (insertID > 0) {
		message = "SUCCESS"
	}

	var users []models.User
	users = append(users, user)

	// format a response object
	res := response{
		ID:      insertID,
		Message: message,
		Data: users,
	}

	// send the response
	respondWithJSON(w, http.StatusOK, res)
}

// GetUser will return a single user by its id
func GetUser(w http.ResponseWriter, r *http.Request) {
	message := "ERROR"

	defer func() {
        if r := recover(); r != nil {
            // format a response object
		res := response{
			Message: message,
		}

		// send the response
		respondWithJSON(w, http.StatusOK, res)
        }
    }()

	// get the userid from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Panicf("Unable to convert the string into int.  %v", err)
	}

	// call the getUser function with user id to retrieve a single user
	user, err := getUser(int64(id))

	if err != nil {
		log.Panicf("Unable to get user. %v", err)
	}

	if (user.ID > 0) {
		message = "SUCCESS"
	} else {
		message = "NOT_FOUND"
	}

	var users []models.User
	users = append(users, user)

	// format a response object
	res := response{
		Data:      users,
		Message: message,
	}

	// send all the users as response
	respondWithJSON(w, http.StatusOK, res)
}

// GetAllUser will return all the users
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	message := "ERROR"

	defer func() {
        if r := recover(); r != nil {
            // format a response object
		res := response{
			Message: message,
		}

		// send the response
		respondWithJSON(w, http.StatusOK, res)
        }
    }()

	// get all the users in the db
	users, err := getAllUsers()

	if err != nil {
		log.Panicf("Unable to get all users. %v", err)
	}

	if (len(users) > 0) {
		message = "SUCCESS"
	} else {
		message = "NOT_FOUND"
	}

	// format a response object
	res := response{
		Data:      users,
		Message: message,
	}

	// send all the users as response
	respondWithJSON(w, http.StatusOK, res)
}

//------------------------- handler functions ----------------
// insert one user in the DB
func insertUser(user models.User) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	// returning userid will return the id of the inserted user
	sqlStatement := `INSERT INTO users (name, location, age) VALUES ($1, $2, $3) RETURNING userid`

	// the inserted id will store in this id
	var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := db.QueryRow(sqlStatement, user.Name, user.Location, user.Age).Scan(&id)

	if err != nil {
		log.Panicf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

// get one user from the DB by its userid
func getUser(id int64) (models.User, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create a user of models.User type
	var user models.User

	// create the select sql query
	sqlStatement := `SELECT * FROM users WHERE userid=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to user
	err := row.Scan(&user.ID, &user.Name, &user.Age, &user.Location)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return user, nil
	case nil:
		return user, nil
	default:
		log.Panicf("Unable to get row data. %v", err)
	}

	// return empty user on error
	return user, err
}

// get one user from the DB by its userid
func getAllUsers() ([]models.User, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var users []models.User

	// create the select sql query
	sqlStatement := `SELECT * FROM users`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Panicf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var user models.User

		// unmarshal the row object to user
		err = rows.Scan(&user.ID, &user.Name, &user.Age, &user.Location)

		if err != nil {
			log.Panicf("Unable to get row data. %v", err)
		}

		// append the user in the users slice
		users = append(users, user)

	}

	// return empty user on error
	return users, err
}
