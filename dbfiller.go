package main


func main() {
	t:=Connect_token{ "dbs.hqhome163.com",5432,"alessandro","<password>","kernel_bin"}
	apply2file(".", analyze_obj, &t)
}
