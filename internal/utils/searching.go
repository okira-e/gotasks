package utils


// IncludesFuzzy takes ("hi there", "itere") and returns true.
func IncludesFuzzy(str string, phrase string) bool {
	phraseIndex := 0
	for _, strChar := range str {
		if phraseIndex == len(phrase) - 1 {
			return true
		}
		
		phrasechar := string(phrase[phraseIndex])
		
		if phrasechar == string(strChar) {
			phraseIndex += 1
		}
	}
	
	return false
}