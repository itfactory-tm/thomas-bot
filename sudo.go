package main

var admins = map[string]bool{
	"687715371255463972": true, // Maartje Eyskens
	"687912036595663051": true, // Ann Hannes
	"633665080994562048": true, // Ward Kerkhofs
}

func isAdmin(userID string) bool {
	if b, exists := admins[userID]; exists {
		return b
	}
	return false
}
