package pronostiek

type Rank struct {
	Name       string `json:"name"`
	Totalscore string `json:"totalscore"`
	AllCorrect int    `json:"allCorrect"`
}
