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

### Optional Flags
- `--file-size b|kb|gb|tb`, your preference for displaying file size (default: b)
- `--group-by bucket|region`, your preference for grouping results together (default: bucket)
- `--timezone`, your prefered timezone to display datetime in (default: Local)
- `--filters bucket-name:bucketname;storage-type:standard|ia|rr|...`, filters to apply of the bucket listing (default: none)

## TODO
- [x] parallelize everything!!! 🧑‍🌾
- [ ] Get and filter by StorageType 🔍
- [ ] Cost helper needs some love 🤑


## Problems
- [ ] Fetch buckets from different regions 