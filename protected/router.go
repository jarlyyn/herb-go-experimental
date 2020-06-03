package protected

import (
	"net/http"

	"github.com/herb-go/herb/user/identifier/httpidentifier"
)

type Server struct {
	Server    *http.Server
	Mux       *http.ServeMux
	Channels  *http.ServeMux
	Protecter *httpidentifier.Protecter
}

func NewServer() *Server {
	return &Server{}
}
func (s *Server) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	SetProtecter(r, s.Protecter)
	s.Mux.ServeHTTP(w, r)
}
