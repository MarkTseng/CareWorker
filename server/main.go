// main.go

package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
	"gopkg.in/mgo.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// global const variable
const (
	configFile = "config.json"
	MongoDBUrl = "mongodb://localhost:27017/articles_demo_dev"
)

type careWorkerServer struct {
	// mgo objs
	db         *mgo.Database
	dbSession  *mgo.Session
	collection map[string]*mgo.Collection
	// gin objs
	router        *gin.Engine
	userRoutes    *gin.RouterGroup
	articleRoutes *gin.RouterGroup
	// private objs
	ConfigSetting map[string]interface{}
}

func deinit(sigs chan os.Signal, cws *careWorkerServer) {
	fmt.Println("Deinit daemon start")
	sig := <-sigs
	fmt.Println(sig)
	fmt.Println("db disconnect")
	cws.dbSession.Close()
	os.Exit(1)
}

func DBconnect(cws *careWorkerServer) {
	session, err := mgo.Dial(MongoDBUrl)
	if err != nil {
		panic(err)
		os.Exit(2)
	}

	cws.dbSession = session
	dbName := fmt.Sprintf("%s", cws.ConfigSetting["DB_DATABASE"])
	db := session.DB(dbName)
	cws.db = db

	// parser BDs collections
	cws.collection = make(map[string]*mgo.Collection)
	collections := fmt.Sprintf("%s", cws.ConfigSetting["DB_COLLECTIONS"])
	for _, collection := range strings.Split(collections[1:len(collections)-1], " ") {
		fmt.Printf("create %s collection\n", collection)
		cws.collection[collection] = cws.db.C(collection)
	}

}

func configParse(cws *careWorkerServer) {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if err := json.Unmarshal(content, &cws.ConfigSetting); err != nil {
		panic(err)
	}

	for k := range cws.ConfigSetting {
		fmt.Printf("%s=%s\n", k, cws.ConfigSetting[k])
	}
}

func genSaltString() string {
	buf := new(bytes.Buffer)
	io.CopyN(buf, rand.Reader, 32)
	return hex.EncodeToString(buf.Bytes())
}

func DoHash(pass, salt string) string {
	h := sha256.New()
	h.Write([]byte(pass))
	h.Write([]byte(salt))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	// initial cws object
	cws := new(careWorkerServer)

	// config parser
	configParse(cws)

	// connect to DB server
	DBconnect(cws)

	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	// Set the router as the default one provided by Gin
	//router = gin.Default()
	cws.router = gin.Default()
	// secure json prefix for angularjs
	cws.router.SecureJsonPrefix(")]}',\n")

	// Set favicon.ico
	cws.router.Use(favicon.New("public/static/photos/favicon.ico"))

	// Set sessions for keeping user info
	store := sessions.NewCookieStore([]byte("secretSession"))
	cws.router.Use(sessions.Sessions("careWorkerSession", store))

	// static html and angularjs
	cws.router.Use(static.Serve("/", static.LocalFile("public", true)))

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	cws.router.LoadHTMLGlob("public/templates/*")

	// Initialize the routes
	initializeRoutes(cws)

	// add deinit when Ctrl+C to term process
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go deinit(sigs, cws)

	// Start serving the application
	//cws.router.Run(fmt.Sprintf(":%s", cws.ConfigSetting["APP_SERVER_PORT"]))
	cws.router.RunTLS(fmt.Sprintf(":%s", cws.ConfigSetting["APP_SERVER_PORT"]), "server/ssldata/certificate.crt", "server/ssldata/private.key")
}

// Render one of HTML, JSON or CSV based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that
// the template name is present
func render(c *gin.Context, data gin.H, templateName string) {
	loggedInInterface, _ := c.Get("is_logged_in")
	data["is_logged_in"] = loggedInInterface.(bool)

	//log.Printf("render Request.Header: %s\n", c.Request.Header.Get("Accept"))
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.SecureJSON(http.StatusOK, data["payload"])
		//c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		// Respond with XML
		c.XML(http.StatusOK, data["payload"])
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}
}
