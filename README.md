## Go Simple Crud API [Gorilla/Mux, Postgress, sqlx]

#### Create Project folder
    mkdir gocrud;cd gocrud 

#### Add .gitignore and readme file
    touch .gitignore;touch README.md

### Install Dependencies
    go get -u github.com/gorilla/mux #Route the incoming HTTP request
    go get -u go get github.com/jmoiron/sqlx # database/sql extensions that helps access Database
    go get github.com/lib/pq #Golang Postgress Driver
    go get github.com/spf13/viper # Configurations manager

### Database Part
1. Install Postgres
2. check postgres status
    `sudo service postgresql <status>or<start>or<restart>or<stop>`
3. Login as postgres user
    `sudo -u postgres psql`
4. Create Seperate User and give access
    `CREATE USER username WITH SUPERUSER PASSWORD 'passwordstring';`
    `alter database dbName owner to userName;`
5. Login as new user
    `psql -U userName dbName`


###### Database migeation
- 1st you have to install `https://github.com/golang-migrate/migrate`

#### POSTGRES DB BONUS PART
###### Database commands 
- \?: show all psql commands.
- \h sql-command: show syntax on SQL command.
- \c dbname [username]: Connect to database, with an optional username (or \connect).
    Display Commands: You can append + to show more details.
- \l: List all database (or \list).
- \dt: Display all tables.
- \d table_name : Display tabes attribute
- \d+ table_name : Display describe tabes attribute 
- \di: Display all indexes.
- \dv: Display all views.
- \ds: Display all sequences.
- \dT: Display all types.
- \dS: Display all system tables.
- \du: Display all users.
- \x auto|on|off: Toggle|On|Off expanded output mode.

###### Commonly-used SQL Data Types
1. INT, SMALLINT: whole number. There is no UNSIGNED attribute in PostgreSQL.
2. SERIAL: auto-increment integer (AUTO_INCREMENT in MySQL).
3. REAL, DOUBLE: single and double precision floating-point number.
4. CHAR(n) and VARCHAR(n): fixed-length string of n characters and variable-length string of up to n characters. 
    String literals are enclosed by single quotes, e.g., 'Peter', 'Hello, world'.
5. NUMERIC(m,n): decimal number with m total digits and n decimal places (DECIMAL(m,n) in MySQL).
6. DATE, TIME, TIMESTAMP, INTERVAL: date and time.
7. User-defined types.
8. NULL: A special value indicates unknown value or no value (of an optional field), which is different from 0 and empty 
    string (that represent known value of 0 and empty string). To test for NULL value, use operator IS NULL or IS NOT NULL 
    (e.g., email IS NULL). Comparing two NULLs with = or != results in unknown.

### Now back to coding part

1. Create src folder [ all coding stuff located inside here ]
    `mkdir src;cd src`

2. Create `common` folder and add `config.go` file
    `mkdir common;touch config.go;touch config.json`

3. Create `database` folder and add `db.go` file
    `mkdir database;touch db.go`
    ```Go
    var Instance *sqlx.DB
    var err error

    func Connect(connectionString string) {
        Instance, err = sqlx.Open("postgres", connectionString)
        if err != nil {
            log.Fatal(err)
            panic("cannot connect to DB")
        }
        log.Println("connceted to DB")
    }
    ```

    - update `main.go`
        Load Configurations
        Initialize Database
        ```
        func main() {
            // Load Configurations
            common.LoadAppConfig()
            // Initialize Database
            database.Connect(common.AppConfig.DbConnectionString)
        }
    ```

4. Create `entities` folder and add `product.go` file
    `mkdir entities;touch product.go`
    ```Go
    type Product struct {
        ID          uint    `json:"id"`
        Name        string  `json:"name"`
        Price       float64 `json:"price"`
        Description string  `json:"description"`
    }
    ```


5. Create `migrations` folder and Create product table migration
- `migrate create -seq -ext=.sql -dir=./src/migrations product_table`

- Run Migration
`migrate -path=./src/migrations -database='postgres://username:password@localhost:5432/bankingapp?sslmode=disable' up`

6. Create `controllers` folder and add `product_controller.go` file
    `mkdir controllers;touch product_controller.go`
    ```Go
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
    ```

7. Create `routes` folder and add `routes.go` file
    `mkdir routes;touch routes.go`
    ```Go
    func RegisterProductRoutes(router *mux.Router) {
        router.HandleFunc("/api/v1/product", controllers.CreateProduct).Methods(http.MethodPost)
        router.HandleFunc("/api/v1/product", controllers.GetAllProducts).Methods(http.MethodGet)
        router.HandleFunc("/api/v1/product/{id}", controllers.GetProduct).Methods(http.MethodGet)
        router.HandleFunc("/api/v1/product/{id}", controllers.UpdateProduct).Methods(http.MethodPut)
        router.HandleFunc("/api/v1/product/{id}", controllers.DeleteProduct).Methods(http.MethodDelete)
    }
    ```

Final Touch
Update `main.go` file
```Go
func main() {
	// Load Configurations
	common.LoadAppConfig()
	// Initialize Database
	database.Connect(common.AppConfig.DbConnectionString)
	// Initialize the router
	+ router := mux.NewRouter().StrictSlash(true)
	+ // Register Routers
	+ routes.RegisterProductRoutes(router)

	+ // start the server
	+ log.Println(fmt.Sprintf("Starting Server on port %s", common.AppConfig.Port))
	+ log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%v", common.AppConfig.Port), router))

}
```

### Our Api Project is Complte
run :  `go run src/main.go`

### Testing
1. POST `http://localhost:8080/api/v1/product`
```
{
    "name": "product",
    "description": "product description",
    "price": 444.00
}
```
2. GET `http://localhost:8080/api/v1/product`
3. GET `http://localhost:8080/api/v1/product/1`
4. PUT `http://localhost:8080/api/v1/product/1`
  ```  
  {
    "name": "product-update",
    "description": "product description",
    "price": 444.00
  }
```
2. GET `http://localhost:8080/api/v1/product`
5. DELETE `http://localhost:8080/api/v1/product/1`



Our Folder Structure Look Like this
.
├── go.mod
├── go.sum
├── README.md
└── src
    ├── common
    │   ├── config.go
    │   └── config.json
    ├── controllers
    │   └── product_controller.go
    ├── database
    │   └── db.go
    ├── entities
    │   └── product.go
    ├── main.go
    ├── migrations
    │   └── migrations.go
    └── routes
        └── routes.go