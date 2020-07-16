package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
	"strings"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "username"
	password = "password"
	dbname   = "covid19-contact"
)

var db *sql.DB

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/register", func(c *gin.Context) {
		query := "INSERT INTO users DEFAULT VALUES RETURNING id;"
		row := db.QueryRow(query)
		var id string
		row.Scan(&id)
		c.JSON(http.StatusOK, gin.H{
			"success" : true,
			"id" : "COVID19-" + id,
		})
		return
	})

	r.POST("/positive/:id", func(c *gin.Context) {
		id := c.Param("id")
		if strings.Contains(id,"COVID19-") {
			id = strings.Split(id, "-")[1]
		}
		query := "INSERT INTO positive_users (id, diagnosed) VALUES ($1, $2);"
		_, err := db.Exec(query, id, c.PostForm("time"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success" : false,
				"error" : "Database error",
			})
			panic(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success" : true,
		})
		return
	})
	
	r.POST("/submitSession/:id", func(c *gin.Context) {
		id := c.Param("id")
		if strings.Contains(id,"COVID19-") {
			id = strings.Split(id, "-")[1]
		}
		ids := strings.Split(c.PostForm("ids"), ",")
		times := strings.Split(c.PostForm("times"), ",")

		fmt.Println("FICK" + c.PostForm("ids"))

		query1 := "INSERT INTO contact (person1, person2, contact_time) VALUES "

		query2 := ""
		for i := 0; i < len(ids); i++ {
			query2 += "($1," + ids[i] + "," + times[i] + ")"
		}
		query2 += ";"

		fmt.Println(query1 + query2)

		_, err := db.Exec(query1 + query2, id)

		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success" : false,
				"error" : "Database error",
			})
			panic(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success" : true,
		})
	})

	r.GET("/alerts/:id", func(c *gin.Context) {
		id := c.Param("id")
		if strings.Contains(id,"COVID19-") {
			id = strings.Split(id, "-")[1]
		}
		query := "SELECT * FROM positive_users WHERE id=(SELECT person1 FROM contact WHERE person2=$1) OR id=(SELECT person2 FROM contact WHERE person1=$1);"
		rows, err := db.Query(query, id)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success" : false,
				"error" : "Database error",
			})
			panic(err)
			return
		}

		ids := ""
		times := ""

		defer rows.Close()
		for rows.Next() {
			var id string
			var time string
			err = rows.Scan(&id, &time)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"success" : false,
					"error" : "Database error",
				})
				panic(err)
				return
			}
			ids += "," + id
			times += "," + time
		}

		if len(ids) > 0 {
			ids = ids[1:]
			times = times[1:]
		}

		c.JSON(http.StatusOK, gin.H{
			"success" : true,
			"ids" : ids,
			"times" : times,
		})
		return
	})

	r.RunTLS(":8080", "debug-key/domain-crt.txt", "debug-key/domain-key.txt")
}
