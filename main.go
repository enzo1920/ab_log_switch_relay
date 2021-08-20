package main

import (
    "fmt"
    "log"
    "time"
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



// структура для реле, заполняем состояниями
type LightRelays struct{
    R_id int
    R_ip string
    R_state int
}


func main() {
        relaySwitch()
        //relaySwitchBytime()
}



func relaySwitch(){
        var light_effect float64 = 50.0
        getsum:= getLightsum()
        getstate:= getRelaystate()
        for _, ls := range getstate {
                if getsum <= light_effect && ls.R_state==0 {
                         sendGet("http://admin:admin@"+ls.R_ip+"/protect/rb0n.cgi")
                         fmt.Println("Relay ON")
                }else if getsum>=light_effect && ls.R_state==1 {
                         sendGet("http://admin:admin@"+ls.R_ip+"/protect/rb0f.cgi") 
                         fmt.Println("Relay OFF")
                }

        }

        fmt.Println(getsum)
}



/*
func relaySwitch(){
        var light_effect float64 = 50.0
        getsum:= getLightsum()
        getstate:= getRelaystate()
        //если из базы вернулся ноль, то ничего не делаем
        if getsum == 0{
           fmt.Println("return sum = 0. Wrong sum query.Check db and inserter service!!! ")
        }else{
               for _, ls := range getstate {
                   if getsum <= light_effect && ls.R_state==0 {
                         sendGet("http://admin:admin@"+ls.R_ip+"/protect/rb0n.cgi")
                         fmt.Println("Relay ON")
                   }else if getsum>=light_effect && ls.R_state==1 {
                         sendGet("http://admin:admin@"+ls.R_ip+"/protect/rb0f.cgi") 
                         fmt.Println("Relay OFF")
                   }

               }
        }

        fmt.Println(getsum)
}
*/

func relaySwitchBytime(){
        //var light_effect float64 = 50.0
        //getsum:= getLightsum()
        night := 22
        day := 4
        now :=time.Now()
        hour :=now.Hour()
        getstate:= getRelaystate()
        for _, ls := range getstate {
                if hour >= night && ls.R_state==0 {
                         sendGet("http://admin:admin@"+ls.R_ip+"/protect/rb0n.cgi")
                         fmt.Println("Relay ON")
                   }else if hour >= day && hour < night && ls.R_state==1 {
                         sendGet("http://admin:admin@"+ls.R_ip+"/protect/rb0f.cgi") 
                         fmt.Println("Relay OFF")
                   }

               }
        fmt.Println(now.Hour(),day, night)
}




//получаем сумму значений датчика за последнии 15 мин
func getLightsum() float64 {
        dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
            DB_USER, DB_PASSWORD, DB_NAME)
        db, err := sql.Open("postgres", dbinfo)
        checkErr(err)
        defer db.Close()

	row := db.QueryRow("SELECT COALESCE(sum(light_val),0) FROM light WHERE light_date BETWEEN NOW()- INTERVAL '15 minutes' and NOW()")
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




//получаем состояние реле
func getRelaystate() []LightRelays {
        dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
            DB_USER, DB_PASSWORD, DB_NAME)
        db, err := sql.Open("postgres", dbinfo)
        checkErr(err)
        defer db.Close()



        rows, err := db.Query("SELECT r_id, r_ip, r_state FROM relays WHERE r_type=1  ORDER BY r_id asc")
        if err != nil {
             log.Fatal(err)
        }
        defer rows.Close()

        lightrelays := []LightRelays{}

        for rows.Next(){
             lr := LightRelays{}
             err := rows.Scan(&lr.R_id, &lr.R_ip, &lr.R_state)
             if err != nil{
                  fmt.Println(err)
                  continue
             }
             lightrelays = append(lightrelays, lr)
        }


    return lightrelays

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

