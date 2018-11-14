// main.go

package main

import (
	"encoding/json"
	"log"
	"net/http"

	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
	"os"
	"os/signal"
	"syscall"
)

// global variable
var gConfigSetting map[string]interface{}
var router *gin.Engine

const (
	configFile = "config.json"
)

func deinit(sigs chan os.Signal) {
	fmt.Println("Deinit daemon start")
	sig := <-sigs
	fmt.Println(sig)
	fmt.Println("db disconnect")
	os.Exit(1)
}

func configParse() {
	// check config file size
	configInfo, err := os.Lstat(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// open config.json
	file, err := os.OpenFile(configFile, os.O_RDONLY, 0)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	// read all to data
	data := make([]byte, configInfo.Size())
	count, err := file.Read(data)
	if err != nil {
		log.Fatal(err, count)
	}

	//fmt.Printf("read %d bytes: %q\n", count, data[:count])

	if err := json.Unmarshal(data, &gConfigSetting); err != nil {
		panic(err)
	}
	fmt.Println(gConfigSetting["APP_NAME"])
	fmt.Println(gConfigSetting["DB_HOST"])
}

func main() {
	// config parser
	configParse()

	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	// Set the router as the default one provided by Gin
	router = gin.Default()

	// Set favicon.ico
	router.Use(favicon.New("public/static/photos/favicon.ico"))

	// Set sessions for keeping user info
	store := sessions.NewCookieStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	router.LoadHTMLGlob("public/templates/*")

	// Initialize the routes
	initializeRoutes()

	// add deinit when Ctrl+C to term process
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go deinit(sigs)

	// Start serving the application
	router.Run(":3000")
}

// Render one of HTML, JSON or CSV based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that
// the template name is present
func render(c *gin.Context, data gin.H, templateName string) {
	loggedInInterface, _ := c.Get("is_logged_in")
	data["is_logged_in"] = loggedInInterface.(bool)

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		// Respond with XML
		c.XML(http.StatusOK, data["payload"])
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}
}
