package routes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

	mux.HandleFunc("POST /search", app.Search)
	mux.HandleFunc("PUT /update", app.Update)
	mux.HandleFunc("DELETE /unstar/{fullname...}", app.Unstar)
	mux.HandleFunc("GET /repos/{i}", app.Repos)
	mux.HandleFunc("/", app.MainRoute)

	return mux
}

func (app *App) Repos(w http.ResponseWriter, r *http.Request) {
	i := r.PathValue("i")

	if i != "" {
		currentIndex, err := strconv.Atoi(i)
		if err != nil {
			return
		}

		fmt.Println(currentIndex)

		repos := providers.GetStarredRepos(app.Db, 30, currentIndex)

		if len(repos) != 0 {
			app.Cache = append(app.Cache, repos...)

			component := components.LazyStarredRepos(repos, len(app.Cache))
			component.Render(r.Context(), w)
		}
	}
}

func (app *App) Search(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("search")

	var repos []github.StarredRepo
	if q == "" {
		repos = app.Cache
	} else {
		repos = providers.SearchRepos(app.Db, q, 30)
	}

	component := components.StarredRepos(repos)
	component.Render(r.Context(), w)
}

func (app *App) MainRoute(w http.ResponseWriter, r *http.Request) {
	repos := providers.GetStarredRepos(app.Db, 30, 0)

	app.Cache = repos

	component := pages.HomePage(repos)
	component.Render(r.Context(), w)
}

func (app *App) Update(w http.ResponseWriter, r *http.Request) {
	repos, err := providers.UpdateCache(app.Db)
	if err != nil {
		log.Fatalln(fmt.Errorf("Couldn't update cache from UI. %w", err))
	}

	repos = append(repos, app.Cache...)
	app.Cache = repos

	reposComponent := components.StarredRepos(repos)
	reposComponent.Render(r.Context(), w)
}

func (app *App) Unstar(w http.ResponseWriter, r *http.Request) {
	fullName := r.PathValue("fullname")
	providers.UnstarRepo(app.Db, fullName)

	var repos []github.StarredRepo
	for _, repo := range app.Cache {
		if repo.FullName != fullName {
			repos = append(repos, repo)
		}
	}

	app.Cache = repos

	reposComponent := components.StarredRepos(repos)
	reposComponent.Render(r.Context(), w)

}
