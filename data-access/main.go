package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

const DATABASE_URL string = "DATABASE_URL"

type Album struct {
	ID int64
	Title string
	Artist string
	Price float64
}

func addAlbum(ctx context.Context, album Album, conn *pgx.Conn) (int64, error) {
	var id int64
	err := conn.QueryRow(
		ctx,
		`
		INSERT INTO album (title, artist, price)
		VALUES ($1, $2, $3)
		RETURNING id
		`,
		album.Title, album.Artist, album.Price,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %w", err)
	}
	return id, nil
}

func albumByID (ctx context.Context, id int64, conn *pgx.Conn) (Album, error) {
	album := Album{}
	row := conn.QueryRow(ctx, "SELECT id, title, artist, price FROM album WHERE ID = $1", id)
	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
        if err == sql.ErrNoRows {
            return album, fmt.Errorf("albumsById %d: no such album", id)
        }
        return album, fmt.Errorf("albumsById %d: %v", id, err)
    }

	return album, nil
}

func albumsByArtist(ctx context.Context, name string, conn *pgx.Conn) ([]Album, error) {
	albums := make([]Album, 0 , 8)

	rows, err := conn.Query(ctx, "SELECT id, title, artist, price FROM album WHERE artist = $1", name)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var album Album
		err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		albums = append(albums, album)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows failed: %w", rows.Err())
	}
	

	return albums, nil
}

func main() {
	err := godotenv.Load()

	if err !=  nil {
		log.Println("warn: .env not found (continue)")
	}

	url := os.Getenv(DATABASE_URL)
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, url)
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	
	defer conn.Close(ctx)

	newAlbum := Album{
		Title: "Satoshi Kosugi 2",
		Artist: "Satoshi Kosugi 2",
		Price: 90.00,
	}
	
	albumID, err := addAlbum(ctx, newAlbum, conn)

	

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID of added album: %v\n", albumID)
	
	fmt.Println("---")
	
	albums, err := albumsByArtist(ctx, "John Coltrane", conn)
	

	if err != nil {
		fmt.Fprintf(os.Stderr, "error")
		os.Exit(1)
	}

	album, err := albumByID(ctx, 1, conn)
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "error")
		os.Exit(1)
	}

	for _, a := range albums {
		fmt.Printf("%3d | %-20s | %-16s | %.2f\n", a.ID, a.Title, a.Artist, a.Price)
	}

	fmt.Println("---")

	fmt.Printf("%3d | %-20s | %-16s | %.2f\n", album.ID, album.Title, album.Artist, album.Price)
}
