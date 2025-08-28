package styles

import "github.com/johnfercher/maroto/v2/pkg/props"

// GetPrimaryColor returns the primary color for PDF styling
func GetPrimaryColorSimplePDF() *props.Color {
	return &props.Color{
		Red:   44,
		Green: 62,
		Blue:  80,
	}
}

// GetSecondaryColor returns the secondary color for PDF styling
func GetSecondaryColorSimplePDF() *props.Color {
	return &props.Color{
		Red:   52,
		Green: 73,
		Blue:  94,
	}
}

// GetTextColor returns the text color for PDF styling
func GetTextColorSimplePDF() *props.Color {
	return &props.Color{
		Red:   52,
		Green: 73,
		Blue:  94,
	}
}
