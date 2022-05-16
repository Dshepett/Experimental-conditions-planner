package models

type Conditions struct {
	Temperature float64 `json:"temperature"`
	Time        float64 `json:"time"`
	CAcid       float64 `json:"c_acid"`
	CTi         float64 `json:"c_ti"`
	Acid        string  `json:"acid"`
	Treatment   float64 `json:"treatment"`
}
