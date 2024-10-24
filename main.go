package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"net/http"
	"sort"
	"strings"
)

type Country struct {
	Name       string `json:"name"`
	Population int    `json:"population"`
	Details    string `json:"details"`
	Images     string `json:"images"`
}

var countries []Country

func scrapeCountries() {
	log.Println("Iniciando o scraping dos países...")
	c := colly.NewCollector()

	c.OnHTML("table.wikitable.sortable tbody tr", func(e *colly.HTMLElement) {
		if e.ChildText("th") != "" {
			return
		}

		country := Country{}
		e.ForEach("td", func(i int, element *colly.HTMLElement) {
			switch i {
			case 1:
				//Pula o index

			case 2:
				country.Name = strings.TrimSpace(element.ChildText("a"))
				country.Details = strings.TrimSpace(element.ChildAttr("a", "href"))
				country.Images = strings.TrimSpace(element.ChildAttr("img", "src"))
			case 3:

				fmt.Sscanf(strings.ReplaceAll(element.Text, " ", ""), "%d", &country.Population)
			}
		})

		if country.Name != "" {
			countries = append(countries, country)
		}
	})

	err := c.Visit("https://pt.wikipedia.org/wiki/Lista_de_pa%C3%ADses_por_popula%C3%A7%C3%A3o")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Total de países coletados: %d", len(countries))
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func getCountries(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	sort.Slice(countries, func(i, j int) bool {
		return countries[i].Name < countries[j].Name
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(countries)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	w.Write([]byte("API para scraping de países. Acesse /api/countries para obter os dados."))
}

func main() {
	scrapeCountries()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/api/countries", getCountries)

	fmt.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
