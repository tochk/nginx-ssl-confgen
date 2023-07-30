package main

import "testing"

func Test_checkPrefix(t *testing.T) {
	type args struct {
		proxyPass string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "valid URL",
			args:    args{proxyPass: "http://localhost:8080"},
			wantErr: false,
		},
		{
			name:    "https URL",
			args:    args{proxyPass: "https://example.com"},
			wantErr: false,
		},
		{
			name:    "ftp URL",
			args:    args{proxyPass: "ftp://localhost:8080"},
			wantErr: true,
		},
		{
			name:    "URL without prefix",
			args:    args{proxyPass: "localhost:8080"},
			wantErr: true,
		},
		{
			name:    "empty URL",
			args:    args{proxyPass: ""},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkPrefix(tt.args.proxyPass); (err != nil) != tt.wantErr {
				t.Errorf("checkPrefix() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
