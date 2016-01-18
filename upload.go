package main

import (
	"os"

	storage "google.golang.org/api/storage/v1"
)

func upload(svc *storage.Service, conf Config, from, to string) error {
	fd, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fd.Close()

	objMeta := &storage.Object{
		Name:            to,
		ContentLanguage: "en",
		ContentType:     getMimeType(from),
	}
	call := svc.Objects.Insert(conf.BucketName, objMeta).PredefinedAcl("publicRead").Media(fd)
	if _, err := call.Do(); err != nil {
		return err
	}
	return nil
}
