package colors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColorConstants(t *testing.T) {
	tests := []struct {
		name  string
		color string
	}{
		// Primary colors
		{"ColorMauve", Mauve},
		{"ColorBlue", Blue},
		{"ColorLavender", Lavender},
		{"ColorSapphire", Sapphire},

		// Text colors
		{"ColorText", Text},
		{"ColorSubtext1", Subtext1},
		{"ColorSubtext0", Subtext0},

		// Surface colors
		{"ColorSurface0", Surface0},
		{"ColorSurface1", Surface1},
		{"ColorSurface2", Surface2},
		{"ColorBase", Base},

		// Border colors
		{"ColorOverlay2", Overlay2},
		{"ColorOverlay1", Overlay1},
		{"ColorOverlay0", Overlay0},

		// Status colors
		{"ColorGreen", Green},
		{"ColorYellow", Yellow},
		{"ColorRed", Red},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.color, "Color constant should not be empty")
			assert.True(t, len(tt.color) >= 4, "Color should be at least 4 characters (minimum hex format)")
			assert.True(t, tt.color[0] == '#', "Color should start with #")

			// Test that it's a valid hex color (basic validation)
			if len(tt.color) == 7 { // #RRGGBB format
				for _, char := range tt.color[1:] {
					assert.True(t, (char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F'),
						"Color should contain only valid hex characters")
				}
			}
		})
	}
}

func TestColorPalette_Consistency(t *testing.T) {
	// Test that all colors follow the same format
	colors := []string{
		Mauve, Blue, Lavender, Sapphire,
		Text, Subtext1, Subtext0,
		Surface0, Surface1, Surface2, Base,
		Overlay2, Overlay1, Overlay0,
		Green, Yellow, Red,
	}

	for _, color := range colors {
		assert.Equal(t, 7, len(color), "All colors should be 7 characters (#RRGGBB format)")
		assert.Equal(t, '#', rune(color[0]), "All colors should start with #")
	}
}

func TestColorPalette_CatppuccinTheme(t *testing.T) {
	// Test specific Catppuccin Mocha palette values
	expectedColors := map[string]string{
		"ColorMauve":    "#cba6f7",
		"ColorBlue":     "#89b4fa",
		"ColorLavender": "#b4befe",
		"ColorSapphire": "#74c7ec",
		"ColorText":     "#cdd6f4",
		"ColorSubtext1": "#bac2de",
		"ColorSubtext0": "#a6adc8",
		"ColorSurface0": "#313244",
		"ColorSurface1": "#45475a",
		"ColorSurface2": "#585b70",
		"ColorBase":     "#1e1e2e",
		"ColorOverlay2": "#9399b2",
		"ColorOverlay1": "#7f849c",
		"ColorOverlay0": "#6c7086",
		"ColorGreen":    "#a6e3a1",
		"ColorYellow":   "#f9e2af",
		"ColorRed":      "#f38ba8",
	}

	actualColors := map[string]string{
		"ColorMauve":    Mauve,
		"ColorBlue":     Blue,
		"ColorLavender": Lavender,
		"ColorSapphire": Sapphire,
		"ColorText":     Text,
		"ColorSubtext1": Subtext1,
		"ColorSubtext0": Subtext0,
		"ColorSurface0": Surface0,
		"ColorSurface1": Surface1,
		"ColorSurface2": Surface2,
		"ColorBase":     Base,
		"ColorOverlay2": Overlay2,
		"ColorOverlay1": Overlay1,
		"ColorOverlay0": Overlay0,
		"ColorGreen":    Green,
		"ColorYellow":   Yellow,
		"ColorRed":      Red,
	}

	for name, expected := range expectedColors {
		actual, exists := actualColors[name]
		assert.True(t, exists, "Color constant %s should exist", name)
		assert.Equal(t, expected, actual, "Color %s should match expected Catppuccin value", name)
	}
}
