package i18n

// func snakeToCamel(s string) string {
// 	// Split the string by underscores
// 	words := strings.Split(s, "_")
// 	if len(words) == 0 {
// 		return ""
// 	}
//
// 	// Process each word
// 	var result strings.Builder
// 	for i, word := range words {
// 		if word == "" {
// 			continue
// 		}
// 		// Convert word to lowercase first
// 		word = strings.ToLower(word)
// 		if i == 0 {
// 			// First word remains lowercase
// 			result.WriteString(strings.Title(word))
// 		} else {
// 			// Capitalize first letter of subsequent words
// 			runes := []rune(word)
// 			runes[0] = unicode.ToUpper(runes[0])
// 			result.WriteString(string(runes))
// 		}
// 	}
// 	return result.String()
// }
