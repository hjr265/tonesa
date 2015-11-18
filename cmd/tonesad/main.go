package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/desertbit/glue"
	"github.com/hjr265/tonesa/api"
	"github.com/hjr265/tonesa/data"
	"github.com/hjr265/tonesa/hub"
	"github.com/hjr265/tonesa/ui"
)

func main() {
	fEnvFile := flag.String("env-file", "", "path to environment file")
	flag.Parse()

	if *fEnvFile != "" {
		err := LoadEnvFile(*fEnvFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	err := data.OpenSession("mongodb://127.0.0.1:27017/tonesa")
	if err != nil {
		log.Fatal(err)
	}
	err = data.MakeIndexes()
	if err != nil {
		log.Fatal(err)
	}

	err = data.InitBucket(os.Getenv("S3_BUCKET_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	err = hub.InitHub(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal(err)
	}

	glueSrv := glue.NewServer(glue.Options{
		HTTPSocketType: glue.HTTPSocketTypeNone,
	})
	glueSrv.OnNewSocket(hub.HandleSocket)

	http.Handle("/", ui.Router)
	http.Handle("/api/", http.StripPrefix("/api", api.Router))
	http.Handle("/assets/", http.StripPrefix("/assets", ui.AssetsFS))
	http.Handle("/glue/", glueSrv)

	port := os.Getenv("PORT")
	log.Printf("Listening on :%s", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
