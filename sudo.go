package main

var admins = map[string]bool{
	"687715371255463972": true, // Maartje Eyskens
	"687912036595663051": true, // Ann Hannes
	"633665080994562048": true, // Ward Kerkhofs
	"688028811677138986": true, // Els Peetermans
	"688028986626146304": true, // Christine Smeets
	"371304151851728896": true, // Michiel Verboven
}

func isAdmin(userID string) bool {
	if b, exists := admins[userID]; exists {
		return b
	}
	return false
}
