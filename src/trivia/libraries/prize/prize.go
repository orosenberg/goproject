package prize

// Prize returns trivia wins
type Prize struct {
	Amount int
}

func (p *Prize) Calculate(change int) {
	if p.Amount+change < 0 {
		p.Amount = 0
	} else {
		p.Amount += change
	}
}
