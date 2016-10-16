package main

import (
	"archive/zip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"rsc.io/letsencrypt"

	"github.com/beevik/etree"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// Metadata metadata struct
type Metadata struct {
	Title      string    `json:"title"`
	Author     string    `json:"author"`
	Identifier string    `json:"identifier"`
	Language   string    `json:"language"`
	Modified   time.Time `json:"modified"`
}

// Link link struct
type Link struct {
	Rel      string `json:"rel,omitempty"`
	Href     string `json:"href"`
	TypeLink string `json:"type"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

// Manifest manifest struct
type Manifest struct {
	Metadata  Metadata `json:"metadata"`
	Links     []Link   `json:"links"`
	Spine     []Link   `json:"spine,omitempty"`
	Resources []Link   `json:"resources"`
}

// Icon icon struct for AppInstall
type Icon struct {
	Src       string `json:"src"`
	Size      string `json:"size"`
	MediaType string `json:"type"`
}

// AppInstall struct for app install banner
type AppInstall struct {
	ShortName string `json:"short_name"`
	Name      string `json:"name"`
	StartURL  string `json:"start_url"`
	Display   string `json:"display"`
	Icons     Icon   `json:"icons"`
}

func main() {

	n := negroni.Classic()
	n.Use(negroni.NewStatic(http.Dir("public")))
	n.UseHandler(loanHandler(false))

	var m letsencrypt.Manager
	if err := m.CacheFile("letsencrypt.cache"); err != nil {
		log.Fatal(err)
	}

	if os.Args[1] == "dev" {
		s := &http.Server{
			Handler:        n,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
			Addr:           ":8080",
		}

		log.Fatal(s.ListenAndServe())
	} else {

		s := &http.Server{
			Handler:        n,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
			Addr:           ":https",
			TLSConfig: &tls.Config{
				GetCertificate: m.GetCertificate,
			},
		}

		log.Fatal(s.ListenAndServeTLS("", ""))
	}
}

func loanHandler(test bool) http.Handler {
	serv := mux.NewRouter()

	serv.HandleFunc("/{filename}/manifest.json", getManifest)
	serv.HandleFunc("/{filename}/webapp.webmanifest", getWebAppManifest)
	serv.HandleFunc("/{filename}/index.html", bookIndex)
	serv.HandleFunc("/{filename}/{asset:.*}", getAsset)
	return serv
}

func getManifest(w http.ResponseWriter, req *http.Request) {
	var opfFileName string
	var scheme string
	var manifestStruct Manifest
	var metaStruct Metadata

	metaStruct.Modified = time.Now()

	vars := mux.Vars(req)
	filename := vars["filename"]

	self := Link{
		Rel:      "self",
		Href:     "http://" + req.Host + "/manifest/" + filename + "/manifest.json",
		TypeLink: "application/epub+zip",
	}
	manifestStruct.Links = make([]Link, 1)
	manifestStruct.Resources = make([]Link, 0)
	manifestStruct.Resources = make([]Link, 0)
	manifestStruct.Links[0] = self

	zipReader, err := zip.OpenReader(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, f := range zipReader.File {
		if f.Name == "META-INF/container.xml" {
			rc, errOpen := f.Open()
			if errOpen != nil {
				fmt.Println("error openging " + f.Name)
			}
			doc := etree.NewDocument()
			_, err = doc.ReadFrom(rc)
			if err == nil {
				root := doc.SelectElement("container")
				rootFiles := root.SelectElements("rootfiles")
				for _, rootFileTag := range rootFiles {
					rootFile := rootFileTag.SelectElement("rootfile")
					if rootFile != nil {
						opfFileName = rootFile.SelectAttrValue("full-path", "")
					}
				}
			} else {
				fmt.Println(err)
			}
			rc.Close()
		}
	}

	if opfFileName != "" {
		for _, f := range zipReader.File {
			if f.Name == opfFileName {
				rc, errOpen := f.Open()
				if errOpen != nil {
					fmt.Println("error openging " + f.Name)
				}
				doc := etree.NewDocument()
				_, err = doc.ReadFrom(rc)
				if err == nil {
					root := doc.SelectElement("package")
					meta := root.SelectElement("metadata")

					titleTag := meta.SelectElement("title")
					metaStruct.Title = titleTag.Text()

					langTag := meta.SelectElement("language")
					metaStruct.Language = langTag.Text()

					identifierTag := meta.SelectElement("identifier")
					metaStruct.Identifier = identifierTag.Text()

					creatorTag := meta.SelectElement("creator")
					metaStruct.Author = creatorTag.Text()

					bookManifest := root.SelectElement("manifest")
					itemsManifest := bookManifest.SelectElements("item")
					for _, item := range itemsManifest {
						linkItem := Link{}
						linkItem.TypeLink = item.SelectAttrValue("media-type", "")
						if os.Args[1] == "dev" {
							scheme = "http://"
						} else {
							scheme = "https://"
						}
						linkItem.Href = scheme + req.Host + "/" + filename + "/" + item.SelectAttrValue("href", "")
						if linkItem.TypeLink == "application/xhtml+xml" {
							manifestStruct.Spine = append(manifestStruct.Spine, linkItem)
						} else {
							manifestStruct.Resources = append(manifestStruct.Resources, linkItem)
						}
					}

					manifestStruct.Metadata = metaStruct
					j, _ := json.Marshal(manifestStruct)
					w.Header().Set("Content-Type", "application/json")
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Write(j)
					return
				}
			}
		}
	}

}

func getAsset(w http.ResponseWriter, req *http.Request) {
	var opfFileName string

	vars := mux.Vars(req)
	filename := vars["filename"]
	assetname := vars["asset"]

	zipReader, err := zip.OpenReader(filename)
	if err != nil {
		fmt.Println(err)
	}

	for _, f := range zipReader.File {
		if f.Name == "META-INF/container.xml" {
			rc, errOpen := f.Open()
			if errOpen != nil {
				fmt.Println("error openging " + f.Name)
			}
			doc := etree.NewDocument()
			_, err = doc.ReadFrom(rc)
			if err == nil {
				root := doc.SelectElement("container")
				rootFiles := root.SelectElements("rootfiles")
				for _, rootFileTag := range rootFiles {
					rootFile := rootFileTag.SelectElement("rootfile")
					if rootFile != nil {
						opfFileName = rootFile.SelectAttrValue("full-path", "")
					}
				}
			} else {
				fmt.Println(err)
			}
			rc.Close()
		}
	}

	resourcePath := strings.Split(opfFileName, "/")[0]

	for _, f := range zipReader.File {
		fmt.Println(f.Name)
		if f.Name == resourcePath+"/"+assetname {
			rc, errOpen := f.Open()
			if errOpen != nil {
				fmt.Println("error openging " + f.Name)
			}
			buff, _ := ioutil.ReadAll(rc)
			defer rc.Close()
			extension := filepath.Ext(f.Name)
			if extension == ".css" {
				w.Header().Set("Content-Type", "text/css")
			}
			if extension == ".xml" {
				w.Header().Set("Content-Type", "application/xhtml+xml")
			}
			if extension == ".js" {
				w.Header().Set("Content-Type", "text/javascript")
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(buff)
			return
		}
	}

}

func getWebAppManifest(w http.ResponseWriter, req *http.Request) {
	var opfFileName string
	var webapp AppInstall

	vars := mux.Vars(req)
	filename := vars["filename"]

	webapp.Display = "standalone"
	webapp.StartURL = "index.html"
	webapp.Icons = Icon{
		Size:      "144x144",
		Src:       "logo.png",
		MediaType: "image/png",
	}

	zipReader, err := zip.OpenReader(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, f := range zipReader.File {
		if f.Name == "META-INF/container.xml" {
			rc, errOpen := f.Open()
			if errOpen != nil {
				fmt.Println("error openging " + f.Name)
			}
			doc := etree.NewDocument()
			_, err = doc.ReadFrom(rc)
			if err == nil {
				root := doc.SelectElement("container")
				rootFiles := root.SelectElements("rootfiles")
				for _, rootFileTag := range rootFiles {
					rootFile := rootFileTag.SelectElement("rootfile")
					if rootFile != nil {
						opfFileName = rootFile.SelectAttrValue("full-path", "")
					}
				}
			} else {
				fmt.Println(err)
			}
			rc.Close()
		}
	}

	if opfFileName != "" {
		for _, f := range zipReader.File {
			if f.Name == opfFileName {
				rc, errOpen := f.Open()
				if errOpen != nil {
					fmt.Println("error openging " + f.Name)
				}
				doc := etree.NewDocument()
				_, err = doc.ReadFrom(rc)
				if err == nil {
					root := doc.SelectElement("package")
					meta := root.SelectElement("metadata")

					titleTag := meta.SelectElement("title")
					webapp.Name = titleTag.Text()
					webapp.ShortName = titleTag.Text()

					j, _ := json.Marshal(webapp)
					w.Header().Set("Content-Type", "application/json")
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Write(j)
					return
				}
			}
		}
	}

}

func bookIndex(w http.ResponseWriter, req *http.Request) {
	var err error

	vars := mux.Vars(req)
	filename := vars["filename"]

	t, err := template.ParseFiles("index.html") // Parse template file.
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, filename) // merge.
}
