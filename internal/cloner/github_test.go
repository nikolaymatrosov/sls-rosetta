package cloner

import (
	"reflect"
	"testing"
)

func TestGetFileList(t *testing.T) {
	type args struct {
		folder string
	}
	tests := []struct {
		name string
		args args
		want []File
	}{
		{
			name: "examples",
			args: args{
				folder: "examples",
			},
			want: []File{
				{
					Name:    "examples.yaml",
					Path:    "examples/examples.yaml",
					RelPath: "examples.yaml",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFileList("nikolaymatrosov/sls-rosetta", tt.args.folder)

			for i := range got {
				if got[i].Name != tt.want[i].Name ||
					got[i].Path != tt.want[i].Path ||
					got[i].RelPath != tt.want[i].RelPath {
					t.Errorf("GetFileList() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestFilterFiles(t *testing.T) {
	type args struct {
		files          []File
		globsToExclude []string
	}
	tests := []struct {
		name string
		args args
		want []File
	}{
		{
			name: "glob pattern",
			args: args{
				files: []File{
					{
						Name:    "examples.yaml",
						Path:    "examples/examples.yaml",
						RelPath: "examples.yaml",
					},
					{
						Name:    "example1.yaml",
						Path:    "examples/foo/example1.yaml",
						RelPath: "foo/example1.yaml",
					},
				},
				globsToExclude: []string{"foo/*"},
			},
			want: []File{
				{
					Name:    "examples.yaml",
					Path:    "examples/examples.yaml",
					RelPath: "examples.yaml",
				},
			},
		},
		{
			name: "file name",
			args: args{
				files: []File{
					{
						Name:    "examples.yaml",
						Path:    "examples/examples.yaml",
						RelPath: "examples.yaml",
					},
					{
						Name:    "example1.yaml",
						Path:    "examples/foo/example1.yaml",
						RelPath: "foo/example1.yaml",
					},
				},
				globsToExclude: []string{"examples.yaml"},
			},
			want: []File{
				{
					Name:    "example1.yaml",
					Path:    "examples/foo/example1.yaml",
					RelPath: "foo/example1.yaml",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterFiles(tt.args.files, tt.args.globsToExclude); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
