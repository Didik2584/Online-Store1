package api

import (
	repo "a21hc3NpZ25tZW50/repository"
	"fmt"
	"net/http"
	"path"
)

type API struct {
	usersRepo    repo.UserRepository
	sessionsRepo repo.SessionsRepository
	products     repo.ProductRepository
	cartsRepo    repo.CartRepository
	mux          *http.ServeMux
}

type Page struct {
	File string
}

func (p Page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filepath := path.Join("views", p.File)
	fmt.Println("template: ", filepath)
	http.ServeFile(w, r, filepath)
}

func NewAPI(usersRepo repo.UserRepository, sessionsRepo repo.SessionsRepository, products repo.ProductRepository, cartsRepo repo.CartRepository) API {
	mux := http.NewServeMux()
	api := API{
		usersRepo,
		sessionsRepo,
		products,
		cartsRepo,
		mux,
	}

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	index := Page{File: "index.html"}
	
	mux.Handle("/", api.Get(index))
	mux.Handle("/page/register", api.Get(Page{File: "register.html"}))
	mux.Handle("/page/login", api.Get(Page{File: "login.html"}))
	mux.Handle("/user/register", api.Post(http.HandlerFunc(api.Register)))
	mux.Handle("/user/login", api.Post(http.HandlerFunc(api.Login)))
	mux.Handle("/user/logout", api.Get(api.Auth(http.HandlerFunc(api.Logout))))
	mux.Handle("/user/img/profile", api.Get(api.Auth(http.HandlerFunc(api.ImgProfileView))))
	mux.Handle("/user/img/update-profile", api.Post(api.Auth(http.HandlerFunc(api.ImgProfileUpdate))))
	mux.Handle("/cart/add", api.Post(api.Auth(http.HandlerFunc(api.AddCart))))

	return api
}

func (api *API) Handler() *http.ServeMux {
	return api.mux
}

func (api *API) Start() {
	fmt.Println("starting web server at http://localhost:8080")
	http.ListenAndServe(":8080", api.Handler())
}
