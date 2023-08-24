package database

import (
	"fmt"
	"log"
)

func Test() {
	db := Open()
	defer db.Close()

	// _, err := db.Exec("insert into segments (name, created_at) values ('AVITO_VOICE_MESSAGES', $1)", time.Now())
	// if err != nil {
	// 	panic(err)
	// }

	rows, err := db.Query("select * from segments")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var created_at string
		var name string

		err := rows.Scan(&id, &created_at, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
}
