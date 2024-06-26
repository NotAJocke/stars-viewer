package database

//
// import (
// 	"database/sql"
// 	"log"
//
// 	_ "github.com/mattn/go-sqlite3"
// )
//
// type Label struct {
// 	Name  string
// 	Color string
// 	Id    int
// }
//
// func CreateLabel(db *sql.DB, name string, color string) {
// 	tx, err := db.Begin()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
//
// 	q := `
//   INSERT INTO labels (name, color) VALUES (?, ?)
//   `
// 	stmt, err := tx.Prepare(q)
// 	if err != nil {
// 		tx.Rollback()
// 		log.Fatal(err)
// 	}
// 	defer stmt.Close()
//
// 	_, err = stmt.Exec(name, color)
//
// 	if err != nil {
// 		log.Println("Couldn't insert in DB")
// 		log.Fatal(err)
// 	}
//
// 	err = tx.Commit()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// }
//
// func DeleteLabelById(db *sql.DB, id int) {
// 	tx, err := db.Begin()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
//
// 	q1 := `DELETE FROM stars_labels WHERE label_id=?;`
// 	q2 := `DELETE FROM labels WHERE id=?;`
//
// 	stmt, err := tx.Prepare(q1)
// 	if err != nil {
// 		tx.Rollback()
// 		log.Fatal(err)
// 	}
// 	defer stmt.Close()
//
// 	_, err = stmt.Exec(id)
// 	if err != nil {
// 		tx.Rollback()
// 		log.Fatal(err)
// 	}
//
// 	stmt, err = tx.Prepare(q2)
// 	if err != nil {
// 		tx.Rollback()
// 		log.Fatal(err)
// 	}
// 	defer stmt.Close()
//
// 	_, err = stmt.Exec(id)
// 	if err != nil {
// 		tx.Rollback()
// 		log.Fatal(err)
// 	}
//
// 	err = tx.Commit()
// 	if err != nil {
// 		tx.Rollback()
// 		log.Fatalln(err)
// 	}
// }
//
// func GetLabels(db *sql.DB) []Label {
// 	q := `SELECT id, name, color FROM labels;`
//
// 	rows, err := db.Query(q)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	defer rows.Close()
//
// 	var labels []Label
//
// 	for rows.Next() {
// 		var label Label
//
// 		if err := rows.Scan(&label.Id, &label.Name, &label.Color); err != nil {
// 			log.Fatal(err)
// 		}
//
// 		labels = append(labels, label)
// 	}
//
// 	if err := rows.Err(); err != nil {
// 		log.Fatalln(err)
// 	}
//
// 	return labels
// }
//
// func AddLabelsToRepoByIds(db *sql.DB, repoId int, labelIds []int) {
// 	tx, err := db.Begin()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
//
// 	q := `INSERT INTO stars_labels (star_id, label_id) VALUES (?, ?);`
//
// 	stmt, err := tx.Prepare(q)
// 	if err != nil {
// 		tx.Rollback()
// 		log.Fatal(err)
// 	}
// 	defer stmt.Close()
//
// 	for _, id := range labelIds {
// 		_, err = stmt.Exec(repoId, id)
// 		if err != nil {
// 			tx.Rollback()
// 			log.Fatal(err)
// 		}
// 	}
//
// 	err = tx.Commit()
// 	if err != nil {
// 		tx.Rollback()
// 		log.Fatalln(err)
// 	}
// }
