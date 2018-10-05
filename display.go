package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const tplStr = `<!DOCTYPE html>
<html>

<head>
    <meta charset="UTF-8">
    <title>Album Detail</title>
</head>

<body>
    <h1><span style="width: 100%">Album Number {{ .AlbumTitle}}</span></h1>
    <table style="width: 100%">
        <thead>
            <th style="text-align:left"><strong>Picture ID</strong></th>
            <th style="text-align:left"><strong>Title</strong></th>
            <th style="text-align:left"><strong>Thumbnail</strong></th>
        </thead>
        <tbody>
            {{ range $idx, $pic := .PicDetails }}
            <tr>
                <td>{{ $pic.ID }}</td>
                <td><a href="{{ $pic.URL }}" target="_blank">{{ $pic.Title }}</a></td>
                <td><img src="{{ $pic.Thumbnail }}" alt="{{ $pic.Thumbnail }}"></td>
            </tr>
            {{else}}
            <tr>
                <td></td>
                <td>
                    <h3 style="text-align:center"><strong>No details found !!!</strong></h3>
                </td>
                <td></td>
            </tr>
            {{end}}
        </tbody>
    </table>
</body>

</html>`

// Struct to unpack json data
type picData struct {
	ID        int    `json:"id"`
	AlbumID   int    `json:"albumId"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnailUrl"`
}

// Struct for pic detail
type picDetail struct {
	ID        int
	Title     string
	URL       string
	Thumbnail string
}

// struct for display
type displayAlbumPage struct {
	AlbumTitle int
	PicDetails []picDetail
}

// map for albums
var album map[int][]picDetail

// Fetch pics details
func fetchPics() []picData {
	resp, respErr := http.Get("https://jsonplaceholder.typicode.com/photos")
	if respErr != nil {
		log.Fatal(respErr)
	}
	respBytes, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	var picdata []picData
	json.Unmarshal(respBytes, &picdata)
	resp.Body.Close()
	return picdata
}

// Map pics in according to album
func mapPics(picList []picData) map[int][]picDetail {
	var key int
	var val []picDetail
	lstLen := len(picList) - 1
	album = make(map[int][]picDetail)
	for idx, pic := range picList {
		if key != pic.AlbumID {
			if len(val) > 0 {
				album[key] = val
				val = nil
			}
			key = pic.AlbumID
		}
		val = append(val, picDetail{pic.ID, pic.Title, pic.URL, pic.Thumbnail})
		if idx == lstLen {
			album[key] = val
		}
	}
	return album
}

// Display the album details for an single album
func displayPic(albumid int) []picDetail {
	fetchedData := fetchPics()
	albumData := mapPics(fetchedData)
	return albumData[albumid]
}

//Handler func for album display
func displayAlbum(w http.ResponseWriter, r *http.Request) {
	id, parseErr := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/album/"))
	if parseErr != nil {
		log.Fatal(parseErr)
	}
	tpl, tplErr := template.New("webpage").Parse(tplStr)
	if tplErr != nil {
		log.Fatal(tplErr)
	}
	picList := displayPic(id)
	disp := displayAlbumPage{id, picList}
	err := tpl.Execute(w, disp)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/album/", displayAlbum)
	http.ListenAndServe(":8000", nil)
}
