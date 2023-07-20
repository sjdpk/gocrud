package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sjdpk/gocrud/src/database"
	"github.com/sjdpk/gocrud/src/entities"
)

// @desc : create product
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product entities.Product
	json.NewDecoder(r.Body).Decode(&product)
	query := "INSERT INTO products (name, price, description) VALUES ($1, $2, $3) RETURNING id"
	var id uint
	err := database.Instance.QueryRow(query, product.Name, product.Price, product.Description).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	product.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// @desc : list all product list
func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	var products []entities.Product
	query := "SELECT * FROM products"
	if err := database.Instance.Select(&products, query); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&products)
}

// @desc : get product
func GetProduct(w http.ResponseWriter, r *http.Request) {
	productId := mux.Vars(r)["id"]

	if !checkIfIdExists(productId) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Product not found")
		return
	}
	var product entities.Product
	query := "SELECT * FROM products WHERE id=$1"
	if err := database.Instance.Get(&product, query, productId); err != nil {
		log.Println("This is Error ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&product)
}

// @desc : update product
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productId := mux.Vars(r)["id"]
	if !checkIfIdExists(productId) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Product not found")
		return
	}

	var product entities.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate the UPDATE query dynamically based on provided fields
	var updateFields []string
	var values []interface{}

	if product.Name != "" {
		updateFields = append(updateFields, "name=$1")
		values = append(values, product.Name)
	}
	if product.Price != 0 {
		updateFields = append(updateFields, "price=$2")
		values = append(values, product.Price)
	}
	if product.Description != "" {
		updateFields = append(updateFields, "description=$3")
		values = append(values, product.Description)
	}

	// Check if any fields were provided for updating
	if len(updateFields) == 0 {
		http.Error(w, "No valid fields provided for update", http.StatusBadRequest)
		return
	}

	// Prepare the SQL update query
	query := fmt.Sprintf("UPDATE products SET %s WHERE id=$%d", strings.Join(updateFields, ", "), len(values)+1)
	values = append(values, productId)

	// Execute the update query
	_, err = database.Instance.Exec(query, values...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the updated product from the database
	query = "SELECT * FROM products WHERE id=$1"
	err = database.Instance.Get(&product, query, productId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&product)
}

// @desc : uDeletepdate product
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productId := mux.Vars(r)["id"]
	if !checkIfIdExists(productId) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Product not found")
		return
	}
	query := "DELETE FROM products WHERE id=$1"
	_, err := database.Instance.Exec(query, productId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode("delete sucess")
}

// check if id is present or not
func checkIfIdExists(id string) bool {
	var product entities.Product
	query := "SELECT * FROM products WHERE id=$1"
	err := database.Instance.Get(&product, query, id)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	if err == sql.ErrNoRows || product.ID == 0 {
		return false
	}
	return true
}
