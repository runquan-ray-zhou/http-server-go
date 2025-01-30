package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums) //serialize the struct into JSON and add it to the response.
}

func main() {

	router := gin.Default()          //Initialize a Gin router using Default.
	router.GET("/albums", getAlbums) //Use the GET function to associate the GET HTTP method and /albums path with a handler function.

	router.Run("127.0.0.1:8080") //Use the Run function to attach the router to an http.Server and start the server.

	// l, err := net.Listen("tcp", "127.0.0.1:4221")
	// if err != nil {
	// 	fmt.Println("Failed to bind to port 4221")
	// 	os.Exit(1)
	// }

	// conn, err := l.Accept()
	// if err != nil {
	// 	fmt.Println("Error accepting connection: ", err.Error())
	// 	os.Exit(1)
	// }

	// fmt.Println(conn)

}
