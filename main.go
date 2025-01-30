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

func main() {

	router := gin.Default()                 //Initialize a Gin router using Default.
	router.GET("/albums", getAlbums)        //Use the GET function to associate the GET HTTP method and /albums path with a handler function.
	router.POST("/albums", postAlbums)      //Associate the POST method at the /albums path with the postAlbums function.
	router.GET("/albums/:id", getAlbumByID) //Associate the /albums/:id path with the getAlbumByID function. In Gin, the colon preceding an item in the path signifies that the item is a path parameter.
	router.Run("127.0.0.1:8080")            //Use the Run function to attach the router to an http.Server and start the server.

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

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums) //serialize the struct into JSON and add it to the response.
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil { //Use Context.BindJSON to bind the request body to newAlbum.
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)            //Append the album struct initialized from the JSON to the albums slice.
	c.IndentedJSON(http.StatusCreated, newAlbum) //Add a 201 status code to the response, along with JSON representing the album you added.
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id") //Use Context.Param to retrieve the id path parameter from the URL. When you map this handler to a path, you’ll include a placeholder for the parameter in the path.

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a) //If it’s found, you serialize that album struct to JSON and return it as a response with a 200 OK HTTP code.
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"}) //Return an HTTP 404 error with http.StatusNotFound if the album isn’t found.
}
