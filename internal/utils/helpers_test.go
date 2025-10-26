package utils

import (
	"testing"
)

func TestHashString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		algorithm string
		want      string
		wantErr   bool
	}{
		{
			name:      "MD5 hash",
			input:     "test",
			algorithm: "md5",
			want:      "098f6bcd4621d373cade4e832627b4f6",
			wantErr:   false,
		},
		{
			name:      "SHA1 hash",
			input:     "test",
			algorithm: "sha1",
			want:      "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3",
			wantErr:   false,
		},
		{
			name:      "SHA256 hash",
			input:     "test",
			algorithm: "sha256",
			want:      "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			wantErr:   false,
		},
		{
			name:      "Invalid algorithm",
			input:     "test",
			algorithm: "invalid",
			want:      "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashString(tt.input, tt.algorithm)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HashString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneratePassword(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantLen int
		wantErr bool
	}{
		{
			name:    "Generate 16 char password",
			length:  16,
			wantLen: 16,
			wantErr: false,
		},
		{
			name:    "Generate 32 char password",
			length:  32,
			wantLen: 32,
			wantErr: false,
		},
		{
			name:    "Generate 8 char password",
			length:  8,
			wantLen: 8,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePassword(tt.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("GeneratePassword() length = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "Existing file (this test file)",
			path: "helpers_test.go",
			want: true,
		},
		{
			name: "Non-existing file",
			path: "nonexistent-file-12345.txt",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileExists(tt.path); got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
