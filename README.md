# CorsXD
Its CorsXD to detect cors in websites. It simply checks the url if it has ```Access-Control-Allow-Origin``` header and if there is, it will send another request but it will change the origin header to ```Origin: https://notexisting.com```. If it has no ```Access-Control-Allow-Origin``` header in the first request, it will not proceed to make it faster

# How to install
```github.com/noobexploiter/CorsXD```

# How to use
```cat urls.txt | CorsXD```
## <br>You can specify the number of threads too
```
Usage of CorsXD:
  -t int
        Specify number of threads to run (default 32)
```
