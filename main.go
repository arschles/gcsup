package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	storage "code.google.com/p/google-api-go-client/storage/v1"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/net/context"
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
		fmt.Printf("Error with configuration [%s]\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	httpClient, err := getAuthenticatedClient(ctx, conf.JWTFileLocation)
	if err != nil {
		fmt.Printf("Error setting up authentication with Google Cloud Storage (%s)\n", err)
		os.Exit(1)
	}

	svc, err := storage.New(httpClient)
	if err != nil {
		fmt.Printf("Error creating GCS client (%s)\n", err)
		os.Exit(1)
	}

	fmt.Printf("Uploading %s to gcs://%s/%s\n", conf.LocalFolder, conf.ProjectName, conf.BucketName)

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
		go func(svc *storage.Service, conf Config, from, to string) {
			defer wg.Done()
			if err := upload(svc, conf, from, to); err != nil {
				fmt.Printf("Error uploading %s to %s (%s)\n", from, to, err)
			}
		}(svc, conf, from, to)
	}
	fmt.Println("Waiting for all uploads to finish...")
	wg.Wait()
	fmt.Println("Done")
}
