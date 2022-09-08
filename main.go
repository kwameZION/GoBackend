package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type todo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

var todos = []todo{}

//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//CSV Datei Reading
func ReadCSV() {
	var newTodos []todo

	//opening Data
	file, err := os.Open("data.csv")
	if err != nil {
		fmt.Println(err)
	}
	//Read Data
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	//changing Records in struct and append array
	//start at the second entry (i=1) to skip the header
	for i := 1; i < len(records); i++ {
		done := false
		if records[i][2] == "true" {
			done = true
		}

		item := todo{Id: records[i][0], Name: records[i][1], Done: done}
		newTodos = append(newTodos, item)
	}
	// Replacing old List with New List
	todos = newTodos

}

//*******************************************************************************************************************************************************
//Writing Todos in CSV
func writeCSV() {
	//Open Data / Create CSV
	csvFile, err := os.Create("data.csv")
	if err != nil {
		fmt.Println(err)
	}

	csvwriter := csv.NewWriter(csvFile)
	// Header
	header := []string{"id", "name", "done"}
	csvwriter.Write(header)

	for _, record := range todos {
		item := []string{record.Id, record.Name, fmt.Sprintf("%t", record.Done)}

		_ = csvwriter.Write(item)
	}

	// Saving the file
	csvwriter.Flush()
	csvFile.Close()
}

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
//CRUD FUNKTION

//Get todo
func getTodos(context *gin.Context) {

	// Storing the Data
	ReadCSV()

	//output
	context.IndentedJSON(http.StatusOK, todos)

}

//***************************************************************************************************************************************
// Post , adding new todos
func postTodo(context *gin.Context) {
	var newTodo todo

	if err := context.BindJSON(&newTodo); err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
		return
	}
	// Loop for the ID
	lastId := 0
	for i := 0; i < len(todos); i++ {
		currentId, _ := strconv.Atoi(todos[i].Id)
		if currentId > lastId {
			lastId = currentId
		}
	}
	newTodo.Id = strconv.Itoa(lastId + 1)
	todos = append(todos, newTodo)

	// Storing the Data
	writeCSV()

	//output
	context.IndentedJSON(http.StatusCreated, newTodo)

}

//**************************************************************************************************************************************************
// getting a specific todo
func getTodo(context *gin.Context) {
	id := context.Param("id")
	todos, err := getTodoById((id))

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})
		return
	}

	// Storing the Data
	ReadCSV()

	//output
	context.IndentedJSON(http.StatusOK, todos)

}

//*********************************************************************************************************************************************
// Put (update/put)
func updateTodo(context *gin.Context) {
	id := context.Param("id")

	var updateData todo

	//Create payload as todo struct
	if err := context.BindJSON(&updateData); err != nil {
		return
	}

	for _, item := range todos {
		if item.Id == id {
			item = updateData

			context.IndentedJSON(http.StatusOK, todos)
			return
		}
	}

	// Determine der highest ID
	max := 0
	for _, item := range todos {
		itemId, _ := strconv.Atoi(item.Id)
		if itemId > max {
			max = itemId
		}
	}

	// Storing the Data
	writeCSV()

	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})

}

//********************************************************************************************************************************************
// Error, if the todo is not found
func getTodoById(id string) (*todo, error) {

	for _, todo := range todos {
		if todo.Id == id {
			return &todo, nil
		}

	}
	return nil, errors.New("todo not found")
}

//*******************************************************************************************************************************************************
// Function defining the remove(delete function)
func RemoveIndex(s []todo, index int) []todo {
	return append(s[:index], s[index+1:]...)
}

// Delete
func deleteTodo(context *gin.Context) {
	id := context.Param("id")

	for i, todo := range todos {
		if todo.Id == id {
			todos = RemoveIndex(todos, i)
			return
		}
	}
	// storing the data
	writeCSV()

	//output
	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Todo not found"})

}

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Main function
func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/todos", getTodos)          // Get todolists
	router.GET("/todos/:id", getTodo)       // get a specific todo
	router.PUT("/todos", updateTodo)        // change to bool of the todo ( its like put/update)
	router.POST("/todos", postTodo)         // adding a todo to the list
	router.DELETE("/todos/:id", deleteTodo) // Deleting a todo

	router.Run("localhost:5000") // Port 5000

	writeCSV() // write CSv
	ReadCSV()  // Read CSV

}
