package luhn

import "testing"

// Checksum by Luhn Algorithm
func TestChecksum(t *testing.T) {
	type args struct {
		nums string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"0 is invalid", args{"0"}, false},
		{"1 is invalid", args{"1"}, false},
		{"11 is invalid", args{"11"}, false},
		{"70483 is invalid", args{"70483"}, false},
		{"349926205465199 is invalid", args{"349926205465199"}, false},
		{"00 is valid", args{"00"}, true},
		{"18 is valid", args{"18"}, true},
		{"70482 is valid", args{"70482"}, true},
		{"349926205465194 is valid", args{"349926205465194"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Checksum(tt.args.nums); got != tt.want {
				t.Errorf("Checksum() = %v, want %v", got, tt.want)
			}
		})
	}
}
