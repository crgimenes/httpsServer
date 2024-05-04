package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"
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

func handler(root string) http.Handler {
	h := http.FileServer(http.Dir(root))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security",
			"max-age=63072000; includeSubDomains")
		h.ServeHTTP(w, r)
	})
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("from %q request %q\n", r.RemoteAddr, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Carrega arquivo JSON com as configurações
	b, err := os.ReadFile("config.json")
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
		http.Handle(site.Pattern, logger(handler(site.Root)))

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
