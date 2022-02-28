package model

// Structs
// The sun
type Sun struct {
	Mass                float64 `json:"mass"`                //solar masses
	Luminosity          float64 `json:"luminosity"`          // by size and effective temperature
	Color               string  `json:"color"`               // it depends on class
	Class               string  `json:"class"`               // star classification
	Maximum_temperature uint16  `json:"maximum_temperature"` // Celsius degree
	Minimum_temperature uint16  `json:"minimum_temperature"` // Celsius degree
	Age                 uint16  `json:"age"`                 // billion years
	R_ecosphere         float64 `json:"r_ecosphere"`
}

// funcs
func r_ecosphere() float64 {
	return 1.15
}

func Sunme() Sun {
	result := Sun{
		Mass:                1,
		Luminosity:          2,
		Color:               "green",
		Class:               "K",
		Maximum_temperature: 1000,
		Minimum_temperature: 500,
		Age:                 126,
		R_ecosphere:         r_ecosphere(),
	}
	return result

}
