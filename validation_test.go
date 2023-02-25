package client

import "testing"

func Test_validateURLSyntax(t *testing.T) {
	type args struct {
		url string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "valid url", args: args{url: "https://www.google.com"}, wantErr: false},
		{name: "valid url", args: args{url: "https://www.google.com/search?q=hello+world&oq=hello+world&aqs=chrome..69i57j0l7.1001j0j7&sourceid=chrome&ie=UTF-8"}, wantErr: false},
		{name: "invalid url", args: args{url: "www.google.com/search?q=hello+world&oq=hello+world&aqs=chrome..69i57j0l7.1001j0j7&sourceid=chrome&ie=UTF-8"}, wantErr: true},
		{name: "invalid url", args: args{url: "not a url"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateURLSyntax(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("validateURLSyntax() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateURLExists(t *testing.T) {
	type args struct {
		url string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "valid url", args: args{url: "https://www.google.com"}, wantErr: false},
		{name: "valid url", args: args{url: "https://www.google.com/search?q=hello+world&oq=hello+world&aqs=chrome..69i57j0l7.1001j0j7&sourceid=chrome&ie=UTF-8"}, wantErr: false},
		{name: "invalid url", args: args{url: "www.gt446oogle.com/search?q=hello+world&oq=hello+world&aqs=chrome..69i57j0l7.1001j0j7&sourceid=chrome&ie=UTF-8"}, wantErr: true},
		{name: "invalid url", args: args{url: "not a url"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateURLExists(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("url = %v; error = %v; wantErr %v", tt.args.url, err, tt.wantErr)
			}
		})
	}
}

func Test_validateURLLength(t *testing.T) {
	type args struct {
		url string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "invalid url", args: args{url: "https:/m"}, wantErr: true},
		{name: "valid url", args: args{url: "https://www.google.com/search?q=hello+world&oq=hello+world&aqs=chrome..69i57j0l7.1001j0j7&sourceid=chrome&ie=UTF-8"}, wantErr: false},
		{name: "valid url", args: args{url: "www.google.com/search?q=hello+world&oq=hello+world&aqs=chrome..69i57j0l7.1001j0j7&sourceid=chrome&ie=UTF-8"}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateURLLength(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("validateURLLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
