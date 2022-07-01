package vat

func (computer) addVATAmount(p plan) float64 {
	vatAmount := p.GetTotalCost() * p.GetVATPercent() / 100
	p.SetTotalCost(p.GetTotalCost() + vatAmount)

	return vatAmount
}
