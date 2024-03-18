package main

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestNewReverseProxy(t *testing.T) {
	type args struct {
		host string
		port string
	}
	tests := []struct {
		name string
		args args
		want *ReverseProxy
	}{
		{
			name: "test",
			args: args{
				host: "localhost",
				port: "8080",
			},
			want: &ReverseProxy{
				host: "localhost",
				port: "8080",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewReverseProxy(tt.args.host, tt.args.port); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReverseProxy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReverseProxy_ReverseProxy(t *testing.T) {
	type fields struct {
		host string
		port string
	}
	type args struct {
		next http.Handler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test",
			fields: fields{
				host: "localhost",
				port: "8080",
			},
			args: args{
				next: nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rp := &ReverseProxy{
				host: tt.fields.host,
				port: tt.fields.port,
			}
			got := rp.ReverseProxy(tt.args.next)
			if got == nil != tt.want {
				t.Errorf("ReverseProxy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_writeToFile(t *testing.T) {
	type args struct {
		path string
		data string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				path: "../hugo/themes/hugo-book/static/_index.md",
				data: content,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeToFile(tt.args.path, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("writeToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewReverseProxy1(t *testing.T) {
	type args struct {
		host string
		port string
	}
	tests := []struct {
		name string
		args args
		want *ReverseProxy
	}{
		{
			name: "test",
			args: args{
				host: "localhost",
				port: "8080",
			},
			want: &ReverseProxy{
				host: "localhost",
				port: "8080",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewReverseProxy(tt.args.host, tt.args.port); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReverseProxy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_geocode(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		http.Post("localhost:8080/api/address/geo", "", io.NopCloser(strings.NewReader("")))
		t.Run(tt.name, func(t *testing.T) {
			geocode(tt.args.w, tt.args.r)
		})
	}
}

func Test_search(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			search(tt.args.w, tt.args.r)
		})
	}
}

func Test_writeCounter(t *testing.T) {
	type args struct {
		counter byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeCounter(tt.args.counter); (err != nil) != tt.wantErr {
				t.Errorf("writeCounter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_writeToFile1(t *testing.T) {
	type args struct {
		path string
		data string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeToFile(tt.args.path, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("writeToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
