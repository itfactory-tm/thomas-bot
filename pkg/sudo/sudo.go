package sudo

// TODO: make this a dynamic configurable solution

var admins = map[string]bool{
	"687715371255463972": true, // Maartje Eyskens
	"687912036595663051": true, // Ann Hannes
	"633665080994562048": true, // Ward Kerkhofs
	"688028811677138986": true, // Els Peetermans
	"688028986626146304": true, // Christine Smeets
	"371304151851728896": true, // Michiel Verboven
	"177531421152247809": true, // Dirk Mervis
}

// IsAdmin checks if a user has admin privileges
func IsAdmin(userID string) bool {
	if b, exists := admins[userID]; exists {
		return b
	}
	return false
}

var itfGameAdmins = map[string]bool{
	"161504618017325057": true, // Victor Welters
	"434499632765075456": true, // Ward Beyens
	"249632139228741632": true, // Brent (Allen?)
	"252083102992695296": true, // Alex Coulon
	"307916386238201856": true, // Jorne Marx
	"687715371255463972": true, // Maartje Eyskens
	"177531421152247809": true, // Dirk Mervis
}

func IsItfAdmin(userID string) bool {
	if value, exists := itfGameAdmins[userID]; exists {
		return value
	}
	return false
}
