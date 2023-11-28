package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
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
	Address struct {
		City string
		State string
		Country string
		Pincode json.Number
	}
	User struct {
		Name string
		Age json.Number
		Contact string
		Company string
		Address Address
	}
	Options struct {
		Logger
	}
)

func stat(path string) (fi os.FileInfo,err error) {
	if fi,err= os.Stat(path);os.IsNotExist(err){
		fi,err= os.Stat(path+".json")
	}
	return
}


func New(dir string, options *Options)(*Driver , error){
		dir= filepath.Clean(dir)
		opts:=Options{}
		if options!=nil{
			opts=*options
		}
		if opts.Logger==nil{
			opts.Logger=lumber.NewConsoleLogger((lumber.INFO))
		}
		driver:=Driver{
			dir: dir,
			mutexes: make(map[string]*sync.Mutex),
			log:opts.Logger,
		}
			if _,err:=os.Stat(dir);err==nil{
				opts.Logger.Debug("Using '%s' (database already exits)\n",dir)
				return &driver,nil
			}
			opts.Logger.Debug("Creating the database at '%s'...\n",dir)
			return &driver, os.MkdirAll(dir,0755)
}

// if mutex exists get it or then create th new mutexes
func (d *Driver) getOrCreateMutex(collection string)(*sync.Mutex){
	d.mutex.Lock()
	defer d.mutex.Unlock()
	m,ok := d.mutexes[collection]
	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}
	return m
}

func (d *Driver) Write(collection,resource string, v interface{})error{
	if collection==""{
		return fmt.Errorf("missing collection - no place to save record")
	}
	if resource==""{
		return fmt.Errorf("missing resource - unable to save record(no name)")
	}
	mutex:= d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir:=filepath.Join(d.dir,collection)
	finalpath:=filepath.Join(dir,resource+".json")
	tmpPath:=finalpath+".tmp"
	if err:=os.MkdirAll(dir,0755);err!=nil{
		return err
	}
	b,err:=json.MarshalIndent(v,"","\t")
	if err!=nil{
		return err
	}
	b=append(b, byte('\n'))
	if err:=os.WriteFile(tmpPath,b,0644);err!=nil{
		return err
	}
	return os.Rename(tmpPath,finalpath)
}



func (d *Driver) Read(collection,resource string,v interface{})error{
	if collection == ""{
		return fmt.Errorf("missing collection - unable to read")
	}
	if resource ==""{
		return fmt.Errorf("missing resource - unable to read record (no name)")
	}
	record := filepath.Join(d.dir,collection,resource)
	if _,err:= stat(record);err!=nil{
		return err
	}
	b,err:= os.ReadFile(record+".json")
	if err!=nil{
		return err
	}
	return json.Unmarshal(b,&v)
}


func (d *Driver) ReadAll(collection string) ([]string,error) {
	if collection==""{
		return nil, fmt.Errorf("missing collection - unable to read")
	}
	dir := filepath.Join(d.dir , collection)

	if _,err:= stat(dir);err!=nil{
		return nil,err
	}
	files,_:= os.ReadDir(dir)
	var records []string
 	for _ , file := range files{
		b, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err!=nil{
			return nil,err
		}
		records= append(records, string(b))
	}

	return records , nil
}

func (d *Driver) Delete(collection,resource string)error{
	path:=filepath.Join(collection,resource)
	mutex:=d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir:=filepath.Join(d.dir,path)
	switch fi,err:= stat(dir);{
		case fi==nil,err!=nil:
			return fmt.Errorf("unable to find file or directory named %v here " , path)
		case fi.Mode().IsDir():
			return os.RemoveAll(dir)
		case fi.Mode().IsRegular():
			return os.RemoveAll(dir+".json")	
	}
return nil
}

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

	for _, value := range employees{
		db.Write("users", value.Name, User{
			Name:value.Name,
			Age: value.Age,
			Contact: value.Contact,
			Company: value.Company,
		})
	} 
	records,err := db.ReadAll("users")
	if err!= nil {
		fmt.Println("Error",err)
	}
	fmt.Println(records)

	allusers := []User{}

	for _,f:= range records{
		employeeFound:=User{}
		if err := json.Unmarshal([]byte(f),&employeeFound);err!=nil{
			fmt.Println("Error",err)
		}
		allusers=append(allusers, employeeFound)
	}
	fmt.Println(allusers)
	// if err:= db.Delete("users","bdbbd");err!=nil{
	// 	fmt.Println("Error",err)
	// }

	// if err:= db.Delete("user","");err!=nil{
	// 	fmt.Println("Error",err)
	// }
}