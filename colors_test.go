package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"kahn/pkg/colors"
)

func TestColorConstants(t *testing.T) {
	tests := []struct {
		name  string
		color string
	}{
		// Primary colors
		{"ColorMauve", colors.Mauve},
		{"ColorBlue", colors.Blue},
		{"ColorLavender", colors.Lavender},
		{"ColorSapphire", colors.Sapphire},

		// Text colors
		{"ColorText", colors.Text},
		{"ColorSubtext1", colors.Subtext1},
		{"ColorSubtext0", colors.Subtext0},

		// Surface colors
		{"ColorSurface0", colors.Surface0},
		{"ColorSurface1", colors.Surface1},
		{"ColorSurface2", colors.Surface2},
		{"ColorBase", colors.Base},

		// Border colors
		{"ColorOverlay2", colors.Overlay2},
		{"ColorOverlay1", colors.Overlay1},
		{"ColorOverlay0", colors.Overlay0},

		// Status colors
		{"ColorGreen", colors.Green},
		{"ColorYellow", colors.Yellow},
		{"ColorRed", colors.Red},
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
		colors.Mauve, colors.Blue, colors.Lavender, colors.Sapphire,
		colors.Text, colors.Subtext1, colors.Subtext0,
		colors.Surface0, colors.Surface1, colors.Surface2, colors.Base,
		colors.Overlay2, colors.Overlay1, colors.Overlay0,
		colors.Green, colors.Yellow, colors.Red,
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
		"ColorMauve":    colors.Mauve,
		"ColorBlue":     colors.Blue,
		"ColorLavender": colors.Lavender,
		"ColorSapphire": colors.Sapphire,
		"ColorText":     colors.Text,
		"ColorSubtext1": colors.Subtext1,
		"ColorSubtext0": colors.Subtext0,
		"ColorSurface0": colors.Surface0,
		"ColorSurface1": colors.Surface1,
		"ColorSurface2": colors.Surface2,
		"ColorBase":     colors.Base,
		"ColorOverlay2": colors.Overlay2,
		"ColorOverlay1": colors.Overlay1,
		"ColorOverlay0": colors.Overlay0,
		"ColorGreen":    colors.Green,
		"ColorYellow":   colors.Yellow,
		"ColorRed":      colors.Red,
	}

	for name, expected := range expectedColors {
		actual, exists := actualColors[name]
		assert.True(t, exists, "Color constant %s should exist", name)
		assert.Equal(t, expected, actual, "Color %s should match expected Catppuccin value", name)
	}
}
