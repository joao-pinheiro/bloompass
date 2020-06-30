## Bloom Filter Password Checker

Beware! PoC! Not to be used in production!

Implementation of a bloom filter available via API to check if a given password exists on a given blacklist.  

Blacklists are text files with lowercase .txt extension, one password per line.  

### Details

The bloom filter is built with separate roaring bitmaps, one for each hashing type;
Currently, murmur3 and fnv are used for hashing;


### Usage
```shell script
$ make
$ bin/bloompass -dir data/
```

(wait until "Start API Server" appears)

### JSON API

The API only supports SEARCH method. 

```shell script
$ curl --header "Content-Type: application/json" \
  --request SEARCH \
  --data '{"password":"xyz"}' \
  http://localhost:3030/check

{"success":true,"exists":1}

```

Valid return codes for "exists":

| code| meaning|
|---|---|
| 0 | Does not exist in the blacklist|
| 1 | Exists in the blacklist|

### TODO
- Debugging
- Unit Tests
- Performance improvements