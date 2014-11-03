package main

import (
    "flag"
    "fmt"
    "os"
    "github.com/jhyle/mmr/web"
)

const(
	APP_VERSION = "0.1"
)

// The flag package provides a default help printer via -h switch
var (
	versionFlag *bool = flag.Bool("v", false, "print the version number")
	portFlag *int = flag.Int("p", 3000, "port to listen on")
	hostFlag *string = flag.String("i", "127.0.0.1", "interface to listen on")
	templateDirFlag *string = flag.String("templateDir", "", "path to templates")
	imageServerFlag *string = flag.String("imageServer", "", "url of image server")
	mongoServerFlag *string = flag.String("mongoServer", "localhost", "url of mongoDb server")
	databaseFlag *string = flag.String("database", "mitmachrepublik", "name of the database")
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

    webServer, err := mmr.NewWebServer(*hostFlag, *portFlag, *templateDirFlag, *imageServerFlag, *mongoServerFlag, *databaseFlag)
    if err != nil {
    	fmt.Println(err.Error())
    } else {
		webServer.Start()
	}
}
