package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type CustomTheme struct{}

func (m CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {

	if variant == theme.VariantLight {
		if name == theme.ColorNameDisabled {
			return color.RGBA{130, 130, 130, 255}
		} else if name == theme.ColorNameDisabledButton {
			return color.RGBA{220, 220, 220, 255}
		} else if name == theme.ColorNameInputBackground {
			return color.RGBA{215, 215, 215, 255}

		}
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (m CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)

}
func (m CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
