package main

type TokenValidator struct {
	Reader StringJsonReader
}

func (tv *TokenValidator) LoadTokens(fileNmae string) {
	tv.Reader.LoadData(fileNmae)
}

func (tv *TokenValidator) CheckToken(token string) bool {

	for _, vToken := range tv.Reader.Data {
		if vToken == token {
			return true
		}
	}

	return false
}

// func main() {
// 	tv := TokenValidator{r: JsonReader{}}
// 	tv.loadTokens("valid_tokens.json")
// 	log.Println(tv.checkToken("token1"))
// 	log.Println(tv.checkToken("token45"))
// }
