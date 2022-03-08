package main

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) { // index page
	// fmt.Println("HomeHandler called!")
	path := path.Dir("./web/index.html")
	http.ServeFile(w, r, path)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) { // handle shorten url requests
	// fmt.Println("ShortenHandler called!")

	if err := r.ParseForm(); err != nil {
		return
	}
	url := r.PostFormValue("url")
	passwd := r.PostFormValue("passwd")

	m := make(map[string]string)
	if passwd == "" && url == "" { // if an empty request, just redirect
		http.Redirect(w, r, "http://"+r.Host, http.StatusTemporaryRedirect)
		return
	} else if passwd != SRV_PASSWD { // validate password
		m["status"] = "failed"
		m["url"] = ""
		m["reason"] = "Invalid Password"
		// fmt.Fprintf(w, "Invalid Password")
		json, _ := json.Marshal(m)
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	} else {
		r_url, err := parseRawURL(url) // convert url to standard url
		if err != nil {
			m["status"] = "failed"
			m["url"] = ""
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
			m["status"] = "failed"
			m["url"] = ""
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
				m["status"] = "failed"
				m["url"] = ""
				m["reason"] = "Unknown error"
				json, _ := json.Marshal(m)
				w.Header().Set("Content-Type", "application/json")
				w.Write(json)
				// w.WriteHeader(http.StatusInternalServerError)
				log.Info("urlInsert error=======", err.Error())
				return
			}
		}

		m["status"] = "success"
		m["url"] = r.Host + "/" + s_key
		m["reason"] = "success"
		json, _ := json.Marshal(m)
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
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
