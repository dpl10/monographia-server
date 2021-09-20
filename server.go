package main

import (
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tkanos/gonfig"

	"github.com/dpl10/monographia-server/handler"
	"github.com/dpl10/monographia-server/model"
	"github.com/dpl10/monographia-server/util"
)

type (
	// Configuration to hold JSON key/value pairs for config/*.json
	Configuration struct {
		AngularPortHTTP       string
		AngularProxy          bool
		BodyLimit             string
		ConnectionStringMySQL string
		DatabaseConnections   int
		Host                  string
		JWTlife               time.Duration // hours
		PrivateKey            string
		PublicKey             string
		ServerPortHTTP        string
		ServerPortHTTPS       string
	}
)

func getConfigFile() string {
	env := os.Getenv("ENV")
	if len(env) == 0 {
		env = "development"
	}
	_, dirname, _, _ := runtime.Caller(0)
	return path.Join(filepath.Dir(dirname), "config", env+".json")
}

func main() {
	// config
	configuration := Configuration{}
	err := gonfig.GetConf(getConfigFile(), &configuration)
	if err != nil {
		panic(err)
	}

	// MySQL/Maria DB
	m := util.MySQL(configuration.ConnectionStringMySQL, configuration.DatabaseConnections)

	// allow database and interface function access
	h := handler.NewHandler(model.NewModel(m))

	// JWT setup and injection
	JWTkey, err := util.GenerateRandomBytes(32) // 256 bit signing key for HMAC-SHA256
	if err != nil {
		panic(err)
	}
	handler.JWTkey = JWTkey
	handler.JWTlife = configuration.JWTlife

	// web server
	e := echo.New()
	e.Pre(middleware.HTTPSRedirect())
	e.Use(middleware.BodyLimit(configuration.BodyLimit))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(util.BrotliWithConfig(util.BrotliConfig{ // switch to the offical version when it becomes available in v5 (probably)
		Level: 4, // setting based on https://blogs.akamai.com/2016/02/understanding-brotlis-potential.html
	}))
	if configuration.AngularProxy == true {
		angular := configuration.AngularPortHTTP
		if len(angular) == 0 {
			angular = "http://0.0.0.0:4200"
		}
		url, err := url.Parse(angular)
		if err != nil {
			panic(err)
		}
		target := []*middleware.ProxyTarget{
			{
				URL: url,
			},
		}
		e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(target)))
	}
	// unrestricted static
	e.File("/favicon.ico", "static/favicon.ico")
	// unrestricted api
	ua := e.Group("/api")
	ua.POST("/logon", h.Logon)
	ua.GET("/status", func(c echo.Context) error {
		return c.String(http.StatusOK, `{"Active":true,"Time":"`+time.Now().UTC().String()+`"}`)
	})
	// restricted api
	ra := ua.Group("/r")
	ra.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &util.JWTclaims{},
		SigningKey: JWTkey,
	}))
	ra.GET("/city", h.GetCity)
	// start web server
	go func(c *echo.Echo) {
		e.Logger.Fatal(e.Start(configuration.ServerPortHTTP))
	}(e)
	e.Logger.Fatal(e.StartTLS(configuration.ServerPortHTTPS, configuration.PublicKey, configuration.PrivateKey))
}
