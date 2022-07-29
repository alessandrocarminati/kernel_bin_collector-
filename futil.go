package main
import (
	"fmt"
	"io/ioutil"
	"regexp"
	"database/sql"
)

func apply2file(path string, fun func(string, *sql.DB),  t *Connect_token){


	db:=Connect_db(t)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
		}

	for _, f := range files {
		if f.IsDir() {
			fmt.Println("directory ",path+"/"+f.Name())
			apply2file(path+"/"+f.Name(), fun, t)
		} else {
			fmt.Println("file ",path+"/"+f.Name())
			match, _ := regexp.MatchString("\\.o$", path+"/"+f.Name())
			if match {
				fmt.Println("match ",path+"/"+f.Name())
				fun(path+"/"+f.Name(), db)
				}
			}
		}
}
