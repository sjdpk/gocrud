## Go Simple Crud API [Gorilla/Mux, Postgress, Gorm]

#### Create Project folder
    mkdir gocrud;cd gocrud 

#### Add .gitignore and readme file
    touch .gitignore;touch README.md

### Install Dependencies
    go get -u github.com/gorilla/mux #Route the incoming HTTP request
    go get -u gorm.io/gorm #ORM that helps access Database
    go get -u gorm.io/driver/postgres #Gorm Postgress Driver
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
    var Instance *gorm.DB
    var err error

    func Connect(connectionString string) {
        Instance, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
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


5. Create `migrations` folder and add `migrations.go` file
    `mkdir migrations;touch migrations.go`
    ```Go
    func Migration() {
        log.Println("ProductModel Migration Start...")
        database.Instance.AutoMigrate(&entities.Product{})
        log.Println("ProductModel Migration End")
    }
    ```
    - update `main.go`
    ```Go
    func main() {
        // Load Configurations
        common.LoadAppConfig()
        // Initialize Database
        database.Connect(common.AppConfig.DbConnectionString)
        + migrations.Migration()
    }
    ```

6. Create `controllers` folder and add `product_controller.go` file
    `mkdir controllers;touch product_controller.go`
    ```Go
    // @desc : create product
    func CreateProduct(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        var product entities.Product
        json.NewDecoder(r.Body).Decode(&product)
        database.Instance.Create(&product)
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(product)
    }
    // @desc : list all product list
    func GetAllProducts(w http.ResponseWriter, r *http.Request) {
        var product []entities.Product
        database.Instance.Find(&product)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(&product)
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
        database.Instance.First(&product, productId)
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
        database.Instance.First(&product, productId)
        json.NewDecoder(r.Body).Decode(&product)
        database.Instance.Save(&product)
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
        var product entities.Product
        database.Instance.Delete(&product, productId)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNoContent)
        json.NewEncoder(w).Encode("delete sucess")
    }

    // check if id is present or not
    func checkIfIdExists(id string) bool {
        var product entities.Product
        database.Instance.First(&product, id)
        if product.ID == 0 {
            return false
        }
        return true
    }

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
	migrations.Migration()

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
