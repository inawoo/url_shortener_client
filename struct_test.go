package client

import "testing"

func Test_ComposeURLString(t *testing.T) {

	tests := []struct {
		obj     URLCollection
		want    string
		wantErr bool
	}{
		{
			obj: URLCollection{
				Host: "www.google.com",
				Path: "/search",
				Params: map[string]string{
					"q": "hello bigman",
				},
			},
			want:    "https://www.google.com/search?q=hello+bigman",
			wantErr: false,
		},
		{
			obj: URLCollection{
				Host: "www.google.com",
				Path: "/search",
				Params: map[string]string{
					"q": "hello world",
				},
			},
			want:    "https://www.google.com/search?q=hello+world",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.obj.Host, func(t *testing.T) {
			got := tt.obj.ComposeURLString()
			if got != tt.want && !tt.wantErr {
				t.Errorf("ComposeURLString() = %v, want %v", got, tt.want)
			}
		})
	}
}
