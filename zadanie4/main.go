package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

type Table struct {
	Type          string
	Mutability    string
	Description   string
	SyntaxExample string
}

func scrape(url string) ([]Table, error) {
	var table_data []Table
	skip_labels := false

	c := colly.NewCollector()

	c.OnHTML(".wikitable > tbody", func(h *colly.HTMLElement) {
		h.ForEach("tr", func(_ int, el *colly.HTMLElement) {

			if !skip_labels {
				skip_labels = true
				return
			}

			data := Table{
				Type:          el.ChildText("td:nth-child(1)"),
				Mutability:    el.ChildText("td:nth-child(2)"),
				Description:   el.ChildText("td:nth-child(3)"),
				SyntaxExample: el.ChildText("td:nth-child(4)"),
			}
			table_data = append(table_data, data)
		})
	})

	err := c.Visit(url)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return nil, err
	}

	return table_data, nil
}

func write_to_csv(data []Table, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	writer := csv.NewWriter(file)
	defer file.Close()
	defer writer.Flush()

	headers := []string{"Type", "Mutability", "Description", "SyntaxExample"}
	writer.Write(headers)
	for _, in_row := range data {
		out_row := []string{in_row.Type, in_row.Mutability, in_row.Description, in_row.SyntaxExample}
		writer.Write(out_row)
	}
}

func main() {
	url := "https://en.wikipedia.org/wiki/Python_(programming_language)"
	out_file := "python_types.csv"

	table_data, err := scrape(url)
	if err != nil {
		return
	}

	write_to_csv(table_data, out_file)
	fmt.Println("Done!")
}
