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

### Database Part
1. Install Postgress
2. check postgress status
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