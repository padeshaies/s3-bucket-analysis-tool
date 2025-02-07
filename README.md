# s3-bucket-analysis-tool
A simple AWS s3 bucket analysis tool in Go. 

Go was chosen as a easy language to set up subroutines and concurency (and because I wanted to test my capabilities in this new language that I just learned 🙃)


## How to use
Windows:
```
go build
.\s3-bucket-analysis-tool.exe
```
Linux:
```
go build -o .out
.\out
```
Requirements: Have AWS config and credentials files set up in advance

## Unit tests and how to run them
Units tests for helpers have been created and can be run with the following command line
```
go test ./...
```

### Optional Flags
- `--file-size b|kb|gb|tb`, your preference for displaying file size (default: b)
- `--group-by bucket|region`, your preference for grouping results together (default: bucket)
- `--timezone`, your prefered timezone to display datetime in (default: Local)
- `--filters 'bucket-name:bucketname;storage-type:standard|intelligent_tiering|...'`, filters to apply on the bucket listing (default: none) (see [documentation](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/s3@v1.75.4/types#ObjectStorageClass) for storage type naming convention)

## TODO
- [x] parallelize everything!!! 🧑‍🌾
- [x] Get and filter by StorageType 🔍
- [x] Change how many objects which storage type
- [x] Cost helper needs some love 🤑 (in progress)
    - [ ] Mising INTELLIGENT_TIERING, SNOW and OUTPOSTS? (will need more precisions on those)


## Problems
- [ ] Fetch buckets from different regions 
- [ ] Calculating cost for Outpost and Snow storageTypes (what's a snow?!) (is Outpost virtually free 'cause it's on prem?)