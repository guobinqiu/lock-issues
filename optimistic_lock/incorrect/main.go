package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const (
	username = "guobin"
)

const schema = `
CREATE TABLE account (
  username varchar(100),
  balance int(11)
)
`

//账户
type account struct {
	username string
	balance  int
}

//并发执行，结果900，1000，1100都有可能，以最后一个提交的事务为准
func main() {
	db, _ := sql.Open("mysql", "root:111111@tcp(localhost:3306)/go_hello?parseTime=true")
	defer db.Close()

	db.Exec("drop table account")
	db.Exec(schema)

	db.Exec("insert into account(username, balance) value (?,?)", username, 1000)
	fmt.Println("余额：1000")

	c1 := make(chan bool)
	c2 := make(chan bool)

	go Withdraw(db, c1)
	go Transfer(db, c2)

	<-c1
	<-c2

	a := queryAccount(db)
	fmt.Println("余额:", a.balance)
}

//取款事务
func Withdraw(db *sql.DB, c chan bool) {
	a := queryAccount(db)
	fmt.Println("[Withdraw]查询到存款余额:", a.balance)

	time.Sleep(10 * time.Millisecond)

	tx, _ := db.Begin()
	fmt.Println("[Withdraw]开始事务")

	//取出100
	tx.Exec("update account set balance = ? where username = ?", a.balance-100, username)
	fmt.Println("[Withdraw]取出100，余额变成:", a.balance-100)

	tx.Commit()
	fmt.Println("[Withdraw]提交事务")

	c <- true
}

//转账事务
func Transfer(db *sql.DB, c chan bool) {
	a := queryAccount(db)
	fmt.Println("[Transfer]查询到存款余额:", a.balance)

	time.Sleep(10 * time.Millisecond)

	tx, _ := db.Begin()
	fmt.Println("[Transfer]开始事务")

	//存入100
	tx.Exec("update account set balance = ? where username = ?", a.balance+100, username)
	fmt.Println("[Transfer]存入100，余额变成", a.balance+100)

	tx.Commit()
	fmt.Println("[Transfer]提交事务")

	c <- true
}

func queryAccount(db *sql.DB) account {
	a := account{}
	db.QueryRow("select username, balance from account where username = ?", username).Scan(&a.username, &a.balance)
	return a
}
