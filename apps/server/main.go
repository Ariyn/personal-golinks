package main

import (
	"encoding/json"
	"errors"
	"fmt"
	gl "github.com/ariyn/golinks"
	"github.com/boltdb/bolt"
	"io"
	"log"
	"net/http"
	"strings"
)

const RoutingBucketName = "routing"

type AddRequestBody struct {
	Name string
	Url  string
}

var db *bolt.DB
var router *gl.Router

func init() {
	router = gl.NewRouter()

	var err error
	db, err = bolt.Open("routings.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin(true)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	bucket, err := tx.CreateBucketIfNotExists([]byte(RoutingBucketName))
	if err != nil {
		log.Fatal(err)
	}

	err = bucket.ForEach(func(name []byte, url []byte) (err error) {
		return router.Redirect(string(name), string(url))
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func parseUrlAndBody(b []byte) (name, url string, err error) {
	var body AddRequestBody
	err = json.Unmarshal(b, &body)
	if err == nil {
		if body.Name != "" && body.Url != "" {
			return body.Name, body.Url, nil
		}

		return
	}

	twoLines := strings.Split(string(b), "\n")
	if len(twoLines) == 2 {
		return twoLines[0], twoLines[1], nil
	}

	return "", "", fmt.Errorf("unidentified format")
}

func saveToDB(db *bolt.DB, name, url string) (err error) {
	tx, err := db.Begin(true)
	if err != nil {
		return
	}
	defer tx.Rollback()

	bucket := tx.Bucket([]byte(RoutingBucketName))

	if url := bucket.Get([]byte(name)); url != nil {
		return fmt.Errorf("already exists key: %s - url: %s", name, url)
	}

	err = bucket.Put([]byte(name), []byte(url))
	if err != nil {
		return
	}

	return tx.Commit()
}

func main() {
	defer db.Close()

	err := router.Redirect("example", "https://example.org")
	if err != nil {
		log.Fatal(err)
	}

	err = router.Post("/", func(writer http.ResponseWriter, request *http.Request) (err error) {
		defer request.Body.Close()
		b, err := io.ReadAll(request.Body)
		if err != nil {
			return
		}

		name, url, err := parseUrlAndBody(b)
		if err != nil {
			return
		}

		if name == "/" || name == "" {
			return fmt.Errorf("you can't add to root")
		}

		err = saveToDB(db, name, url)
		if err != nil {
			return
		}

		err = router.Redirect(name, url)
		if err != nil {
			return
		}

		writer.WriteHeader(http.StatusOK)
		_, err = writer.Write([]byte("ok"))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	err = router.Get("/", func(writer http.ResponseWriter, request *http.Request) error {
		routeList := make([]string, 0)
		for _, r := range router.Routings {
			routeList = append(routeList, fmt.Sprintf("%s: %s", r.Method, r.Path))
		}

		b := []byte(strings.Join(routeList, "\n"))
		l, err := writer.Write(b)
		if err != nil {
			return err
		}

		if len(b) != l {
			return errors.New("not fully respond")
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = router.Put("/", func(writer http.ResponseWriter, request *http.Request) error {
		http.Error(writer, "not implemented yet", http.StatusTeapot)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = router.Delete("/", func(writer http.ResponseWriter, request *http.Request) error {
		http.Error(writer, "not implemented yet", http.StatusTeapot)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Addr:    ":80",
		Handler: http.HandlerFunc(router.Route),
	}

	log.Fatal(s.ListenAndServe())
}
