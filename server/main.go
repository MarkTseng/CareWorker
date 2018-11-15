// main.go

package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	"syscall"
)

// global const variable
const (
	configFile = "config.json"
	MongoDBUrl = "mongodb://localhost:27017/articles_demo_dev"
)

type careWorkerServer struct {
	// mgo objs
	db        *mgo.Database
	dbSession *mgo.Session
	articles  *mgo.Collection
	users     *mgo.Collection
	counters  *mgo.Collection
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
	cws.articles = cws.db.C("articles")
	cws.users = cws.db.C("users")
	cws.counters = cws.db.C("counters")
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

	// Set favicon.ico
	cws.router.Use(favicon.New("public/static/photos/favicon.ico"))

	// Set sessions for keeping user info
	store := sessions.NewCookieStore([]byte("secret"))
	cws.router.Use(sessions.Sessions("mysession", store))

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
	cws.router.Run(fmt.Sprintf(":%s", cws.ConfigSetting["APP_SERVER_PORT"]))
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
