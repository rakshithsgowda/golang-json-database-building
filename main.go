package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

type (
	Logger interface{
		Fatal(string,...interface{})
		Error(string,...interface{})
		Warn(string,...interface{})
		Info(string,...interface{})
		Debug(string,...interface{})
		Trace(string,...interface{})
	}
	Driver struct{
		mutex sync.Mutex
		mutexes map[string]*sync.Mutex
		dir string
		log Logger
	}
)

type Address struct {
	City string
	State string
	Country string
	Pincode json.Number
}


type User struct{
	Name string
	Age json.Number
	Contact string
	Company string
	Address Address
}

func New(dir,options)(*Driver,error){}
func (d *Driver)Write(){}
func (d *Driver)Read(){}
func (d *Driver)ReadAll(){}
func (d *Driver)Delete(){}

func main(){
	// define the directory for the folder creation.
	dir := "./"

	db,err:= New(dir,nil)
	if err!=nil{
		fmt.Println("Error",err)
	}

	employees:=[]User{
		{"Jsohn","23","2954419","greatcompany",Address{"bangalore","karnataka","india","560083"}},
		{"ccwvd","25","5445466","greddcompany",Address{"bangalore","karnataka","india","560013"}},
		{"bdbbd","26","4552455","gfrfbcompany",Address{"bangalore","karnataka","india","560056"}},
		{"ghkgv","30","7858754","gfvfvcompany",Address{"bangalore","karnataka","india","560093"}},
		{"wrfgv","33","8988966","grvftcompany",Address{"bangalore","karnataka","india","560043"}},
	}

	for _, value:= range employees{
		db.Write("users",value.Name,User{
			Name:value.Name,
			Age: value.Age,
			Contact: value.Contact,
			Company: value.Company,
		})
	} 
	records,err := db.ReadAll("users")
	if err!=nil{
		fmt.Println("Error",err)
	}
	fmt.Println(records)

	allusers:=[]User{}

	for _,f:= range records{
		employeeFound:=User{}
		if err:=json.Unmarshal([]byte(f),&employeeFound);err!=nil{
			fmt.Println("Error",err)
		}
		allusers=append(allusers, employeeFound)
	}
	fmt.Println(allusers)
	// if err:= db.Delete("users","");err!=nil{
	// 	fmt.Println("Error",err)
	// }
}