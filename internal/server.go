package internal

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"test/config"
	"text/template"

	"github.com/gorilla/mux"
)

//Server describes file server struct
type Server struct {
	router *mux.Router
	config *config.Config
	mu     sync.Mutex
}

// New gets Apps params
func New(cfg *config.Config) *Server {

	return &Server{
		router: mux.NewRouter(),
		config: cfg,
	}
}

// Start method start test service
func (s *Server) Start() error {

	s.configureRouter()

	addr := s.config.Host + ":" + s.config.Port
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/api/download/{id}", s.downloadFile()).Methods("GET")
	s.router.HandleFunc("/api/delete/{id}", s.deleteFile()).Methods("DELETE")
	s.router.HandleFunc("/api/upload", s.uploadFile()).Methods("GET", "POST")
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
				log.Println("File not uploaded successfully", err)
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

		err := os.Remove(s.config.Directory + s.config.Route + u[:2] + "/" + u)
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

		input, err := os.Open(s.config.Directory + s.config.Route + u[:2] + "/" + u)

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
