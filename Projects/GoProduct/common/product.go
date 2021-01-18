package common

import "fmt"

func (s *sqliteHandler) createProducts(name string, sessionid string) {
	stmt, err := s.db.Prepare("INSERT INTO products (name, sessionid, createdat)  VALUES(?,?, datetime('now')")
	if err != nil {
		panic(err)
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, sessionid)
	if err != nil {
		panic(err)
	}

	rowcnt, err := result.RowsAffected()
	if err != nil {
		fmt.Println("[ERR] RowsAffected err : ", err)
		panic(err)
	}

	if rowcnt == 0 {
		fmt.Println("[LOG] anything is not affected")
	}

}

func (s *sqliteHandler) readProducts(sessionid string) []*Products {
	// read html file
	productlist := []*Products{}
	rows, err := s.db.Query("SELECT id, name, sessionid FROM products WHERE sessionid=?", sessionid)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var product Products
		rows.Scan(&product.ID, &product.Name, &product.SessionID)
		productlist = append(productlist, &product)
	}
	return productlist
}

func (s *sqliteHandler) updateProducts(name string, sessionid string) {
	stmt, err := s.db.Prepare("UPDATE products SET name=? WHERE sessionid=?")
	if err != nil {
		panic(err)
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, sessionid)
	if err != nil {
		panic(err)
	}

	rowcnt, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}

	if rowcnt == 0 {
		fmt.Println("[LOG] anything is not affected")
	}

}

func (s *sqliteHandler) deleteProducts(name string, sessionid string) {
	stmt, err := s.db.Prepare("DELETE FROM products where name=? AND sessionid=?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, sessionid)
	if err != nil {
		panic(err)
	}

	rowcnt, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}

	if rowcnt == 0 {
		fmt.Println("[LOG] anything is not affected")
	}
}
