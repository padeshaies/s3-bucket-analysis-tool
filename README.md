# s3-bucket-analysis-tool
A simple AWS s3 bucket analysis tool in Go


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
- `--filters bucket-name:bucketname;storage-type:standard|ia|rr|...`, filters to apply of the bucket listing (default: none)

## Problems
- [ ] Fetch buckets from different regions