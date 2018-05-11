package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"fmt"
)

import (
	"log"
	"encoding/json"
)
type bpHit struct {
	bp string
	color string
}
func bpHitRatio(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var msg string
	var color string
	var b = bpHit{}
	err := db.QueryRow(`SELECT
	min(TOTAL_HIT_RATIO_PERCENT) as minHR,
		CASE
	when min(TOTAL_HIT_RATIO_PERCENT) < 85.00 then 'red'
        when min(TOTAL_HIT_RATIO_PERCENT) > 98.00 then 'green'
        else 'yellow'
    end v_Color
	from sysibmadm.bp_hitratio bphr join table (mon_get_bufferpool(NULL,-2)) mgbp on bphr.bp_name=mgbp.bp_name
	join syscat.bufferpools b on bphr.bp_name=b.bpname
	group by 1`).Scan(&msg, &color)
	if err != nil {
		log.Println(err)
	}
	b.bp = msg
	b.color = color
	j, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	log.Println(j)
	w.Write(j)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	fmt.Fprint(w, string(j))
}

