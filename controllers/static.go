package controllers

import (
	"html/template"
	"net/http"
)

type Static struct {
	Template Template
}

func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}

func FAQ(tpl Template) http.HandlerFunc {
	rules := []struct {
		Rule     template.HTML
		RuleInfo template.HTML
	}{
		{
			Rule:     "Nelaužykite Lietuvos ar Vokietijos įstatymų.",
			RuleInfo: "(Serveris yra Vokietijoje)",
		},
		{
			Rule:     "Neužtvindykite, nespaminkite, neskelbkite nereikalingų įrašų.",
			RuleInfo: "",
		},
		{
			Rule:     "Reklama yra draudžiama.",
			RuleInfo: "",
		},
		{
			Rule:     "Kalbėkite tik lietuviškai visose lentose, išskyrus /int/.",
			RuleInfo: "",
		},
		{
			Rule:     "Laikykitės kiekvienos lentos temos.",
			RuleInfo: "(pvz.: nedėkite to, kas nesusiję su anime /a/ lentoje)",
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, rules)
	}
}
