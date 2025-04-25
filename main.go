package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "example"
	dbname   = "postgres"
)

func HomeRouteHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Try /ping",
	})
}

type Data struct {
	From   string `form:"from"`
	Target string `form:"target"`
}

func PingRouteHandler(c *gin.Context) {
	var data Data
	c.Bind(&data)
	fmt.Println(data)

	c.JSON(http.StatusOK, gin.H{
		"message": "ponging " + data.From,
	})
}

func is_numeric(word string) bool {
	return regexp.MustCompile(`\d`).MatchString(word)
}

func CheckPingable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data Data
		c.Bind(&data)

		if is_numeric(data.From) || is_numeric(data.Target) {
			fmt.Println("Invalid")
			c.AbortWithStatusJSON(401, gin.H{"message": "Please recheck the from and target parameters"})
		}

		c.Next()
	}
}

func connectToDB() error {
	connectionString := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%d sslmode=disable", user, password, dbname, port)

	conn, err := sql.Open("postgres", connectionString)

	if err != nil {
		return err
	}

	err = conn.Ping()
	if err != nil {
		return err
	}
	// rows, err := conn.Query("SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = 'public'")
	// for rows.Next() {
	// 	var tableName string
	// 	if err := rows.Scan(&tableName); err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println(tableName)
	// }

	defer conn.Close()
	return nil
}

func main() {

	err := connectToDB()
	if err != nil {
		fmt.Println("Error connecting to database")
	}
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	r.GET("/", HomeRouteHandler)
	r.GET("/ping", CheckPingable(), PingRouteHandler)

	r.Run(":8000")
}
