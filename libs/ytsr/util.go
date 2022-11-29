package youtubego

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func CreateRequest(searchWord string, options SearchOptions) []SearchResult {
	Url, err := url.Parse("http://youtube.com/results")
	if err != nil {
		panic("The URL is incorrect!")
	}

	query := url.Values{}
	query.Add("search_query", searchWord)

	if strings.ToLower(options.Type) == "video" {
		query.Add("sp", "EgIQAQ%253D%253D")
	} else if strings.ToLower(options.Type) == "playlist" {
		query.Add("sp", "EgIQAw%253D%253D")
	} else if strings.ToLower(options.Type) == "channel" {
		query.Add("sp", "EgIQAg%253D%253D")
	}

	Url.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", Url.String(), nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:78.0) Gecko/20100101 Firefox/78.0")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic("Something went wrong, the request cannot be sent to the URL!")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	bodyResp, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Cannot read the body stream.")
	}

	return ParseHTML(string(bodyResp), options.Limit)
}

func ParseHTML(html string, limit int) []SearchResult {
	index := len(strings.Split(html, `{"itemSectionRenderer":`)) - 1
	items := strings.Split(html, `{"itemSectionRenderer":`)[index]
	parsed := strings.Split(items, `},{"continuationItemRenderer":{`)[0]

	var out map[string]interface{}
	err := json.Unmarshal([]byte(string(parsed)), &out)
	if err != nil {
		panic("Something went wrong, the problem was encountered while analyzing JSON!")
	}
	arr := out["contents"].([]interface{})
	output := []SearchResult{}

	if limit == 0 {
		limit = len(arr)
	}

	for i := 0; limit > i; i++ {
		sdata := arr[i].(map[string]interface{})

		if val, ok := sdata["videoRenderer"]; ok {
			parsed := ParseVideo(val)

			if parsed.IsSuccess {
				output = append(output, SearchResult{
					Video: parsed.Video,
				})
			}
		} else if val, ok := sdata["playlistRenderer"]; ok {
			parsed := ParsePlaylist(val)

			if parsed.IsSuccess {
				output = append(output, SearchResult{
					Playlist: parsed.Playlist,
				})
			}
		} else if val, ok := sdata["channelRenderer"]; ok {
			parsed := ParseChannel(val)

			if parsed.IsSuccess {
				output = append(output, SearchResult{
					Channel: parsed.Channel,
				})
			}
		}
	}

	return output
}
