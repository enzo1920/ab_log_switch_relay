package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "database/sql"
   _ "github.com/lib/pq"
)

const (
     DB_USER     = "postgres"
     DB_PASSWORD = ""
     DB_NAME     = "ab_log_db"
)


func main() {
        var light_effect float64 = 50.0
        //fmt.Println(light)
        getsum:=getLightsum()
        fmt.Println(getsum)
        if getsum <= light_effect  {
            sendGet("http://admin:admin@192.168.71.117/protect/rb0n.cgi")
            fmt.Println("Relay ON")
        } else {
            sendGet("http://admin:admin@192.168.71.117/protect/rb0f.cgi") 
            fmt.Println("Relay OFF")
        }
}

//получаем сумму значений датчика за последнии 15 мин
func getLightsum() float64 {
        dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
            DB_USER, DB_PASSWORD, DB_NAME)
        db, err := sql.Open("postgres", dbinfo)
        checkErr(err)
        defer db.Close()

	row := db.QueryRow("SELECT COALESCE(sum(light_val),0) light_date FROM light WHERE light_date BETWEEN NOW()- INTERVAL '15 minutes' and NOW()")
	if err != nil {
		log.Fatal(err)
	}
	//defer row.Close()
	var lightsum float64
	if err := row.Scan(&lightsum); err != nil {
	     // Check for a scan error.
	     // Query rows will be closed with defer.
		log.Fatal(err)
       }

	return lightsum
}

//отправка запросов
func sendGet(host string){
    resp, err := http.Get(host) 
    if err != nil { 
        fmt.Println(err) 
        return
    } 
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
          fmt.Println(err)
          return
    }
    //fmt.Println(string(body))
    fmt.Println(string(body))
}





func checkErr(err error) {
        if err != nil {
            panic(err)
        }
    }
