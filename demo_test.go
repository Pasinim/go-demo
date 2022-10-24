package main

import (
	"context"
	"database/sql"
	"demo/utility"
	"github.com/google/go-cmp/cmp"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"testing"
)

type Item struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
	Sku  string `db:"sku"`
}

type Collection struct {
	Id       string
	Name     string
	Articles []Item
}

type demoDBContainer struct {
	instance testcontainers.Container
}

func initDemoDB(ctx context.Context, db sql.DB) error {
	const query = `CREATE SCHEMA demo;
   CREATE TABLE demo.articolo(
       id serial4 primary key not null,
       nome varchar(20), 
       sku varchar(20),
       collezione_id serial4 
       
       
       
   );`

	_, err := db.ExecContext(ctx, query)
	return err
}

func truncateDemoDB(ctx context.Context, db sql.DB) error {
	const query = `TRUNCATE demo.articolo`
	_, err := db.ExecContext(ctx, query)
	return err
}

func TestIntegratonDBInsertSelect(t *testing.T) {
	if testing.Short() {
		t.Skip("skipp")
	}

	ctx := context.Background()
	dbContainer := utility.NewTestDatabase(t)
	defer dbContainer.Close(t)
	db, err := sql.Open("postgres", dbContainer.ConnectionString(t))
	if err != nil {
		log.Fatal(err)
	}

	err = initDemoDB(ctx, *db)
	if err != nil {
		log.Fatal(err)
	}

	//aggiungo item
	item := Item{Id: 1, Name: "prova testing", Sku: "111"}
	const insertQuery = `INSERT INTO articolo (id, nome, sku ) values ($1, $2, $3)`
	_, err = db.ExecContext(
		ctx,
		insertQuery,
		item.Id,
		item.Name,
		item.Sku)
	if err != nil {
		t.Fatal(err)
	}

	//select
	savedItem := Item{Id: item.Id}
	const findQuery = `SELECT id, nome, sku FROM articolo where id = $1`
	row := db.QueryRowContext(ctx, findQuery, item.Id)
	err = row.Scan(&savedItem.Id, &savedItem.Name, &savedItem.Sku)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(item, savedItem) {
		t.Fatalf("Gli elementi non corrispondono:\n%s", cmp.Diff(item, savedItem))
	}
}