package main

import (
	"flag"
	"fmt"
	"github.com/jhyle/mitmachrepublik/web"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

const (
	APP_VERSION = "0.1"
)

// The flag package provides a default help printer via -h switch
var (
	versionFlag              *bool   = flag.Bool("v", false, "print the version number")
	envFlag                  *string = flag.String("env", "dev", "dev, test or www")
	portFlag                 *int    = flag.Int("p", 3000, "port to listen on")
	hostFlag                 *string = flag.String("i", "127.0.0.1", "interface to listen on")
	templateDirFlag          *string = flag.String("templateDir", "", "path to templates")
	indexDirFlag             *string = flag.String("indexDir", "", "path for search index")
	imageServerFlag          *string = flag.String("imageServer", "", "url of image server")
	mongoServerFlag          *string = flag.String("mongoServer", "localhost", "url of mongoDb server")
	databaseFlag             *string = flag.String("database", "mitmachrepublik", "name of the database")
	smtpPassFlag             *string = flag.String("smtpPassword", "", "password for mitmachrepublik@gmail.com")
	fbAppSecret              *string = flag.String("fbAppSecret", "", "Facebook App Secret")
	fbUser                   *string = flag.String("fbUser", "", "Facebook user")
	fbPassword               *string = flag.String("fbPassword", "", "Facebook password")
	twitterApiSecret         *string = flag.String("twitterApiSecret", "", "Twitter Api Secret")
	twitterAccessTokenSecret *string = flag.String("twitterAccessTokenSecret", "", "Twitter Api Access Token Secret")
	gibUser                  *string = flag.String("gibUser", "", "Gratis in Berlin user")
	gibPassword              *string = flag.String("gibPassword", "", "Gratis in Berlin password")
	scrapersFlag             *bool   = flag.Bool("scrapers", false, "run scrapers")
	postEventFlag            *bool   = flag.Bool("postEvent", false, "run post event")
)

func IsFolder(path string) bool {

	folder, err := os.Stat(path)
	if err != nil {
		return false
	}

	return folder.IsDir()
}

func main() {
	flag.Parse() // Scan the arguments list

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
	}

	if *envFlag != "dev" && *envFlag != "test" && *envFlag != "www" {
		fmt.Println("Invalid environment specified!")
		os.Exit(1)
	}

	if *hostFlag == "" {
		fmt.Println("You need to specify an interface to listen on (-i)!")
		os.Exit(1)
	}

	if *templateDirFlag == "" {
		fmt.Println("You need to specify a template directory (-templateDir)!")
		os.Exit(1)
	}

	if !IsFolder(*templateDirFlag) {
		fmt.Println("Given template directory (-templateDir=" + *templateDirFlag + ") is not a directory!")
		os.Exit(1)
	}

	if *indexDirFlag == "" {
		fmt.Println("You need to specify a search index directory (-indexDir)!")
		os.Exit(1)
	}

	if !IsFolder(*indexDirFlag) {
		fmt.Println("Given search index directory (-indexDir=" + *indexDirFlag + ") is not a directory!")
		os.Exit(1)
	}

	if *imageServerFlag == "" {
		fmt.Println("You need to specify an image server url (-imageServer)!")
		os.Exit(1)
	}

	if *mongoServerFlag == "" {
		fmt.Println("You need to specify a mongo server url (-mongoServer)!")
		os.Exit(1)
	}

	if *databaseFlag == "" {
		fmt.Println("You need to specify a database (-database)!")
		os.Exit(1)
	}

	if *smtpPassFlag == "" {
		fmt.Println("You need to specify a SMTP password (-smtpPassword)!")
		os.Exit(1)
	}

	app, err := mmr.NewMmrApp(*envFlag, *hostFlag, *portFlag, *templateDirFlag, *indexDirFlag, *imageServerFlag, *mongoServerFlag, *databaseFlag, *smtpPassFlag, *fbAppSecret, *fbUser, *fbPassword, *twitterApiSecret, *twitterAccessTokenSecret, *gibUser, *gibPassword)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	if *scrapersFlag == false && *postEventFlag == false {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGTERM)
		go func() {
			<-c
			app.Stop()
			os.Exit(1)
		}()

		// start debugging server
		go func() {
			http.ListenAndServe("localhost:6060", nil)
		}()

		app.Start()
	} else if *scrapersFlag {
		err = app.RunScrapers()
		if err != nil {
			fmt.Println(err)
		}
		app.Stop()
	} else if *postEventFlag {
		err = app.RunPostEvent()
		if err != nil {
			fmt.Println(err)
		}
		app.Stop()
	}
}
