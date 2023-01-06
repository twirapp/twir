package youtubego

func Search(searchq string, options SearchOptions) ([]SearchResult, error) {
	res, err := CreateRequest(searchq, options)

	return res, err
}
