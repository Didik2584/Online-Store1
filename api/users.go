package api

import (
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"path"
	"text/template"
	"time"
)

func (api *API) Register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	creds := model.Credentials{
		Username: username,
		Password: password,
	}

	if creds.Username == "" || creds.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Username or Password empty"})
		return
	}

	err = api.usersRepo.AddUser(creds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	filepath := path.Join("views", "status.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	var data = map[string]string{"name": creds.Username, "message": "register success!"}
	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
	}

}

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	creds := model.Credentials{
		Username: username,
		Password: password,
	}

	if creds.Username == "" || creds.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Username or Password empty"})
		return
	}

	listUser, err := api.usersRepo.ReadUser()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	if !api.usersRepo.LoginValid(listUser, creds) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Wrong User or Password!"})
		return
	}

	sessionKey := uuid.New().String()
	expiry := time.Now().Add(time.Hour * 5)

	cookie := http.Cookie{
		Path:    "/",
		Name:    "session_token",
		Value:   sessionKey,
		Expires: expiry,
	}
	http.SetCookie(w, &cookie)

	session := model.Session{
		Token:    sessionKey,
		Username: creds.Username,
		Expiry:   expiry,
	}
	err = api.sessionsRepo.AddSessions(session)

	filepath := path.Join("views", "dashboard.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	listProducts, err := api.products.ReadProducts()
	data := model.Dashboard{
		Product: listProducts,
		Cart: model.Cart{
			Name:       creds.Username,
			Cart:       []model.Product{},
			TotalPrice: 0,
		},
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
	}
}

func (api *API) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}
	sessionToken := cookie.Value

	api.sessionsRepo.DeleteSessions(sessionToken)

	cookie = &http.Cookie{
		Path:    "/",
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	}

	filepath := path.Join("views", "login.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
	}
}
