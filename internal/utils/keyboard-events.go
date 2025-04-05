package utils

// ParseEventId converts keyboard event identifiers to their character equivalents.
func ParseEventId(key string) string {
	if len(key) == 1 {
		toBytes := []byte(key)
		
		// Printable ASCII characters (33-126 inclusive)
		if toBytes[0] < 33 || toBytes[0] > 126 {
			return ""
		}
		
		return key
	}
	
	// Handle special key events
	switch key {
	case "<Space>":
		return " "
	case "<Tab>":
		return "\t"
	case "<Enter>":
		return "\n"
	default:
		return ""
	}
}