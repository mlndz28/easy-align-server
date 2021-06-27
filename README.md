Easy Align Server
=====

Web server for Praat's plugin, Easy Align.


## Dependencies

* Easy Align's CLI must available from the system's path. For further information about installation, please refer to [praat-easy-align-linux](https://github.com/mlndz28/praat-easy-align-linux).

* [praatgo](https://github.com/mlndz28/praatgo)

```bash
go get github.com/mlndz28/praatgo
```

## Usage

```bash
cd ./easy-align-server
go build -o server
./server
```

By default, the API will be available at port 7728, and an alignment request will have this structure:

```bash
curl --location --request POST 'localhost:7728/align' \
     --form 'transcript="<transcription string>"' \
     --form 'audio=@"<audio file>"'
```


## Docker

To get a Docker container up and running:

```bash
cd ./easy-align-server
docker build --tag easy-align .
docker run --publish 7728:7728 --detach --name ea easy-align
```
