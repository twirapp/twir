package youtubego

import "fmt"

func ParsePlaylist(data interface{}) PlaylistParser {
	if data != nil {
		return PlaylistParser{
			IsSuccess: false,
		}
	} else {
		return PlaylistParser{
			IsSuccess: false,
		}
	}
}

func ParseChannel(data interface{}) ChannelParser {
	if data != nil {
		thumbnails := data.(map[string]interface{})["thumbnail"].(map[string]interface{})["thumbnails"]
		thumbnail := thumbnails.([]interface{})[len(thumbnails.([]interface{}))-1]

		var out ChannelParser
		out = ChannelParser{
			Channel: Channel{
				Id:   data.(map[string]interface{})["channelId"].(string),
				Name: data.(map[string]interface{})["title"].(map[string]interface{})["simpleText"].(string),
				Icon: Thumbnail{
					Url:    thumbnail.(map[string]interface{})["url"].(string),
					Width:  thumbnail.(map[string]interface{})["width"].(float64),
					Height: thumbnail.(map[string]interface{})["height"].(float64),
				},
				Subscribers: data.(map[string]interface{})["subscriberCountText"].(map[string]interface{})["simpleText"].(string),
			},
			IsSuccess: true,
		}

		return out
	} else {
		return ChannelParser{
			IsSuccess: false,
		}
	}
}

func ParseVideo(data interface{}) VideoParser {
	if data != nil {
		thumbnail := data.(map[string]interface{})["thumbnail"].(map[string]interface{})["thumbnails"].([]interface{})

		var out VideoParser
		out = VideoParser{
			Video: Video{
				Id:    data.(map[string]interface{})["videoId"].(string),
				Title: data.(map[string]interface{})["title"].(map[string]interface{})["runs"].([]interface{})[0].(map[string]interface{})["text"].(string),
				Url:   fmt.Sprintf("https://www.youtube.com/watch?v=%s", data.(map[string]interface{})["videoId"].(string)),
				Thumbnail: Thumbnail{
					Id:     data.(map[string]interface{})["videoId"].(string),
					Url:    thumbnail[len(thumbnail)-1].(map[string]interface{})["url"].(string),
					Width:  thumbnail[len(thumbnail)-1].(map[string]interface{})["width"].(float64),
					Height: thumbnail[len(thumbnail)-1].(map[string]interface{})["height"].(float64),
				},
			},
			IsSuccess: true,
		}

		return out
	} else {
		return VideoParser{
			IsSuccess: false,
		}
	}
}
