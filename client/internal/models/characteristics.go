package models

type Characteristics struct {
	Size        float64 `json:"size"`
	Consistence float64 `json:"consistance"`
	Stability   float64 `json:"stability"`
}

func (c *Characteristics) IsValid() bool {
	return c.Size >= 0 && c.Consistence >= 0 && c.Consistence <= 100 && c.Stability >= 0
}
