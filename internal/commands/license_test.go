package commands

import (
	"testing"
)

func TestGenerateLicense(t *testing.T) {
	tests := []struct {
		name        string
		licenseType string
		author      string
		year        string
		wantErr     bool
		wantContain string
	}{
		{
			name:        "MIT License",
			licenseType: "mit",
			author:      "Test Author",
			year:        "2024",
			wantErr:     false,
			wantContain: "MIT License",
		},
		{
			name:        "Apache2 License",
			licenseType: "apache2",
			author:      "Test Author",
			year:        "2024",
			wantErr:     false,
			wantContain: "Apache License",
		},
		{
			name:        "GPL3 License",
			licenseType: "gpl3",
			author:      "Test Author",
			year:        "2024",
			wantErr:     false,
			wantContain: "GNU GENERAL PUBLIC LICENSE",
		},
		{
			name:        "BSD3 License",
			licenseType: "bsd3",
			author:      "Test Author",
			year:        "2024",
			wantErr:     false,
			wantContain: "BSD 3-Clause License",
		},
		{
			name:        "Invalid License",
			licenseType: "invalid",
			author:      "Test Author",
			year:        "2024",
			wantErr:     true,
			wantContain: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateLicense(tt.licenseType, tt.author, tt.year)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateLicense() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("generateLicense() returned empty string")
			}
			if tt.wantContain != "" && !containsString(got, tt.wantContain) {
				t.Errorf("generateLicense() does not contain %v", tt.wantContain)
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr || containsString(s[1:], substr)))
}
