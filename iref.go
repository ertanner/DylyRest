package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"fmt"
	"log"
	"encoding/json"
)
type Data struct {
	StatTime string
	Iref string
}
type Overall struct {
	Name string `json:"name"`
	Data [] Data `json:"data"`
	}
var base [] Overall

func irefOverall(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rows, err := db.Query(`select 'Overall' StatType, to_char(STAT_TIME, 'DD-MON-YYYY') as STAT_TIME, to_char(IREF) as IREF 
									from TMWIN.DYLT_IREF_STATS
									order by StatType, stat_time
									FETCH FIRST 10 ROWS ONLY`)
	if err != nil {
		log.Println(err)
	}
	var statType string
	var statTime string
	var iref string

	d := Overall{}
	d.Name = "Overall"
	d.Data = make([]Data, 10)
	i:= 0
	for rows.Next() {
		err := rows.Scan(&statType, &statTime, &iref)
		if err != nil {
			log.Println(err)
		}
		d.Data[i].StatTime = statTime
		d.Data[i].Iref = iref
		i++
	}
	fmt.Println(d)
	fmt.Println(base)
	base = append(base, d)
	fmt.Println(base)
	b, err := json.Marshal(base)
	if err != nil {
		panic(err)
	}

	// log.Println(b)
	w.Write(b)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	fmt.Fprint(w )
}

func irefPeriod(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rows, err := db.Query(`select 'Period' as StatType, to_char(STAT_TIME, 'DD-MON-YYYY') as STAT_TIME, to_char(IREF_PERIOD) as IREF
									from TMWIN.DYLT_IREF_STATS
									order by StatType, stat_time
									FETCH FIRST 10 ROWS ONLY`)
	if err != nil {
		log.Println(err)
	}
	var statType string
	var statTime string
	var iref string

	d := Overall{}
	d.Name = "Period"
	d.Data = make([]Data, 10)
	i:= 0
	for rows.Next() {
		err := rows.Scan(&statType, &statTime, &iref)
		if err != nil {
			log.Println(err)
		}
		d.Data[i].StatTime = statTime
		d.Data[i].Iref = iref
		i++
	}
	//log.Println(d)
	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		panic(err)
	}
	//log.Println(b)
	//w.Write(b)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	fmt.Fprintf(w, string(b) )
}

func iref(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	d := base
	d = append(d, getOverall())
	d = append(d, getPeriod())
	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		panic(err)
	}
	//log.Println(d)
	log.Println(d)
	w.Write(b)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	fmt.Fprint(w)
}

func getOverall() Overall {

	rows, err := db.Query(`select StatType, STAT_TIME, IREF from (
									select 'Overall' StatType, to_char(STAT_TIME, 'DD-MON-YYYY') as STAT_TIME, to_char(IREF) as IREF 
									from TMWIN.DYLT_IREF_STATS
									order by StatType, stat_time desc
									FETCH FIRST 10 ROWS ONLY
									) order by stat_time asc`)
	if err != nil {
		log.Println(err)
	}
	var statType string
	var statTime string
	var iref string

	d := Overall{}
	d.Name = "Overall"
	d.Data = make([]Data, 10)
	i:= 0
	for rows.Next() {
		err := rows.Scan(&statType, &statTime, &iref)
		if err != nil {
			log.Println(err)
		}
		d.Data[i].StatTime = statTime
		d.Data[i].Iref = iref
		i++
	}

	return  d
}

func getPeriod() Overall {
	rows, err := db.Query(`select StatType, STAT_TIME, IREF from (
select 'Period' as StatType, to_char(STAT_TIME, 'DD-MON-YYYY') as STAT_TIME, to_char(IREF_PERIOD) as IREF
from TMWIN.DYLT_IREF_STATS
order by StatType, stat_time desc
FETCH FIRST 10 ROWS ONLY
) order by stat_time asc`)
	if err != nil {
		log.Println(err)
	}
	var statType string
	var statTime string
	var iref string

	d := Overall{}
	d.Name = "Period"
	d.Data = make([]Data, 10)
	i:= 0
	for rows.Next() {
		err := rows.Scan(&statType, &statTime, &iref)
		if err != nil {
			log.Println(err)
		}
		d.Data[i].StatTime = statTime
		d.Data[i].Iref = iref
		i++
	}

	return  d
}

func irefData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	iRef := 0.0
	iRefP := 0.0
	type ir struct{
		IRef float64 `json:"IREF_Overall"`
		IRefP float64 `json:"IREF_Period"`
	}
	err := db.QueryRow(`select IREF, IREF_PERIOD from tmwin.DYLT_IREF_STATS
									where STAT_TIME = (select  max(stat_time) from tmwin.DYLT_IREF_STATS)
									`).Scan(&iRef, &iRefP)
	if err != nil {
		log.Println(err)
	}
	i := ir{IRef: iRef, IRefP: iRefP }

	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	fmt.Fprintf(w, string(b))
}
