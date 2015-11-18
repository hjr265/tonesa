# TONESA

This sample app is published as part of the blog article at [www.toptal.com/blog](http://www.toptal.com/blog).

## Trying It Out

Clone the repository into a Go environment:

~~~
$ mkdir tonesa
$ cd tonesa
$ export GOPATH=`pwd`
$ mkdir -p src/github.com/hjr265/tonesa
$ cd src/github.com/hjr265/tonesa
$ git clone https://github.com/hjr265/tonesa.git .
$ go get ./...
~~~

Compile tonesad:

~~~
$ go build ./cmd/tonesad
~~~

Create .env file:

~~~
$ cp env-sample.txt .env
$ nano .env
~~~

~~~
MONGO_URL=mongodb://127.0.0.1/tonesa
REDIS_URL=redis://127.0.0.1
AWS_ACCESS_KEY_ID={Your-AWS-Access-Key-ID-Goes-Here}
AWS_SECRET_ACCESS_KEY={And-Your-AWS-Secret-Access-Key}
S3_BUCKET_NAME={And-S3-Bucket-Name}
~~~

Finally run it:

~~~
$ PORT=9091 ./tonesad -env-file=.env
~~~

## License

Node WHOIS is available under the [BSD (2-Clause) License](http://opensource.org/licenses/BSD-2-Clause).

## What's up with this name?

Believe me, I would have given this a better name if I could. I was just **t**ired **o**f **n**aming **e**very **s**ingle **a**pp.
