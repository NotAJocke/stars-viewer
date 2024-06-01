package routes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/NotAJocke/stars-viewer/internal/github"
	"github.com/NotAJocke/stars-viewer/internal/providers"
	"github.com/NotAJocke/stars-viewer/internal/templates/components"
	"github.com/NotAJocke/stars-viewer/internal/templates/pages"
)

type App struct {
	Db    *sql.DB
	Cache []github.StarredRepo
}

func (app *App) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("internal/public"))
	mux.Handle("/public/", http.StripPrefix("/public/", fs))

	mux.HandleFunc("PUT /update", app.Update)
	mux.HandleFunc("DELETE /unstar", app.Unstar)
	mux.HandleFunc("/", app.MainRoute)

	return mux
}

func (app *App) MainRoute(w http.ResponseWriter, r *http.Request) {
	repos := providers.GetStarredRepos(app.Db, 30)

	app.Cache = repos

	component := pages.HomePage(repos)
	component.Render(r.Context(), w)
}

func (app *App) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		repos, err := providers.UpdateCache(app.Db)
		if err != nil {
			log.Fatalln(fmt.Errorf("Couldn't update cache from UI. %w", err))
		}

		repos = append(repos, app.Cache...)
		app.Cache = repos

		reposComponent := components.StarredRepos(repos)
		reposComponent.Render(r.Context(), w)

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (app *App) Unstar(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		if err := r.ParseForm(); err != nil {
			log.Println(fmt.Errorf("[WARN] couldn't parse form from unstar request. %w", err))
		}

		fullName := r.FormValue("fullname")
		providers.UnstarRepo(app.Db, fullName)

		fmt.Printf("%s\n", fullName)

		var repos []github.StarredRepo
		for _, repo := range app.Cache {
			if repo.FullName != fullName {
				repos = append(repos, repo)
			}
		}

		app.Cache = repos

		reposComponent := components.StarredRepos(repos)
		reposComponent.Render(r.Context(), w)

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

}
