package main

import (
	"fmt"
	"reflect"
	//	"time"
	//	"os"
	"errors"
)

const usrgrad = `Please enter number of operation:
1.Create new account
2.Show datail of account
3.Deposit
4.Withdraw
5.Make transfer
6.List account by Id
7.List account by balance
8.Delete account
9.Exit`

var printFunc = func(idx int, been interface{}) error {
	fmt.Printf("idx=%d :%#v\n", idx, been)
	return nil
}

type user struct {
	id   int64
	Name string
	age  int
}

func (u user) Hello(to string, time int) error {
	fmt.Println("Hello", to, ",My name is", time, u.Name)
	return errors.New("Print")
}
func Info(o interface{}) {
	v := reflect.ValueOf(o)
	if v.Kind() != reflect.Ptr {
		fmt.Println("o Not Ptr")
		return
	}

	if !v.Elem().CanSet() {
		fmt.Println("XX")
		return
	} else {
		fmt.Println("1")
		v = v.Elem()
	}
	f := v.FieldByName("Name")
	if !f.IsValid() {
		fmt.Println("Byebye")
		return
	}
	if f.Kind() == reflect.String {
		f.SetString("Lucy")
	}
	ff := v.MethodByName("Hello")
	arg1 := []reflect.Value{reflect.ValueOf("Lily"), reflect.ValueOf(32)}
	err := ff.Call(arg1)
	if nil != err {
		fmt.Println(err)
		return
	}
	return
}
func setReflect() {
	x := 123
	v := reflect.ValueOf(&x)
	//只有指针才有Elem
	if v.Kind() == reflect.Ptr || !v.Elem().CanSet() {
		v.Elem().SetInt(1000)
		fmt.Println(x)
	} else {
		fmt.Println("v not can set")
	}
	u := user{100, "Joe", 17}
	Info(u)
	fmt.Println(u)

	//只有指针才有Elem
	xx := 22
	vv := reflect.ValueOf(xx)
	fmt.Println("-->", vv.Kind())
}

func main() {
	setReflect()
	fmt.Println("Welcome Bank of xorm By Lu!")

Exit:
	for {
		fmt.Println(usrgrad)
		var num int
		fmt.Scanf("%d\n", &num)
		switch num {
		case 1:
			fmt.Println("Please input <name><Balance>")
			var err error
			var name string
			var balance float64
			fmt.Scanf("%s %f\n", &name, &balance)
			if " " != name || balance > 0.0 {
				err = newAccount(name, balance)
			} else {
				fmt.Printf("Warning,%s,%f\n", name, balance)
			}
			if nil != err {
				fmt.Println(err.Error())
			}
		case 2:
			fmt.Println("Please input <ID>")
			var id int64
			fmt.Scanf("%d\n", &id)
			a, err := ShowAccount(id)
			if nil != err {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("%#v\n\n\n", a)
			}
			break

		case 3:
			fmt.Println("Please Input <ID><Deposit>")
			var id int64
			var deposit float64
			fmt.Scanf("%d %f\n", &id, &deposit)
			_, err := DepositAccount(deposit, id)
			if nil != err {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Deposit Success!\n\n\n")
			}
		case 4:
			fmt.Println("Please Input <ID><Withdraw>")
			var id int64
			var withdraw float64
			fmt.Scanf("%d %f\n", &id, &withdraw)
			_, err := WithdrawAccount(withdraw, id)
			if nil != err {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Withdraw Success!\n\n\n")
			}
		case 5:
			fmt.Println("Please Input <Self-ID><Other-ID><Funds>")
			var self_id, other_id int64
			var funds float64
			fmt.Scanf("%d %d %f\n", &self_id, &other_id, &funds)
			err := MakeTransferAccount(self_id, other_id, funds)
			if nil != err {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Make transfer Success!\n\n\n")
			}

		//List account by Id
		case 6:
		//
		case 7:
			err := FindName(printFunc)
			if nil != err {
				fmt.Println(err.Error())
			}
		case 8:
			fmt.Println("Please input <ID>")
			var name string
			fmt.Scanf("%s\n", &name)
			err := DeleteNameAccount(name)
			if nil != err {
				fmt.Println(err.Error())
			}
		case 9:
			fmt.Println("Bank of xorm system exit")
			break Exit
		case 10:
			_, _ = ShowAllAccount()
		default:
			fmt.Println("Input num Error!!Re-input")
		}
	}
}
