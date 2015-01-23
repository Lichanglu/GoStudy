package main

import (
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	//	"reflect"
	"time"
)

type Account struct {
	Id      int64
	Name    string `xorm:"unique"`
	Balance float64
	Version int `xorm:"version"`
}

var x *xorm.Engine

func init() {
	var err error
	x, err = xorm.NewEngine("sqlite3", "bank.db")
	if nil != err {
		log.Fatalf("Create engine failed:%v\n", err)
	}
	err = x.Sync(new(Account))
	if nil != err {
		log.Fatalf("Sync DataBase failed:%v\n", err)
	}
}
func newAccount(name string, balance float64) error {
	if "" == name || balance <= 0 {
		return errors.New("Failed to input")
	}
	_, err := x.Insert(&Account{Name: name, Balance: balance})
	if nil != err {
		log.Println(err.Error())
	} else {
		log.Println("Insert success\n")
	}
	return err
}
func ShowAccount(Id int64) (*Account, error) {
	//a := &Account{}
	defer func() {
		if x := recover(); x != nil {
			fmt.Println(x)
			for {
				time.Sleep(10 * time.Second)
			}
		}
	}()
	if Id <= 0 {
		return nil, errors.New("Failed to input")
	}
	a := &Account{}
	has, err := x.Id(Id).Get(a)
	if nil != err {
		panic(has)
		fmt.Println(err.Error())
		return nil, err
	}
	if !has {
		return nil, errors.New("Account not found")
	}
	return a, err
}
func DeleteAccount(id int64) error {
	/*defer func() {
		if x := recover(); x != nil {
			fmt.Println(x)
			for {
				time.Sleep(10 * time.Second)
			}
		}
	}()*/
	if id <= 0 {
		return errors.New("Failed to input")
	}
	num, err := x.Delete(&Account{Id: id})
	if nil != err {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("num=", num)
	return err
}
func DeleteNameAccount(name string) error {
	/*defer func() {
		if x := recover(); x != nil {
			fmt.Println(x)
			for {
				time.Sleep(10 * time.Second)
			}
		}
	}()*/
	if "" == name {
		return errors.New("Failed to input")
	}
	num, err := x.Delete(&Account{Name: name})
	if nil != err {
		fmt.Println(err.Error())
	}
	fmt.Println("num=", num)
	return err
}

func ShowAllAccount() ([]Account, error) {
	a := make([]Account, 10)
	err := x.Asc("balance").Find(&a)
	for _, value := range a {
		if value.Version > 0 {
			fmt.Printf("--->:%#v\n", value)
		}
	}
	fmt.Printf("\n\n")
	return nil, err
}
func DepositAccount(deposit float64, id int64) (*Account, error) {
	if deposit <= 0 || id <= 0 {
		return nil, errors.New("Failed to input")
	}

	a, err := ShowAccount(id)
	if nil != err {
		fmt.Println(err.Error())
		return nil, err
	}
	a.Balance += deposit
	_, err = x.Update(a)
	return a, err
}
func WithdrawAccount(Withdraw float64, id int64) (*Account, error) {
	if Withdraw <= 0 || id <= 0 {
		return nil, errors.New("Failed to input")
	}
	a, err := ShowAccount(id)
	if nil != err {
		fmt.Println(err.Error())
		return nil, err
	}
	if a.Balance <= Withdraw {
		err := errors.New("Account balance not enough!!")
		return a, err
	}
	a.Balance -= Withdraw
	_, err = x.Update(a)
	if nil != err {
		fmt.Println(err.Error())
		return nil, err
	}
	return a, err
}

func MakeTransferAccount(src_id int64, dst_id int64, funds float64) error {
	if src_id <= 0 || dst_id <= 0 || funds <= 0 {
		return errors.New("Failed to input")
	}
	src_a, err := ShowAccount(src_id)
	if nil != err {
		fmt.Println("1", err.Error())
		return err
	}
	if src_a.Balance < funds {
		err := errors.New("Account balance not enough!!")
		return err
	}
	dst_a, err := ShowAccount(dst_id)
	if nil != err {
		fmt.Println("2", err.Error())
		return err
	}
	src_a.Balance -= funds
	dst_a.Balance += funds
	//创建回滚事务
	secc := x.NewSession()
	defer secc.Close()
	//启动事务
	if err = secc.Begin(); nil != err {
		return err
	}

	_, err = secc.Update(src_a)
	if nil != err {
		//发生错误，进行回滚
		secc.Rollback()
		fmt.Println("3", err.Error())
		return err
	}
	_, err = secc.Update(dst_a)
	if nil != err {
		secc.Rollback()
		fmt.Println("4", err.Error())
		return err
	}
	//提交事务
	return secc.Commit()
}
func FindName(f xorm.IterFunc) error {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println(x)
			for {
				time.Sleep(10 * time.Second)
			}
		}
	}()

	//	v := reflect.ValueOf(f)
	//	if xorm.IterFunc != v.Type() {
	//		return errors.New("Failed to arg ")
	//	}
	x.Cols("name").Iterate(new(Account), f)
	return nil
}
