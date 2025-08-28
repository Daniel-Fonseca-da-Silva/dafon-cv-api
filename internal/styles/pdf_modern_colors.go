package styles

import "github.com/johnfercher/maroto/v2/pkg/props"

// GetPrimaryColor returns the primary color for PDF styling
func GetPrimaryColorModernPDF() *props.Color {
	return &props.Color{
		Red:   17,
		Green: 24,
		Blue:  39,
	}
}

// GetSecondaryColor returns the secondary color for PDF styling
func GetSecondaryColorModernPDF() *props.Color {
	return &props.Color{
		Red:   59,
		Green: 130,
		Blue:  246,
	}
}

// GetTextColor returns the text color for PDF styling
func GetTextColorModernPDF() *props.Color {
	return &props.Color{
		Red:   55,
		Green: 65,
		Blue:  81,
	}
}

// GetAccentColor returns the accent color for highlights
func GetAccentColorModernPDF() *props.Color {
	return &props.Color{
		Red:   99,
		Green: 102,
		Blue:  241,
	}
}

// GetLightGrayColor returns a light gray color for backgrounds
func GetLightGrayColorModernPDF() *props.Color {
	return &props.Color{
		Red:   249,
		Green: 250,
		Blue:  251,
	}
}

// GetBorderColor returns the border color for sections
func GetBorderColorModernPDF() *props.Color {
	return &props.Color{
		Red:   229,
		Green: 231,
		Blue:  235,
	}
}

// GetMutedTextColor returns a muted text color for secondary information
func GetMutedTextColorModernPDF() *props.Color {
	return &props.Color{
		Red:   107,
		Green: 114,
		Blue:  128,
	}
}
