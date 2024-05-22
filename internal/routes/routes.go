package routes

import (
	"database/sql"
	"net/http"

	"github.com/NotAJocke/stars-viewer/internal/database"
	"github.com/NotAJocke/stars-viewer/internal/templates/pages"
)

type App struct {
	Db *sql.DB
}

func (app *App) NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.MainRoute)

	return mux
}

func (app *App) MainRoute(w http.ResponseWriter, r *http.Request) {
	repos := database.GetRepos(app.Db, 0)

	component := pages.HomePage(repos)
	component.Render(r.Context(), w)
}
