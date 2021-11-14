# magazi
simple in memory key-value store.

## Local Installation
Clone the repository. Build the program with `go build`.
There are several options for this program. You can see the options with `./magazi -help`
```
Usage of ./magazi:
  -file string
    	name of the backup file (default "backup-data")
  -port int
    	port to listen (default 9242)
  -time int
    	interval time as minutes (default 1)
```

Each parameter have default value. If you want to define your own you can fill the variable like this: 
```
./magazi -file=<file> -port=<port>  -time=<time>
```

## Containerization
Run build.sh with specfying version: `./build.sh <my_version>`
Now we can run our docker image as container. Run:
```
docker run -d -p 9242:9242 magazi:<my_version>
```
You can see the running container with `docker -ps`.

## API Usage
At this time our app is running in `localhost:9242`.
### AddData /POST
You can add data to store with a request. Request body should be the following way:
```
{
    "key": "data-key",
    "value": "data-value"
}
```
### GetData /GET
You can get a data from store with a GET request. In the request parameter you need to specify the key. For the key `data-key` you can get the value with:
``` 
localhost:9242/?key=data-key
```
Json response will be the following structure:
```
{
    "key": "data-key",
    "value": "data-value"
}
```

### FlushData /POST
You can flush the data to backup json file with a POST request. You need POST this request:
```
localhost:9242/flush
```
Now backup file will be synced with current data store.

## Requests Log
Each http request will be saved to file called `requests.log` under the `/tmp` folder. You can analyze all the requests with their properties.

## Tests
You can run the unit tests with `./run-tests.sh`