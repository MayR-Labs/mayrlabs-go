package commands

import (
	"testing"
)

func TestGenerateEditorConfig(t *testing.T) {
	tests := []struct {
		name        string
		language    string
		wantContain string
	}{
		{
			name:        "Go config",
			language:    "go",
			wantContain: "indent_style = tab",
		},
		{
			name:        "JavaScript config",
			language:    "javascript",
			wantContain: "indent_size = 2",
		},
		{
			name:        "Python config",
			language:    "python",
			wantContain: "indent_size = 4",
		},
		{
			name:        "PHP config",
			language:    "php",
			wantContain: "indent_size = 4",
		},
		{
			name:        "General config",
			language:    "general",
			wantContain: "root = true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateEditorConfig(tt.language)
			if got == "" {
				t.Errorf("generateEditorConfig() returned empty string")
			}
			if !contains(got, tt.wantContain) {
				t.Errorf("generateEditorConfig() does not contain %v", tt.wantContain)
			}
		})
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
