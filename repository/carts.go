package repository

import (
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
)

type CartRepository struct {
	db db.DB
}

func NewCartRepository(db db.DB) CartRepository {
	return CartRepository{db}
}

func (u *CartRepository) ReadCart() (model.Cart, error) {
	records, err := u.db.Load("carts")
	if err != nil {
		return model.Cart{}, err
	}

	if len(records) == 0 {
		return model.Cart{}, fmt.Errorf("Cart not found!")
	}

	var cart model.Cart
	err = json.Unmarshal(records, &cart)
	if err != nil {
		return model.Cart{}, err
	}

	return cart, nil
}

func (u *CartRepository) AddCart(cart model.Cart) error {
	var carts model.Cart

	jsonData, err := u.db.Load("carts")
	if err != nil {
		return err
	}

	if string(jsonData) == "" {
		curly := []byte("{}")
		jsonData = curly
	}

	err = json.Unmarshal(jsonData, &carts)
	if err != nil {
		return err
	}

	carts.Name = cart.Name
	carts.Cart = append(carts.Cart, cart.Cart...)
	carts.TotalPrice = cart.TotalPrice

	data, err := json.Marshal(carts)
	if err != nil {
		return err
	}

	err = u.db.Save("carts", data)
	return err
}

func (u *CartRepository) ResetCarts() error {
	err := u.db.Reset("carts", []byte(""))
	if err != nil {
		return err
	}

	return nil
}
