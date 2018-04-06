package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type site struct {
	Pattern string `json:"pattern"`
	Root    string `json:"root"`
	Key     key    `json:"key"`
}

type key struct {
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

func handler(root string) http.HandlerFunc {
	h := http.FileServer(http.Dir(root))
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security",
			"max-age=63072000; includeSubDomains")
		h.ServeHTTP(w, r)
	}
}

func main() {
	// Carrega arquivo JSON com as configurações
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Faz parse do JSON para a siteList
	var siteList []site
	err = json.Unmarshal(b, &siteList)
	if err != nil {
		log.Fatal(err)
	}

	// Prepara strutura com as condigurações TLS
	tlscfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Carrega as assinaturas dos sites contidos em siteList
	for _, site := range siteList {
		http.Handle(site.Pattern, handler(site.Root))

		cert, err := tls.LoadX509KeyPair(
			site.Key.CertFile,
			site.Key.KeyFile)
		if err != nil {
			log.Fatal(err)
		}
		tlscfg.Certificates = append(tlscfg.Certificates, cert)
	}

	// Cria os nomes para os certificados
	// basicamente cria os nomes DNS para
	// cada na lista
	tlscfg.BuildNameToCertificate()

	// Prepara as configurações para o servidor
	httpServer := &http.Server{
		Addr:      ":443",
		TLSConfig: tlscfg,
	}

	// Entra em loop escutando a porta HTTPS e servindo
	// os sites da lista
	log.Fatal(httpServer.ListenAndServeTLS("", ""))
}
