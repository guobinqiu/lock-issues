package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//前置知识：更新操作会在行上加互斥锁，只有当事务提交后才会释放锁互斥锁
//并发执行A,B两个事务
func main() {
	db, _ := sql.Open("mysql", "root:111111@tcp(localhost:3306)/go_hello?parseTime=true")
	defer db.Close()

	db.Exec("drop table customer")
	db.Exec("drop table ord")

	db.Exec(`CREATE TABLE customer (
  id int(10),
  username varchar(100),
  PRIMARY KEY (id)
  )`)
	db.Exec(`CREATE TABLE ord (
  id int(10),
  orderno varchar(100),
  PRIMARY KEY (id)
)`)

	db.Exec("insert into customer(id, username) value (?,?)", 1, "guobin")
	db.Exec("insert into ord(id, orderno) value (?,?)", 1, "001")

	c1 := make(chan bool)
	c2 := make(chan bool)

	go TransA(db, c1)
	go TransB(db, c2)

	<-c1
	<-c2
}

//事务A
func TransA(db *sql.DB, c chan bool) {
	tx, _ := db.Begin()
	fmt.Println("[TransA]开始事务")

	tx.Exec("update customer set username = ? where id = ?", "gb", 1)
	fmt.Println("[TransA]锁住customer.id=1这条记录")

	time.Sleep(10 * time.Millisecond)

	if _, err := tx.Exec("update ord set orderno = ? where id=?", "001", 1); err != nil {
		fmt.Println("[TransA]" + err.Error())
	} else {
		fmt.Println("[TransA]等待TransB释放锁")
	}

	tx.Commit()
	fmt.Println("[TransA]提交事务")

	c <- true
}

//事务B
func TransB(db *sql.DB, c chan bool) {
	tx, _ := db.Begin()
	fmt.Println("[TransB]开始事务")

	tx.Exec("update ord set orderno = ? where id=?", "002", 1)
	fmt.Println("[TransB]锁住ord.id=1这条记录")

	time.Sleep(10 * time.Millisecond)

	if _, err := tx.Exec("update customer set username = ? where id = ?", "guobin", 1); err != nil {
		fmt.Println("[TransB]" + err.Error())
	} else {
		fmt.Println("[TransB]等待TransA释放锁")
	}

	tx.Commit()
	fmt.Println("[TransB]提交事务")

	c <- true
}
