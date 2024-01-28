package main

import (
	"os"

	"github.com/hasbiasshidiq/simple-profile/generated"
	"github.com/hasbiasshidiq/simple-profile/handler"
	"github.com/hasbiasshidiq/simple-profile/repository"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, newServer())

	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	dbDsn := os.Getenv("DATABASE_URL")
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})

	opts := handler.NewServerOptions{
		Repository: repo,
	}
	return handler.NewServer(opts)
}
