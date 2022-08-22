package urls

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		s   string
		sep string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "user story example",
			args: args{
				s:   "ya.ru, google.com, mts.ru",
				sep: ",",
			},
			want:    []string{"https://ya.ru", "https://google.com", "https://mts.ru"},
			wantErr: false,
		},
		{
			name: "caret",
			args: args{
				s:   "ya\nru",
				sep: ",",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "good string",
			args: args{
				s:   "https://httpbin.org/delay/10, https://httpbin.org/delay/11, http://ya.ru, google.com, mts.ru, https://pkg.go.dev/search?q=http",
				sep: ",",
			},
			want:    []string{"https://httpbin.org/delay/10", "https://httpbin.org/delay/11", "http://ya.ru", "https://google.com", "https://mts.ru", "https://pkg.go.dev/search?q=http"},
			wantErr: false,
		},
		{
			name: "single url",
			args: args{
				s:   "https://chel.mts.ru/personal",
				sep: ",",
			},
			want:    []string{"https://chel.mts.ru/personal"},
			wantErr: false,
		},
		{
			name: "wrong string",
			args: args{
				s:   "\\ya.ru, https://goo gle.com, mts^ru",
				sep: ",",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty string",
			args: args{
				s:   "",
				sep: ",",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrong sep",
			args: args{
				s:   "ya.ru, google.com, mts.ru",
				sep: ";",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.s, tt.args.sep)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
