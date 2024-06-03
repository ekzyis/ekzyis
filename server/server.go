package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	sn "github.com/ekzyis/snappy"
)

func main() {
	bindAddr := "127.0.0.1:1337"
	http.HandleFunc("/api/content_meta", HandleContentMeta)

	log.Printf("server running on %s\n", bindAddr)
	http.ListenAndServe(bindAddr, logRequest(http.DefaultServeMux))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func HandleContentMeta(w http.ResponseWriter, r *http.Request) {
	var (
		// stacker.news item id
		qid  = r.URL.Query().Get("sn_id")
		c    = sn.NewClient()
		item *sn.Item
		id   int
		err  error
	)

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	w.Header().Add("Access-Control-Allow-Headers", "*")

	if r.Method == "OPTIONS" {
		return
	}

	if qid == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "sn_id query param required\n")
		return
	}

	if id, err = strconv.Atoi(qid); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "sn_id must be numeric\n")
		return
	}

	if item, err = c.Item(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	fmt.Fprintf(w,
		""+
			"<span class=\"px-1\">|</span>"+
			"<span>"+
			"<a"+
			" class=\"underline\""+
			" href=\"https://stacker.news/items/%d\""+
			" target=\"_blank\""+
			" rel=\"noopener noreferrer me\">"+
			"%d comments"+
			"</a>"+
			"</span>"+
			"<span class=\"px-1\">|</span>"+
			"<span>%d sats</span>",
		id, item.NComments, item.Sats)
}
