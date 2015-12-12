package main

import (
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
)

// Config is the envconfig compatible struct to store config values that are input to gcsup
type Config struct {
	JWTFileLocation string `required:"true" envconfig:"jwt_file_location"`
	ProjectName     string `required:"true" envconfig:"project_name"`
	BucketName      string `required:"true" envconfig:"bucket_name"`
	LocalFolder     string `required:"true" envconfig:"local_folder"`
}

// FilePath represents a file to upload to GCS
type FilePath struct {
	RootDir      string
	RelativePath string
	Name         string
}

// AbsolutePath returns the absolute path to fp on disk
func (fp FilePath) AbsolutePath() string {
	return fmt.Sprintf("%s/%s", fp.RootDir, fp.RelativePath)
}

// String is an alias to AbsolutePath
func (fp FilePath) String() string {
	return fp.AbsolutePath()
}

func main() {
	var conf Config
	if err := envconfig.Process("gcsup", &conf); err != nil {
		fmt.Printf("config error [%s]", err)
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(conf.JWTFileLocation)
	if err != nil {
		fmt.Printf("error reading file [%s]\n", err)
		os.Exit(1)
	}

	jwtConf, err := google.JWTConfigFromJSON(data, storage.ScopeFullControl)
	if err != nil {
		fmt.Printf("error creating JWT config [%s]\n", err)
		os.Exit(1)
	}

	ctx := cloud.NewContext(conf.ProjectName, jwtConf.Client(oauth2.NoContext))
	var files []FilePath
	if err := filepath.Walk(conf.LocalFolder, func(path string, fInfo os.FileInfo, err error) error {
		if fInfo.IsDir() {
			return nil
		}
		relPath := strings.TrimPrefix(path, conf.LocalFolder+"/")
		files = append(files, FilePath{RootDir: conf.LocalFolder, RelativePath: relPath, Name: fInfo.Name()})
		return nil
	}); err != nil {
		fmt.Printf("error gathering all files [%s]", err)
		os.Exit(1)
	}
	for _, file := range files {
		from := file.AbsolutePath()
		to := file.RelativePath
		fmt.Printf("uploading %s to %s\n", from, to)
		upload(ctx, conf, from, to)
	}
}

func upload(ctx context.Context, conf Config, from string, to string) error {
	w := storage.NewWriter(ctx, conf.BucketName, to)
	defer w.Close()
	w.ACL = []storage.ACLRule{
		storage.ACLRule{Entity: storage.AllUsers, Role: storage.RoleReader},
	}
	extension := from[strings.LastIndex(from, "."):]
	w.ContentType = mime.TypeByExtension(extension)

	fileBytes, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}
	if _, err := w.Write(fileBytes); err != nil {
		return err
	}
	return nil
}
