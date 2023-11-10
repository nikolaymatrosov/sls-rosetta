package cloner

import (
	"path/filepath"
	"strings"

	"github.com/go-resty/resty/v2"
)

type File struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Sha         string `json:"sha"`
	Size        int    `json:"size"`
	Url         string `json:"url"`
	HtmlUrl     string `json:"html_url"`
	GitUrl      string `json:"git_url"`
	DownloadUrl string `json:"download_url"`
	Type        string `json:"type"`
	Links       struct {
		Self string `json:"self"`
		Git  string `json:"git"`
		Html string `json:"html"`
	} `json:"_links"`
	RelPath string
}

func (f *File) SetRelPath(basePath string) error {
	relPath, err := filepath.Rel(basePath, f.Path)
	if err != nil {
		return err
	}
	f.RelPath = relPath
	return nil
}

// function that calls github api to get list files in folder of repo
func GetFileList(repo string, folder string) []File {
	var files []File

	client := resty.New()
	_, err := client.R().
		SetHeader("Accept", "application/vnd.github.v3+json").
		ForceContentType("application/json").
		SetResult(files).
		Get(strings.Join([]string{"https://api.github.com/repos", repo, "contents/examples", folder}, "/"))
	if err != nil {
		panic(err)
	}
	for i := range files {
		err := files[i].SetRelPath(folder)
		if err != nil {
			return nil
		}
	}
	return files
}

// function that downloads files from github repo
func DownloadFiles(downloadPath string, files []File) {
	client := resty.New()
	client.SetOutputDirectory(downloadPath)

	for _, file := range files {

		// HTTP response gets saved into file, similar to curl -o flag
		_, err := client.R().
			SetOutput(file.RelPath).
			Get(file.DownloadUrl)

		if err != nil {
			panic(err)
		}
	}
}

// function that filters files based deploy type
func FilterFiles(files []File, globsToExclude []string) []File {
	var filteredFiles []File

FILE_LOOP:
	for _, file := range files {

		// check if file is not in globToExclude
		for _, globToExclude := range globsToExclude {
			matched, err := filepath.Match(globToExclude, file.RelPath)
			if err != nil {
				panic(err)
			}
			if matched {
				continue FILE_LOOP
			}
		}
		filteredFiles = append(filteredFiles, file)
	}
	return filteredFiles
}

// function that clones files from github repo
func CloneFiles(
	repo string,
	folder string,
	downloadPath string,
	globsToExclude []string,
) {
	files := GetFileList(repo, folder)
	filteredFiles := FilterFiles(files, globsToExclude)
	DownloadFiles(downloadPath, filteredFiles)
}
