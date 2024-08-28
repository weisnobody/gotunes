package gotunes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var ReturnExtra bool = false

type ItunesSearchRequest struct {
	Term      string `json:"term"`
	Country   string `json:"country"`
	Media     string `json:"media"`
	Entity    string `json:"entity"`
	Attribute string `json:"attribute"`
	Limit     int32  `json:"limit"`
	Version   string `json:"version"`
	Explicit  bool   `json:"explicit"`
}

type ItunesFindRequest struct {
	ItunesId    string `json:"itunes_id"`
	AmgArtistId string `json:"amg_artist_id"`
	AmgAlbumId  string `json:"amg_album_id"`
	AmgVideoId  string `json:"amg_video_id"`
	Entity      string `json:"entity"`
	Limit       int32  `json:"limit"`
	Isbn        string `json:"isbn"`
	Upc         string `json:"upc"`
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Search URL
func SearchUrl(request ItunesSearchRequest) string {
	Url, err := url.Parse("https://itunes.apple.com")
	if err != nil {
		log.Println("There was an error parsing the iTunes API Endpoint")
	}
	Url.Path += "/search"
	parameters := url.Values{}
	addParameter := func(key string, value string) {
		if value != "" {
			parameters.Add(key, value)
		}
	}
	addParameter("term", request.Term)
	addParameter("country", request.Country)
	addParameter("media", request.Media)
	addParameter("entity", request.Entity)
	addParameter("attribute", request.Attribute)
	addParameter("limit", fmt.Sprintf("%v", request.Limit))
	if request.Explicit == false {
		parameters.Add("explicit", "no")
	}
	Url.RawQuery = parameters.Encode()
	return Url.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Find URL
func FindUrl(request ItunesFindRequest) string {
	Url, err := url.Parse("http://itunes.apple.com")
	if err != nil {
		log.Println("There was an error parsing the iTunes API Endpoint")
	}
	Url.Path += "/lookup"
	parameters := url.Values{}
	addParameter := func(key string, value string) {
		if value != "" {
			parameters.Add(key, value)
		}
	}
	addParameter("id", request.ItunesId)
	addParameter("amgArtistId", request.AmgArtistId)
	addParameter("amgAlbumId", request.AmgAlbumId)
	addParameter("amgVideoId", request.AmgVideoId)
	addParameter("entity", request.Entity)
	addParameter("limit", fmt.Sprintf("%v", request.Limit))
	addParameter("isbn", request.Isbn)
	addParameter("upc", request.Upc)
	Url.RawQuery = parameters.Encode()
	return Url.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Make search request
func ItunesSearch(request ItunesSearchRequest) ItunesResponse {
	var url = SearchUrl(request)
	res, err := http.Get(url)
	if res.StatusCode == 503 {
		log.Println("WARNING: 503 Response; pausing and trying again")
		time.Sleep(15)
		res, err = http.Get(url)
	}
	var response ItunesResponse

	if err != nil {
		log.Println("An error occurred when making the request: ", err)
		response.Raw.Err = append(response.Raw.Err, "An error occurred when making the request")
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("An error occurred when getting the contents of the response: ", err)
		response.Raw.Err = append(response.Raw.Err, "An error occurred when getting the contents of the response")
	}

	jsonerr := json.Unmarshal(contents, &response)
	if jsonerr != nil {
		log.Println("An error occurred unmarshaling the JSON: ", jsonerr)
		log.Println(res.StatusCode, contents)
		response.Raw.Err = append(response.Raw.Err, "An error occurred unmarshaling the JSON")
		response.Raw.Err = append(response.Raw.Err, jsonerr.Error())
		//if len(contents) > 10 {
		//	log.Println(contents)
		//}

	}

	if ReturnExtra {
		response.Raw.Content = contents
		response.Raw.Status = res.StatusCode
		response.Raw.Header = res.Header
		response.Raw.Url = url
	}
	return response
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Make find request
func ItunesFind(request ItunesFindRequest) ItunesResponse {
	var url = FindUrl(request)
	res, err := http.Get(url)
	if res.StatusCode == 503 {
		log.Println("WARNING: 503 Response; pausing and trying again")
		time.Sleep(15)
		res, err = http.Get(url)
	}
	var response ItunesResponse

	if err != nil {
		log.Println("An error occurred when making the request: ", err)
		response.Raw.Err = append(response.Raw.Err, "An error occurred when making the request")
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("An error occurred when getting the contents of the response: ", err)
		response.Raw.Err = append(response.Raw.Err, "An error occurred when getting the contents of the response")
	}
	jsonerr := json.Unmarshal(contents, &response)
	if jsonerr != nil {
		log.Println("An error occurred unmarshaling the JSON: ", jsonerr)
		log.Println(res.StatusCode, contents)
		response.Raw.Err = append(response.Raw.Err, "An error occurred unmarshaling the JSON")
		response.Raw.Err = append(response.Raw.Err, jsonerr.Error())
	}

	if ReturnExtra {
		response.Raw.Content = contents
		response.Raw.Status = res.StatusCode
		response.Raw.Header = res.Header
		response.Raw.Url = url
	}
	return response
}
