package routes

import (
	"database/sql"
	"net/http"

	"github.com/NotAJocke/stars-viewer/internal/providers"
	"github.com/NotAJocke/stars-viewer/internal/templates/pages"
)

type App struct {
	Db *sql.DB
}

func (app *App) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// mux.HandleFunc("/public/main.css", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "internal/public/css/output.css")
	// })
	fs := http.FileServer(http.Dir("internal/public"))
	mux.Handle("/public/", http.StripPrefix("/public/", fs))

	mux.HandleFunc("/", app.MainRoute)

	return mux
}

func (app *App) MainRoute(w http.ResponseWriter, r *http.Request) {
	repos := providers.GetStarredRepos(app.Db, 1000)

	component := pages.HomePage(repos)
	component.Render(r.Context(), w)
}
