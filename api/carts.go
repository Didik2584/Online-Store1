package api

import (
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"text/template"
)

func (api *API) AddCart(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Request Product Not Found"})
		return
	}

	product := r.Form.Get("product")
	if product == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Request Product Not Found"})
		return
	}

	var totalPrice float64
	var list []model.Product
	for _, formList := range r.Form {
		for _, v := range formList {
			item := strings.Split(v, ",")
			p, _ := strconv.ParseFloat(item[2], 64)
			q, _ := strconv.ParseFloat(item[3], 64)
			total := p * q
			totalPrice += total
			list = append(list, model.Product{
				Id:       item[0],
				Name:     item[1],
				Price:    item[2],
				Quantity: item[3],
				Total:    total,
			})
		}
	}

	cart := model.Cart{
		Name:       fmt.Sprintf("%s", r.Context().Value("username")),
		Cart:       list,
		TotalPrice: totalPrice,
	}

	api.cartsRepo.AddCart(cart)

	filepath := path.Join("views", "dashboard.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

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
