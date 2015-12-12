# gcsup

[![Build Status](https://travis-ci.org/arschles/gcsup?branch=master)](https://travis-ci.org/arschles/gcsup)
[![Docker Repository on Quay](https://quay.io/repository/arschles/gcsup/status "Docker Repository on Quay")](https://quay.io/repository/arschles/gcsup)
[![Go Report Card](http://goreportcard.com/badge/arschles/gcsup)](http://goreportcard.com/report/arschles/gcsup)


A utility for uploading folders to Google Cloud Storage. Run it simply with `./gcsup` and configure it with the following environment variables:

- `GCSUP_JWT_FILE_LOCATION` - the location of the JSON Web Token file for access to The Google Cloud Storage API
- `GCSUP_PROJECT_NAME` - the name of the Google Cloud project. The JSON Web Token should be for this project
- `GCSUP_BUCKET_NAME` - the name of the bucket inside the given project
- `GCSUP_LOCAL_FOLDER` - the name of the local folder to upload

When you run `./gcsup`, the folder at `GCSUP_LOCAL_FOLDER` will be entirely uploaded to the given bucket in the given project. The bucket will have the exact same hierarchy as the local folder, and MIME types for each file will be inferred by the Go standard library's [`TypeByExtension`](https://godoc.org/mime#TypeByExtension) function, which guesses based on each file's extension.
