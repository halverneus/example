package router

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/halverneus/example/api/file"
	"github.com/halverneus/example/api/user"
	"github.com/halverneus/example/config"
	"github.com/halverneus/example/database"
	"github.com/halverneus/example/lib/authenticate"
	"github.com/halverneus/example/lib/web"
	"github.com/halverneus/example/model"
)

// Run the web server. Function blocks.
func Run() (err error) {
	// Setup concurrent channels for shutdown and final error handling.
	errChan := make(chan error)
	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt)

	// Setup HTTP routes.
	router := httprouter.New()

	// V1 of the API.
	router.DELETE("/api/v1/file/*filepath", web.Wrap(authenticate.User(file.DELETE)))
	router.GET("/api/v1/file/*filepath", web.Wrap(authenticate.User(file.GET)))
	router.PUT("/api/v1/file/*filepath", web.Wrap(authenticate.User(file.PUT)))
	router.DELETE("/api/v1/user", web.Wrap(authenticate.User(user.DELETE)))
	router.POST("/api/v1/user", web.Wrap(authenticate.User(user.POST)))
	router.PUT("/api/v1/user", web.Wrap(authenticate.User(user.PUT)))

	// Latest version of the API.
	router.DELETE("/api/latest/file/*filepath", web.Wrap(authenticate.User(file.DELETE)))
	router.GET("/api/latest/file/*filepath", web.Wrap(authenticate.User(file.GET)))
	router.PUT("/api/latest/file/*filepath", web.Wrap(authenticate.User(file.PUT)))
	router.DELETE("/api/latest/user", web.Wrap(authenticate.User(user.DELETE)))
	router.POST("/api/latest/user", web.Wrap(authenticate.User(user.POST)))
	router.PUT("/api/latest/user", web.Wrap(authenticate.User(user.PUT)))

	// Setup HTTP server.
	server := &http.Server{Addr: config.Get.Example.Bind, Handler: router}

	// Start the file model and retrieve wait group to verify all deletions in
	// progress finish before exiting.
	wg := model.Start()

	// Start HTTP server.
	go func() {
		log.Printf("Server is listening at %s.\n", config.Get.Example.Bind)
		errChan <- server.ListenAndServe()
		close(errChan)
		database.Shutdown() // Start shutting down the database.
	}()

	// Wait for the server to be shutdown or to fail.
	select {
	case <-exitChan: // Expected exit. Provide time for uploads/downloads to finish.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		server.Shutdown(ctx)

	case err = <-errChan: // Server shutdown unexpectedly.
		return
	}

	// Wait for and retrieve error from shutdown server. Ignore 'server closed'
	// error (caused by interrupt).
	err = <-errChan
	if http.ErrServerClosed == err {
		err = nil
	}

	// Wait for all file deletions in progress before returning.
	wg.Wait()
	return
}
