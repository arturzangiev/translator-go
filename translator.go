package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var apiKey string

type Translations struct {
	Translations []Translation `json:"translations"`
}

type Translation struct {
	Text string `json:"text"`
}

func main() {
	path, err := os.Executable()
	checkError(err)
	dir := filepath.Dir(path)
	fmt.Println(path) // the path to the current file
	fmt.Println(dir)  // the directory of the current file

	// Open source CSV file
	source, err := os.Open(filepath.Join(dir, "source.csv"))
	checkError(err)
	defer source.Close()

	// Read source file into a variable
	lines, err := csv.NewReader(source).ReadAll()
	checkError(err)

	// Open output CSV file
	output, err := os.OpenFile(filepath.Join(dir, "output.csv"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	checkError(err)

	// Writer for output CSV file
	writer := csv.NewWriter(output)

	// Loop through lines & turn into object
	for _, line := range lines {
		productID := line[0]
		description := line[1]
		translarion := translateText(line[1])
		fmt.Println("------------- SOURCE ---------------\n" + description)
		fmt.Println("------------- TRANSLATION ---------------\n" + translarion)
		writer.Write([]string{productID, description, translarion})
	}

	writer.Flush()

}

func translateText(text string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.deepl.com/v2/translate", nil)
	checkError(err)

	// Adding params to the request
	q := req.URL.Query()
	q.Add("auth_key", apiKey)
	q.Add("text", text)
	q.Add("target_lang", "NL")
	req.URL.RawQuery = q.Encode()

	// Send request to the server
	resp, err := client.Do(req)
	checkError(err)
	defer resp.Body.Close()

	// Read response into bytes
	content, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	// Return translation result
	var translations Translations
	json.Unmarshal(content, &translations)
	return translations.Translations[0].Text

}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
