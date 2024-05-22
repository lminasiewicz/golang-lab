package main

import (
	"fmt";
	"net/http";
	"os";
	"io";
	"encoding/json";
	"sync";
	"math/rand"
)


var mutex sync.Mutex


type Entry struct {
	Id int
	Data Data
}

type Data struct {
    Year     string `json:"Year"`
    Type     string `json:"Type"`
    Country  string `json:"Country"`
    Activity string `json:"Activity"`
	Sex      string `json:"Sex"`
	Fatal    string `json:"Fatal_y_n"`
}


func get_db() []Entry {
	jsonFile, err := os.Open("global-shark-attack.json")
    if err != nil {
        fmt.Println(err)
    }
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

    var jsonList []Data
    json.Unmarshal(byteValue, &jsonList)

	var db []Entry
	for i := 0; i < len(jsonList); i++ {
		var newEntry Entry
		newEntry.Id = i
		newEntry.Data = jsonList[i]
		db = append(db, newEntry)
	}
	return db
}

func choose_random_ten(db []Entry) []Entry {
	len_db := len(db)
	random_post_ids := [10]int{}
	for i := 0; i < len(random_post_ids); i++ {
		for {
			unique := true
			index := rand.Intn(len_db-1)
			for _, elem := range random_post_ids {
				if index == elem {
					unique = false
					break
				}
			}
			if unique == true {
				random_post_ids[i] = index
				break
			}
		}
	}

	random_posts := make([]Entry, 10)
	for i := 0; i < len(random_posts); i++ {
		random_posts[i] = db[random_post_ids[i]]
	}
	return random_posts
}


func entriesHandler (w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		mutex.Lock()
		defer mutex.Unlock()

		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(random_posts)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}


func main() {
	db := choose_random_ten(get_db())

	http.HandleFunc("/entries", entriesHandler)

	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}