package utils


// @Incomplete: Only accepts English characters right now.
func ParseEventId(key string) string {
	if len(key) == 1 {
		toBytes := []byte(key)
		
		if toBytes[0] < 33 && toBytes[0] > 126 {
			return ""
		}
		
		return key
	}
	
	if key == "<Space>" {
		return " "
	}
	
	if key == "<Tab>" {
		return "\t"
	}
	
	if key == "<Enter>" {
		return "\n"
	}
	
	return ""
}