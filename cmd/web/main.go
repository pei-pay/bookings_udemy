package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/pei-pay/bookings_udemy/pkg/config"
	"github.com/pei-pay/bookings_udemy/pkg/handlers"
	"github.com/pei-pay/bookings_udemy/pkg/render"
)

const portNumber = ":8080"

// templateCacheをglobalState?(wideConfig)に保存して、毎回作成する必要をなくす
var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {
	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	// development mode
	app.UseCache = false

	// TODO: これいるの?(まだ使ってない?)
	repo := handlers.NewRepo(&app)
	handlers.NewHanders(repo)

	// appの値をrenderFileでも使えるようにする
	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Starting application on port: %s", portNumber))
	// _ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
