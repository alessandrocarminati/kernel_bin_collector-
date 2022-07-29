package main
 
import (
	"database/sql"
	"fmt"
	"strings"
	_ "github.com/lib/pq"
)
type Connect_token struct{
	Host	string
	Port	int
	User	string
	Pass	string
	Dbname	string
}

func Connect_db(t *Connect_token) (*sql.DB){
	fmt.Println("connect")
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", (*t).Host, (*t).Port, (*t).User, (*t).Pass, (*t).Dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err!= nil {
		panic(err)
		}
	fmt.Println("connected")
	return db
}

func Insert_data(db *sql.DB, symbol string, fn string, xrefs []string, test bool){

	//insert container file
	fmt.Println("insert file")
	query:="insert into files (file_name) select CAST($1 AS VARCHAR) where not exists (select * from files where file_name=$1);"
	if !test {
		_, err := db.Exec(query, fn)
		if err!= nil {
			panic(err)
			}
		} else {
			query=strings.ReplaceAll(query, "$1", fn)
			fmt.Println(query)
			}

	//insert main symbol, update if already present (inserted by xref)
	fmt.Println("insert main symb")
	query="update symbols set file_ref_id=(select file_id from files where file_name=$2) where symbol_name=$1 and file_ref_id is null;";
	if !test {
		_, err := db.Exec(query, symbol, fn)
		if err!= nil {
			panic(err)
			}
		} else {
			query=strings.ReplaceAll(query, "$1", fn)
			query=strings.ReplaceAll(query, "$2", symbol)
			fmt.Println(query)
			}


	query="insert into symbols (symbol_name, file_ref_id) select CAST($1 AS VARCHAR), (select file_id from files where file_name=$2) where not exists (select * from symbols where symbol_name=$1 and file_ref_id=(select file_id from files where file_name=$2));";
	if !test {
		_, err := db.Exec(query, symbol, fn)
		if err!= nil {
			panic(err)
			}
		} else {
			query=strings.ReplaceAll(query, "$1", fn)
			query=strings.ReplaceAll(query, "$2", symbol)
			fmt.Println(query)
			}

	//insert symbols from xrefs. they are incomplete and they will eventually update to final state when their definition is found
	//in the same cycle, fetch all references, and populate the xrefs table
	fmt.Println("insert xrefs")
	for _,s:= range xrefs {
		query="insert into symbols (symbol_name) select CAST($1 AS VARCHAR) where not exists (select * from symbols where symbol_name=$1);"
		if !test {
			_, err := db.Exec(query, s)
			if err!= nil {
				panic(err)
				}
			} else {
				query=strings.ReplaceAll(query, "$1", s)
				fmt.Println(query)
				}
		query="insert into xrefs (caller, callee) select (select symbol_id from symbols where symbol_name=$1), (select symbol_id from symbols where symbol_name=$2);"
		if !test {
			_, err := db.Exec(query, symbol, s)
			if err!= nil {
				panic(err)
				}
			} else {
				query=strings.ReplaceAll(query, "$1", symbol)
				query=strings.ReplaceAll(query, "$2", s)
				fmt.Println(query)
				}

		}

}


/*
create table files (file_id SERIAL PRIMARY KEY, file_name varchar(100) UNIQUE);
create table symbols (symbol_id SERIAL PRIMARY KEY, symbol_name varchar(100) unique, file_ref_id int);
create table xrefs (caller int, callee int); 


insert into files (file_name) Select 'pippo.com' Where not exists (select * from files where file_name='pippo.com');
insert into symbols (symbol_name, file_ref_id) select 'peppe', (select file_id from files where file_name='pippo.com');





update symbols set file_ref_id=(select file_id from files where file_name='pippo.com') where symbol_name='peppe' and file_ref_id=null;
insert into symbols (symbol_name, file_ref_id) select 'pippo', (select file_id from files where file_name='pippo.com') where not exists (select * from symbols where symbol_name='pippo' and file_ref_id=(select file_id from files where file_name='pippo.com'));



insert into symbols (
        symbol_name, file_ref_id
        ) select 'pippo', (
        select file_id from files where file_name='pippo.com'
        ) where not exists (
                select * from symbols where symbol_name='pippo' and file_ref_id=(
                        select file_id from files where file_name='pippo.com'
                        )
                );



*/
