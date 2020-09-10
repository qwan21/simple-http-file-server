package internal

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"test/config"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

//Server describes file server struct
type Server struct {
	mu          sync.Mutex
	router      *mux.Router
	config      *config.Config
	server      *http.Server
	shutdownReq chan bool
}

// New gets Apps params
func New(cfg *config.Config) *Server {

	return &Server{
		router:      mux.NewRouter(),
		config:      cfg,
		shutdownReq: make(chan bool),
	}
}

// Start method start test service
func (s *Server) Start() {

	s.configureRouter()

	s.server = &http.Server{
		Addr:    s.config.Host + ":" + s.config.Port,
		Handler: s.router,
	}

	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			log.Println("Listen and serve:", err)
		}

	}()
	s.waitShutdown()

}

//waitShutdown waits shutdownReq
func (s *Server) waitShutdown() {

	<-s.shutdownReq
	log.Println("Shutdown...")

	//Create shutdown context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)

	if err != nil {
		log.Println("Shutdown request error:", err)
	}

	log.Println("Stopping http file server...")

}
func (s *Server) configureRouter() {
	s.router.HandleFunc("/api/download/{id}", s.downloadFile()).Methods("GET")
	s.router.HandleFunc("/api/delete/{id}", s.deleteFile()).Methods("DELETE")
	s.router.HandleFunc("/api/upload", s.uploadFile()).Methods("GET", "POST")
	s.router.HandleFunc("/api/shutdown", s.shutdown()).Methods("GET", "POST")
}

func (s *Server) shutdown() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Shutdown server"))
		go func() {
			s.shutdownReq <- true
		}()
	}

}

func (s *Server) uploadFile() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			t, _ := template.ParseFiles("fileUp.html")
			t.Execute(w, nil)
		} else {

			log.Println("Uploading file...")
			input, handler, err := r.FormFile("myFile")
			if err != nil {
				log.Println("Error Retrieving the File", err)
				return
			}
			defer input.Close()

			hash, err := s.upload(input, handler)

			if err != nil {
				log.Println("File not uploaded successfully:", err)
				w.Write([]byte(err.Error() + "\r\n"))
				return
			}

			w.Write([]byte(hash + "\r\n"))

			w.Write([]byte(http.StatusText(http.StatusOK) + "\r\n"))
			log.Println("File was uploaded successfully")
		}
	}
}

func (s *Server) deleteFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		log.Println("Deleting file...")
		u := path.Base(r.URL.String())

		filePath := s.config.Directory + s.config.Route + u[:2] + "/" + u
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			//path does not exist
			log.Printf("This file not found\r\n")
			w.Write([]byte(fmt.Errorf("This file not found").Error()))
			return
		}

		err := os.Remove(filePath)
		if err != nil {
			log.Println(err)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(http.StatusText(http.StatusOK)))
		log.Println("File was deleted successfully")
	}
}

func (s *Server) downloadFile() http.HandlerFunc {

	// Get the data
	return func(w http.ResponseWriter, r *http.Request) {

		log.Println("Downloading file...")
		u := path.Base(r.URL.String())

		filePath := s.config.Directory + s.config.Route + u[:2] + "/" + u

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			//path does not exist
			log.Printf("This file not found\r\n")
			w.Write([]byte(fmt.Errorf("This file not found").Error()))
			return
		}

		input, err := os.Open(filePath)

		if err != nil {
			log.Println(err)
			w.Write([]byte(err.Error()))
			return
		}
		defer input.Close()

		w.Header().Set("Content-Disposition", "attachment; filename="+u)
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.Header().Set("Content-Length", r.Header.Get("Content-Length"))

		if _, err := io.Copy(w, input); err != nil {
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			return
		}
		w.Write([]byte(http.StatusText(http.StatusOK)))

		log.Println("File was downloaded successfully")
	}
}
