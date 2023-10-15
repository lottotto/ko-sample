package kosample

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Server struct {
	http.Server
	DBPool *sqlx.DB
}

func NewServer() *Server {
	db := GetConnect()

	return &Server{
		Server: http.Server{
			Addr: ":8080",
		},
		DBPool: db,
	}

}
