package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	connection string
	host       string
	port       int
	dbname     string
	user       string
	password   string
)

// init function called by Go before main execution and after variables definition
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connection = os.Getenv("DB_CONNECTION")
	host = os.Getenv("DB_HOST")
	port, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Printf("Error: %v\nUsing default port\n", err)
		port = 5432
	}
	dbname = os.Getenv("DB_DATABASE")
	user = os.Getenv("DB_USERNAME")
	password = os.Getenv("DB_PASSWORD")
}

func connectDatabase() *sql.DB {
	// Define connection string for lib/pq
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	fmt.Println(psqlInfo)

	// Open db connection
	db, err := sql.Open(connection, psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	// Check db connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create table and insert default data if not exist
	err = initDb(db)
	if err != nil {
		log.Fatalf("Error during db initialization: %v\n", err)
	}

	return db
}

func initDb(db *sql.DB) error {
	var products = []Product{
		{
			Name:        "Cyberpunk 2077",
			Description: "Cyberpunk 2077 est un jeu vidéo d''action-RPG en vue à la première personne développé par CD Projekt RED, fondé sur la série de jeu de rôle sur table Cyberpunk 2020 conçue par Mike Pondsmith.",
			Quantity:    100,
			Weight:      10.0,
			Price:       50.0,
			AssetUrl:    "https://images-na.ssl-images-amazon.com/images/I/81%2BpdXH3fjL._AC_SY500_.jpg",
		},
		{
			Name:        "Assassin''s Creed Valhalla",
			Description: "Assassin''s Creed Valhalla est un jeu vidéo d''action-aventure et de rôle, développé par Ubisoft Montréal et édité par Ubisoft, sorti en novembre 2020 sur Microsoft Windows.",
			Quantity:    100,
			Weight:      8.0,
			Price:       59.99,
			AssetUrl:    "https://image.jeuxvideo.com/medias/158826/1588264397-5261-jaquette-avant.jpg",
		},
		{
			Name:        "Red Dead Redemption II",
			Description: "Red Dead Redemption II est un jeu vidéo d''action-aventure et de western multiplateforme, développé par Rockstar Studios et édité par Rockstar Games, sorti le 26 octobre 2018 sur PlayStation 4 et Xbox One et le 5 novembre 2019 sur Windows.",
			Quantity:    100,
			Weight:      10.0,
			Price:       50.0,
			AssetUrl:    "https://lh3.googleusercontent.com/HCUkD69MAHEOj84Yi7Kb5vxHpCePTsmQI4g9vYuVPUo-87cWE6ZZIk0tiyYzaiS9zaAFMTXRNYJaaRczRN-yQYw",
		},
	}

	_, err := db.Query(
		"CREATE TABLE IF NOT EXISTS products (" +
			"id SERIAL," +
			"name varchar(40) NOT NULL PRIMARY KEY," +
			"description text NOT NULL," +
			"quantity integer NOT NULL," +
			"weight real NOT NULL," +
			"price real NOT NULL," +
			"asset_url text)",
	)
	if err != nil {
		return err
	}

	for _, p := range products {
		_, err = db.Query(
			"INSERT INTO products (name, description, quantity, weight, price, asset_url) VALUES ('" +
				p.Name + "', '" +
				p.Description + "', " +
				fmt.Sprintf("%d", p.Quantity) + ", " +
				fmt.Sprintf("%f", p.Weight) + ", " +
				fmt.Sprintf("%f", p.Price) + ", '" +
				p.AssetUrl + "') " +
				"ON CONFLICT DO NOTHING",
		)
		if err != nil {
			return err
		}
	}

	return nil
}
