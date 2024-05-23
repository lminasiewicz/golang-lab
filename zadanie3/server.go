package main

import (
	"fmt";
	"net/http";
	"os";
	"io";
	"encoding/json";
	"sync";
	"math/rand";
	"strconv"
)


var mutex sync.Mutex
var db []Entry = choose_random_ten(get_db())


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


func remove_from_db(idx int) bool {
	if idx == -1 {
		return false
	}
    db[idx] = db[len(db)-1]
    db = db[:len(db)-1]
	return true
}


func find_entry(id int) int {
	for index, entry := range db {
		if entry.Id == id {
			return index
		}
	}
	return -1
}

func get_next_id() int {
	current := -1
	for _, entry := range db {
		if entry.Id > current {
			current = entry.Id
		}
	}
	return current + 1
}


func post_entry_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var data Data
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, "Error parsing request body", http.StatusBadRequest)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		id := get_next_id()
		var entry Entry
		entry.Id = id
		entry.Data = data

		db = append(db, entry)

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(entry)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}


func entries_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		mutex.Lock()
		defer mutex.Unlock()

		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(db)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}


func entry_handler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/entries/"):])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		mutex.Lock()
		defer mutex.Unlock()

		var idx int = find_entry(id)
		if idx == -1 {
			http.Error(w, "Post not found", http.StatusBadRequest)
		} else {
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(db[idx])
		}

	case "DELETE":
		mutex.Lock()
		defer mutex.Unlock()

		idx := find_entry(id)
		if remove_from_db(idx) == false {
			http.Error(w, "Post not found", http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}


func main() {
	http.HandleFunc("/entries", entries_handler)
	http.HandleFunc("/entries/", entry_handler)
	http.HandleFunc("/entries/submit", post_entry_handler)

	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}