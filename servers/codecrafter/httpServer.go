package codecrafter

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"net"
	"os"
	"slices"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "127.0.0.1:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleRequest(conn) // go routine
	}
}

func handleRequest(conn net.Conn) {
	for {
		req := make([]byte, 1024) // creates a slice of 1024 bytes initialized to zero.
		//It is useful for scenarios requiring temporary storage of binary data.
		conn.Read(req)

		stringRequest := string(req)

		// Check if client wants to close connection
		shouldClose := strings.Contains(stringRequest, "Connection: close")

		if strings.HasPrefix(stringRequest, "GET /files") {
			directoryResponse := handleDirectoryEndpoint(stringRequest, shouldClose)
			conn.Write([]byte(directoryResponse))
			if shouldClose {
				conn.Close()
				return
			}
			continue
		}

		if strings.HasPrefix(stringRequest, "GET /user-agent") {
			userAgentResponse := handleUserAgentEndpoint(stringRequest, shouldClose)
			conn.Write([]byte(userAgentResponse))
			if shouldClose {
				conn.Close()
				return
			}
			continue
		}

		if strings.HasPrefix(stringRequest, "GET /echo") {
			echoResponse := handleEchoEndpoint(stringRequest, shouldClose)
			conn.Write([]byte(echoResponse))
			if shouldClose {
				conn.Close()
				return
			}
			continue
		}

		if strings.HasPrefix(stringRequest, "POST /files") {
			postResponse := handlePostRequest(stringRequest, shouldClose)
			conn.Write([]byte(postResponse))
			if shouldClose {
				conn.Close()
				return
			}
			continue
		}
		indexURL := stringRequest[:6]
		if indexURL != "GET / " {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			conn.Close()
			return
		} else {
			if shouldClose {
				conn.Write([]byte("HTTP/1.1 200 OK\r\nConnection: close\r\n\r\n"))
				conn.Close()
				return
			} else {
				conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
				continue
			}
		}
	}

}

func handlePostRequest(request string, shouldClose bool) string {
	directory := os.Args[2]
	directoryString := request[12:]
	fileName := strings.Split(directoryString, " ")[0]
	requestBody := strings.Split(request, "\n")
	postContent := fmt.Sprint(requestBody[len(requestBody)-1])
	trimmedPostContent := strings.ReplaceAll(postContent, "\x00", "")

	err := os.WriteFile(directory+fileName, []byte(trimmedPostContent), 0664)
	if err != nil {
		log.Fatal(err)
	}
	if shouldClose {
		return "HTTP/1.1 201 Created\r\nConnection: close\r\n\r\n"
	}
	return "HTTP/1.1 201 Created\r\n\r\n"
}

func handleDirectoryEndpoint(request string, shouldClose bool) string {
	directoryString := request[11:]
	fileName := strings.Split(directoryString, " ")[0]
	directory := os.Args[2] // grabs the directory in the test

	file, err := os.ReadFile(directory + fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "HTTP/1.1 404 Not Found\r\n\r\n"
	}
	connectionHeader := "Connection: open"
	if shouldClose {
		connectionHeader = "Connection: close"
	}

	return fmt.Sprintf(
		"HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\n%s\r\nContent-Length: %d\r\n\r\n%s", connectionHeader, len(file), string(file))
}

func handleUserAgentEndpoint(request string, shouldClose bool) string {
	responseBody := strings.Split(request, "User-Agent: ")
	correctBody := strings.Split(responseBody[1], " ")
	endPoint := strings.Split(correctBody[0], "\r")
	connectionHeader := "Connection: open"
	if shouldClose {
		connectionHeader = "Connection: close"
	}
	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n%s\r\nContent-Length: %v\r\n\r\n%v", connectionHeader, len(endPoint[0]), endPoint[0])
}

func handleEchoEndpoint(request string, shouldClose bool) string {
	echoedString := request[10:]
	connectionHeader := "Connection: open"
	if shouldClose {
		connectionHeader = "Connection: close"
	}

	if strings.Contains(echoedString, "Accept-Encoding:") {
		var buf bytes.Buffer
		gzipWriter := gzip.NewWriter(&buf)
		requestBody := strings.Split(echoedString, " ")[0]

		requestBodyBytes := []byte(requestBody)
		_, err := gzipWriter.Write(requestBodyBytes)
		if err != nil {
			panic(err)
		}

		err = gzipWriter.Close()
		if err != nil {
			panic(err)
		}

		encodingMethod := strings.Split(echoedString, "Accept-Encoding: ")[1]
		cleanedUpEncodingMethod := strings.Split(encodingMethod, "\r\n\r\n")[0]
		containsGzipMethod := slices.Contains(strings.Split(cleanedUpEncodingMethod, ", "), "gzip")
		if !containsGzipMethod {
			return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n%s\r\n\r\n", connectionHeader)
		} else {
			return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Encoding: gzip\r\nContent-Type: text/plain\r\n%s\r\nContent-Length: %v\r\n\r\n%v", connectionHeader, buf.Len(), buf.String())
		}
	} else {
		requestBody := strings.Split(echoedString, " ")[0]
		return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n%s\r\nContent-Length: %v\r\n\r\n%v", connectionHeader, len(requestBody), requestBody)
	}

}
