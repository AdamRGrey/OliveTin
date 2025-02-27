package httpservers

import (
	"encoding/json"
	//	cors "github.com/jamesread/OliveTin/internal/cors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"

	config "github.com/jamesread/OliveTin/internal/config"
)

type webUISettings struct {
	Rest      string
	ThemeName string
}

func findWebuiDir() string {
	directoriesToSearch := []string{
		"./webui",
		"/var/www/olivetin/",
	}

	for _, dir := range directoriesToSearch {
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			log.Infof("Found the webui directory here: %v", dir)

			return dir
		}
	}

	log.Warnf("Did not find the webui directory, you will probably get 404 errors.")

	return "./webui" // Should not exist
}

func generateWebUISettings(w http.ResponseWriter, r *http.Request) {
	restAddress := ""

	if !cfg.UseSingleHTTPFrontend {
		restAddress = cfg.ExternalRestAddress
	}

	jsonRet, _ := json.Marshal(webUISettings{
		Rest:      restAddress + "/api/",
		ThemeName: cfg.ThemeName,
	})

	w.Write([]byte(jsonRet))
}

func startWebUIServer(cfg *config.Config) {
	log.WithFields(log.Fields{
		"address": cfg.ListenAddressWebUI,
	}).Info("Starting WebUI server")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(findWebuiDir())))
	mux.HandleFunc("/webUiSettings.json", generateWebUISettings)

	srv := &http.Server{
		Addr:    cfg.ListenAddressWebUI,
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())
}
