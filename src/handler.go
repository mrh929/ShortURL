package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func shortenHandler(w http.ResponseWriter, r *http.Request) { // handle shorten url requests
	// fmt.Println("ShortenHandler called!")

	if err := r.ParseForm(); err != nil {
		return
	}
	url := r.PostFormValue("url")
	passwd := r.PostFormValue("passwd")

	m := make(map[string]string)
	if passwd == "" || url == "" { // if an empty request, just redirect
		http.Error(w, "url or passwd cannot be empty", http.StatusBadRequest)
		return
	} else if passwd != SRV_PASSWD { // validate password
		m["status"] = "1"
		m["reason"] = "Invalid Password"
		// fmt.Fprintf(w, "Invalid Password")
		json, _ := json.Marshal(m)
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		log.Info("wrong password")
		return
	} else {
		r_url, err := parseRawURL(url) // convert url to standard url
		if err != nil {
			m["status"] = "2"
			m["reason"] = "Invalid URL"
			// fmt.Fprintf(w, "Invalid URL")
			json, _ := json.Marshal(m)
			w.Header().Set("Content-Type", "application/json")
			w.Write(json)
			log.Info("parseRawURL error=======", err.Error())
			return
		}

		s_key, needInsert, err := keyGenerate(r_url)
		if err != nil {
			m["status"] = "3"
			m["reason"] = "Unknown error"
			json, _ := json.Marshal(m)
			w.Header().Set("Content-Type", "application/json")
			w.Write(json)
			// w.WriteHeader(http.StatusInternalServerError)
			log.Info("keyGenerate error=======", err.Error())
			return
		}

		if needInsert { // need to insert key and url
			err := urlInsert(s_key, r_url)
			if err != nil {
				m["status"] = "3"
				m["reason"] = "Unknown error"
				json, _ := json.Marshal(m)
				w.Header().Set("Content-Type", "application/json")
				w.Write(json)
				// w.WriteHeader(http.StatusInternalServerError)
				log.Info("urlInsert error=======", err.Error())
				return
			}
		}

		var protocol string
		var base_path string

		_, ok := r.Header["X-Forwarded-Proto"]
		if ok { // decide http or https
			protocol = r.Header["X-Forwarded-Proto"][0]
		} else {
			protocol = SRV_PROTO
		}

		if SRV_BASE_PATH != "" { // decide base path
			base_path = SRV_BASE_PATH
		} else {
			base_path = r.Host
		}

		m["shortenURL"] = fmt.Sprintf("%s://%s/%s", protocol, base_path, s_key)
		m["realURL"] = r_url
		m["status"] = "0"
		m["reason"] = "success"
		json, _ := json.Marshal(m)
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)

		log.Info(s_key, "----->", r_url)
		// fmt.Fprintf(w, s_url)
		return
	}

}

func urlHandler(w http.ResponseWriter, r *http.Request) { // handle url requets
	// fmt.Println("URLHandler called!")

	vars := mux.Vars(r)

	if url, err := urlSelect(vars["key"]); url != "" {
		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	} else if err == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		log.Info("urlSelect error=======", err.Error())
		return
	}
}
