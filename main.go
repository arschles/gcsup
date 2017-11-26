package main

import (
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/context"

	"cloud.google.com/go/storage"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// Config is the envconfig compatible struct to store config values that are input to gcsup
type Config struct {
	JWTFileLocation string `required:"true" envconfig:"jwt_file_location"`
	ProjectName     string `required:"true" envconfig:"project_name"`
	BucketName      string `required:"true" envconfig:"bucket_name"`
	LocalFolder     string `required:"true" envconfig:"local_folder"`
}

func main() {
	var conf Config
	if err := envconfig.Process("gcsup", &conf); err != nil {
		fmt.Printf("Error with configuration [%s]", err)
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(conf.JWTFileLocation)
	if err != nil {
		fmt.Printf("Error reading file [%s]\n", err)
		os.Exit(1)
	}

	jwtConf, err := google.JWTConfigFromJSON(data, storage.ScopeFullControl)
	if err != nil {
		fmt.Printf("Error creating JWT config [%s]\n", err)
		os.Exit(1)
	}

	fmt.Printf("Uploading %s to gcs://%s/%s\n", conf.LocalFolder, conf.ProjectName, conf.BucketName)

	ctx := cloud.NewContext(conf.ProjectName, jwtConf.Client(oauth2.NoContext))
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Printf("Error creating GCS client [%s]\n", err)
		os.Exit(1)
	}
	bucket := client.Bucket(conf.BucketName)

	files, err := getAllFiles(conf.LocalFolder)
	if err != nil {
		fmt.Printf("Error gathering all files [%s]\n", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	for _, file := range files {
		from := file.AbsolutePath
		to := strings.TrimPrefix(file.RelativePath, conf.LocalFolder)
		fmt.Printf("Uploading %s to %s\n", from, to)
		wg.Add(1)
		go func(ctx context.Context, bucket *storage.BucketHandle, conf Config, from, to string) {
			defer wg.Done()
			if err := upload(ctx, client.Bucket(conf.BucketName), conf, from, to); err != nil {
				fmt.Printf("Error uploading %s to %s (%s)\n", from, to, err)
			}
		}(ctx, bucket, conf, from, to)
	}
	fmt.Println("Waiting for all uploads to finish...")
	wg.Wait()
	fmt.Println("Done")
}

func upload(ctx context.Context, bucket *storage.BucketHandle, conf Config, from, to string) error {

	fileBytes, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	obj := bucket.Object(to)
	w := obj.NewWriter(ctx)
	if _, err := w.Write(fileBytes); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		fmt.Printf("ERROR closing writer for upload %s => %s (%s)\n", from, to, err)
		return err
	}

	extension := from[strings.LastIndex(from, "."):]
	attrs := storage.ObjectAttrs{
		ACL: []storage.ACLRule{
			storage.ACLRule{
				Entity: storage.AllUsers,
				Role:   storage.RoleReader,
			},
		},
		ContentType: mime.TypeByExtension(extension),
	}

	if _, err := obj.Update(ctx, attrs); err != nil {
		return err
	}

	return nil
}
