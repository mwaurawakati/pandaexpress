package translate

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

func TranslateTextToPreferredLanguage(text, preferredLang string) (string, error) {
	key := "6573363319e24a9d9421db5dda13be7c"
	endpoint := "https://api.cognitive.microsofttranslator.com/"
	uri := endpoint + "/translate?api-version=3.0"
	location := "westeurope"
	u, err := url.Parse(uri)
	if err != nil {
		return text, err
	}
	q := u.Query()
	q.Add("to", preferredLang)
	u.RawQuery = q.Encode()

	// Create an anonymous struct for your request body and encode it to JSON
	body := []struct {
		Text string
	}{
		{Text: text},
	}
	b, err := json.Marshal(body)
	if err != nil {
		return text, err
	}

	// Build the HTTP POST request
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(b))
	if err != nil {
		return text, err
	}
	// Add required headers to the request
	req.Header.Add("Ocp-Apim-Subscription-Key", key)
	// location required if you're using a multi-service or regional (not global) resource.
	req.Header.Add("Ocp-Apim-Subscription-Region", location)
	req.Header.Add("Content-Type", "application/json")

	// Call the Translator API
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return text, err
	}

	// Decode the JSON response
	var result []Res
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return text, err
	}
	if len(result) == 0 {
		return text, errors.New("no results")
	}
	return result[0].Transtations[0].Text, nil
}

type Res struct {
	Transtations []Translation `json:"translations"`
}

type Translation struct {
	Text string `json:"text"`
	To   string `json:"to"`
}
