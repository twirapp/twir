package youtubego

func Search(searchq string, options SearchOptions) []SearchResult {
	res := CreateRequest(searchq, options)

	return res
}
