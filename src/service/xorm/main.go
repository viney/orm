package main

import (
	"errors"
	"fmt"
	"sync"
	"time"

	_ "engine/pq"
	"engine/xorm"
)

var (
	mux    sync.Mutex
	engine *xorm.Engine
)

type Users struct {
	Uid     uint64 `xorm:"pk"`
	Name    string
	Email   string    `xorm:"text"`
	Created time.Time `xorm:"created"`
}

func init() {
	mux.Lock()
	defer mux.Unlock()

	// check
	if engine != nil {
		return
	}

	pq, err := xorm.NewEngine("postgres", "host=192.168.1.241 port=4932 user=viney password=admin dbname=test sslmode=disable")
	if err != nil {
		panic(err)
	}

	// new db
	engine = pq

	// sync
	if err := engine.Sync(new(Users)); err != nil {
		panic(err)
	}

	// open log
	engine.ShowSQL = true
	engine.ShowErr = true
	engine.ShowDebug = true
	engine.ShowWarn = true
}

// insert data to table
func Insert() error {
	user := &Users{
		Uid:     1,
		Name:    "viney",
		Email:   "viney.chow@gmail.com",
		Created: time.Now(),
	}

	id, err := engine.InsertOne(user)
	if err != nil {
		return err
	} else if id <= 0 {
		return errors.New("插入失败")
	}

	return nil
}

// update data of table
func Update() error {
	user := &Users{
		Name:    "维尼",
		Created: time.Now(),
	}

	i, err := engine.Id(1).Update(user)
	if err == nil {
		return nil
	} else if i <= 0 {
		return errors.New("更新失败")
	}

	return nil
}

// query data from table
func Query() error {
	user := &Users{}
	b, err := engine.Id(1).Get(user)
	if err == nil {
		return nil
	} else if !b {
		return errors.New("查询失败")
	}

	fmt.Println("Query: ", user)

	return nil
}

// delete data from table
func Delete() error {
	user := &Users{}
	i, err := engine.Id(1).Delete(user)
	if err == nil {
		return nil
	} else if i <= 0 {
		return errors.New("删除失败")
	}

	return nil
}

func main() {
	// insert
	if err := Insert(); err != nil {
		fmt.Println(err)
		return
	}

	// update
	if err := Update(); err != nil {
		fmt.Println(err)
		return
	}

	// query
	if err := Query(); err != nil {
		fmt.Println(err)
		return
	}

	// delete
	if err := Delete(); err != nil {
		fmt.Println(err)
		return
	}
}
