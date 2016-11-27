package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/colindev/osenv"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
)

var (
	jobs  *gorm.DB
	cards *gorm.DB
	env   = &environments{}
)

func init() {

	flag.StringVar(&env.path, "env", ".env", "env file")
	flag.Parse()

	if err := godotenv.Load(env.path); err != nil {
		log.Fatal(err)
	}

	if err := osenv.LoadTo(env); err != nil {
		log.Fatal(err)
	}

	log.Println(env)

	db, err := gorm.Open("sqlite3", env.DSN)
	if err != nil {
		log.Fatal(err)
	}

	if env.Debug {
		db = db.Debug()
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}
	jobs = db.Model(Job{})
	cards = db.Model(Card{})

	if err := db.AutoMigrate(Job{}, Card{}).Error; err != nil {
		log.Fatal(err)
	}
}

func main() {

	go func() {
		api := rest.NewApi()
		api.Use(rest.DefaultDevStack...)
		router, err := rest.MakeRouter(
			rest.Get("/jobs", getJobs),
			rest.Get("/cards", getCards),
		)
		if err != nil {
			log.Fatal(err)
		}
		api.SetApp(router)

		server := http.NewServeMux()
		listener, err := net.Listen("tcp", env.PublicAddr)
		if err != nil {
			log.Fatal(err)
		}
		server.Handle("/ui/", http.StripPrefix("/ui", http.FileServer(http.Dir(env.PublicUI))))
		server.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
		log.Println("admin server:", http.Serve(listener, server))
	}()
	go func() {
		api := rest.NewApi()
		api.Use(rest.DefaultDevStack...)
		router, err := rest.MakeRouter(
			rest.Get("/jobs", getJobs),
			rest.Post("/job", postJob),
			rest.Delete("/job/:name", deleteJob),
			rest.Get("/cards", getCards),
			rest.Post("/card", postCard),
			rest.Delete("/card/:name", deleteCard),
		)
		if err != nil {
			log.Fatal(err)
		}
		api.SetApp(router)

		server := http.NewServeMux()
		listener, err := net.Listen("tcp", env.AdminAddr)
		if err != nil {
			log.Fatal(err)
		}
		server.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))
		log.Println("admin server:", http.Serve(listener, server))
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGQUIT)
	log.Printf("shutdown with: %s\n", <-shutdown)
}

// Error ...
type Error struct {
	Content error `json:"error"`
}

// Ok ...
type Ok struct {
	Content interface{} `json:"success"`
}

var ok = Ok{"ok"}

// job
func getJobs(w rest.ResponseWriter, r *rest.Request) {
	var (
		db    = jobs
		list  []Job
		query = r.URL.Query()
	)

	if classes := strings.TrimSpace(query.Get("classes")); classes != "" {
		db = db.Where("class IN (?)", strings.Split(classes, ","))
	}

	if err := db.Find(&list).Error; err != nil {
		w.WriteJson(Error{err})
		return
	}

	w.WriteJson(list)
}
func postJob(w rest.ResponseWriter, r *rest.Request) {
	var (
		job Job
	)

	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		w.WriteJson(Error{err})
		return
	}

	if err := jobs.Save(job).Error; err != nil {
		w.WriteJson(Error{err})
		return
	}

	w.WriteJson(job)
}
func deleteJob(w rest.ResponseWriter, r *rest.Request) {
	name, err := url.QueryUnescape(r.PathParams["name"])
	if err != nil {
		w.WriteJson(Error{err})
	}

	if err := jobs.Where("name = ?", name).Delete(Job{}).Error; err != nil {
		w.WriteJson(Error{err})
		return
	}

	w.WriteJson(ok)
}

// card
func getCards(w rest.ResponseWriter, r *rest.Request) {
	var (
		db    = cards
		list  []Card
		query = r.URL.Query()
	)

	if classes := strings.TrimSpace(query.Get("classes")); classes != "" {
		db = db.Where("class IN (?)", strings.Split(classes, ","))
	}

	if elements := strings.TrimSpace(query.Get("elements")); elements != "" {
		db = db.Where("element IN (?)", strings.Split(elements, ","))
	}

	if err := db.Find(&list).Error; err != nil {
		w.WriteJson(Error{err})
		return
	}

	w.WriteJson(list)
}
func postCard(w rest.ResponseWriter, r *rest.Request) {
	var (
		card Card
	)

	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		w.WriteJson(Error{err})
		return
	}

	if err := jobs.Save(card).Error; err != nil {
		w.WriteJson(Error{err})
		return
	}

	w.WriteJson(card)
}
func deleteCard(w rest.ResponseWriter, r *rest.Request) {
	name, err := url.QueryUnescape(r.PathParams["name"])
	if err != nil {
		w.WriteJson(Error{err})
	}

	if err := jobs.Where("name = ?", name).Delete(Card{}).Error; err != nil {
		w.WriteJson(Error{err})
		return
	}

	w.WriteJson(ok)
}
