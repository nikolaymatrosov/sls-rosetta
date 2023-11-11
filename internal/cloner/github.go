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

type Cloner struct {
	repo           string
	root           string
	downloadPath   string
	globsToExclude []string
}

func NewCloner(repo string, root string, downloadPath string, globsToExclude []string) Cloner {
	return Cloner{
		repo:           repo,
		root:           root,
		downloadPath:   downloadPath,
		globsToExclude: globsToExclude,
	}
}

// function that calls github api to get list files in folder of repo
func (c Cloner) GetFileList(folder string) []File {
	var files []File
	url := JoinNotEmpty("https://api.github.com/repos", c.repo, "contents", c.root, folder)

	client := resty.New()
	_, err := client.R().
		SetHeader("Accept", "application/vnd.github.v3+json").
		ForceContentType("application/json").
		SetResult(&files).
		Get(url)

	if err != nil {
		panic(err)
	}
	for i := range files {
		err := files[i].SetRelPath(c.root)
		if err != nil {
			return nil
		}
	}
	return files
}

// function that downloads files from github repo
func (c Cloner) DownloadFiles(folder string, files []File) {
	client := resty.New()
	downloadPath := c.downloadPath
	if folder != "" {
		downloadPath = JoinNotEmpty(downloadPath, folder)
	}
	client.SetOutputDirectory(downloadPath)

	for _, file := range files {
		if file.Type == "dir" {
			c.Clone(JoinNotEmpty(folder, file.Name))
			continue
		}
		// HTTP response gets saved into file, similar to curl -o flag
		_, err := client.R().
			SetOutput(file.Name).
			Get(file.DownloadUrl)

		if err != nil {
			panic(err)
		}
	}
}

// function that filters files based deploy type
func (c Cloner) FilterFiles(files []File, globsToExclude []string) []File {
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
func (c Cloner) Clone(folder string) {
	files := c.GetFileList(folder)
	filteredFiles := c.FilterFiles(files, c.globsToExclude)
	c.DownloadFiles(folder, filteredFiles)
}

func JoinNotEmpty(val ...string) string {
	var res []string
	for _, v := range val {
		if v != "" {
			res = append(res, v)
		}
	}
	return strings.Join(res, "/")
}
