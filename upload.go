package main

import (
	"os"

	storage "code.google.com/p/google-api-go-client/storage/v1"
)

func upload(svc *storage.Service, conf Config, from, to string) error {
	fd, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fd.Close()

	objMeta := &storage.Object{
		ContentLanguage: "en",
		ContentType:     getMimeType(from),
	}
	call := svc.Objects.Insert(conf.BucketName, objMeta).Media(fd).PredefinedAcl("publicRead")
	if _, err := call.Do(); err != nil {
		return err
	}
	return nil
}
