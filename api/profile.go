package api

import (
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	"text/template"
)

func (api *API) ImgProfileView(w http.ResponseWriter, r *http.Request) {
	filepath := path.Join("assets", "images", "img-avatar.png")
	http.ServeFile(w, r, filepath)
}

func (api *API) ImgProfileUpdate(w http.ResponseWriter, r *http.Request) {
	image := path.Join("assets", "images", "img-avatar.png")
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	f, _, err := r.FormFile("file-avatar")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	file, err := os.OpenFile(image, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if _, err := io.Copy(file, f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filepath := path.Join("views", "dashboard.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	cart, err := api.cartsRepo.ReadCart()
	listProducts, err := api.products.ReadProducts()
	data := model.Dashboard{
		Product: listProducts,
		Cart:    cart,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
	}
}
