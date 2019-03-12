package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/productcatalogservice/genproto"
	_ "github.com/go-sql-driver/mysql"
)

//Database represents a... database!
type Database struct {
	user            string
	pass            string
	host            string
	port            string
	name            string
	DB              *sql.DB
	createStatement *sql.Stmt
	insertStatement *sql.Stmt
	selectStatement *sql.Stmt
}

//NewDatabase gives you a new database with initialized tables
func NewDatabase(user, pass, host, port, name string) (*Database, error) {
	var err error
	newDb := &Database{user, pass, host, port, name, nil, nil, nil, nil}
	newDb.DB, err = sql.Open("mysql", newDb.uri())
	if err != nil {
		return nil, fmt.Errorf("cannot create SQL connection: %v", err)
	}
	err = newDb.prepareStatements()
	return newDb, err
}

//Store products in the database
func (db *Database) Store(products []*pb.Product) error {
	for _, p := range products {
		data, _ := json.Marshal(p)
		_, err := db.insertStatement.Exec(p.Id, data)
		if err != nil {
			return err
		}
	}
	return nil
}

//Read products from the database
func (db *Database) Read() ([]*pb.Product, error) {
	var data []byte
	var product *pb.Product
	products := []*pb.Product{}

	rows, err := db.selectStatement.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&data)
		if err != nil {
			log.Warnf("database error retrieving product: %v", err)
		}
		err = json.Unmarshal(data, &product)
		if err != nil {
			log.Warnf("database error serializing product: %v", err)
		}
		products = append(products, product)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return products, nil
}

//Close closes the database connection
func (db *Database) Close() error {
	return db.DB.Close()
}

func (db *Database) uri() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db.user, db.pass, db.host, db.port, db.name)
}

func (db *Database) prepareStatements() error {
	var err error
	db.createStatement, err = db.DB.Prepare("CREATE TABLE IF NOT EXISTS product (id VARCHAR(64) PRIMARY KEY, data BLOB);")
	if err != nil {
		return err
	}
	// Exec create right away so we can prepare insert and select
	_, err = db.createStatement.Exec()
	if err != nil {
		return err
	}
	db.insertStatement, err = db.DB.Prepare("INSERT INTO product (id, data) VALUES (?, ?) ON DUPLICATE KEY UPDATE data = VALUES(data)")
	if err != nil {
		return err
	}
	db.selectStatement, err = db.DB.Prepare("SELECT data FROM product;")
	if err != nil {
		return err
	}
	return nil
}
