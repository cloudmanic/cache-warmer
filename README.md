# Cache Warmer

The goal of this project is to rebuild a cache from varnish or nginx cache from a sitemap. This app will download a sitemap and then access each url in the sitemap. This is useful for when you do a complete cache flush from varnish or nginx cache.

# How to Run

Simply post your shared key, and the url to the sitemap that you want to cache warm.

```curl -H "Content-Type: application/json" -X POST -d '{"token": "<ChangeMe>","sitemap": "https://somesite.com/sitemap.xml"}' http://cache-warm.yourdomain.com```

![https://www.evernote.com/l/ABvO9A_KZkRCD4-bI98iHcjj-gD6nu8Dh8UB/image.png](https://www.evernote.com/l/ABvO9A_KZkRCD4-bI98iHcjj-gD6nu8Dh8UB/image.png)

# Docker

You do not need to use docker, but the project is tested against using docker and using [https://github.com/jwilder/nginx-proxy](https://github.com/jwilder/nginx-proxy).

# Build

To build this project and the docker container that goes with it simply run ```./build.sh```

# Configure

Just so everyone in the world can't use this we suggest you create a shared key. Open ```src/main.go``` and update ```const accessKey = "<ChangeMe>"``` with something only you know.

By default the site cache builder will try to access 10 urls at once when processing your site map. If your server can handle more you can turn this up for faster builds. If your server can't handle 10 concurrent connections you should tune this down. You configure this by changing ```const workerCount = 10``` in ```src/main.go```

More or less this is a simple app. Just change it to meet your needs. If you do something cool please share with a pull request.