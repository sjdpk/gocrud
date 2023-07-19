## Go-Lang Simple Crud API [Gorilla/Mux, Postgress, Gorm]

#### Create Project folder
    mkdir gocrud;cd gocrud 

#### Add .gitignore and readme file
    touch .gitignore;touch README.md

### Install Dependencies
###### Route the incoming HTTP request 
    go get -u github.com/gorilla/mux
###### ORM that helps access Database
    go get -u gorm.io/gorm
###### Gorm Postgress Driver
    go get -u gorm.io/driver/postgres
###### Configurations manager that helps us to load file and environment values
    go get github.com/spf13/viper

