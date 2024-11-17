package cmd

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"../status/spotify"
	"../utils"
)

var (
	colorCache []byte
	cssCache   []byte
)

// watch
func Watch() {
	status := spotifyStatus.Get(spotifyPath)

	if status != spotifystatus.APPLIED {
		utils.PrintError(`você não se inscreveu. rode "spot apply" uma vez antes de entrar no modo de observação.`)
		
		os.Exit(1)
	}

	themeName, err := settingSection.GetKey("current_theme")

	if err != nil {
		log.Fatal(err)
	}

	themeFolder := getThemeFolder(themeName.MustString("SpotDefault"))
	colorPath := filepath.Join(themeFolder, "color.ini")
	cssPath := filepath.Join(themeFolder, "user.css")

	for {
		shouldUpdate := false
		currColor, err := ioutil.ReadFile(colorPath)
		if err != nil {
			utils.PrintError(err.Error())
			os.Exit(1)
		}

		currCSS, err := ioutil.ReadFile(cssPath)
		if err != nil {
			utils.PrintError(err.Error())
			os.Exit(1)
		}

		if !bytes.Equal(colorCache, currColor) {
			shouldUpdate = true
			colorCache = currColor
		}

		if !bytes.Equal(cssCache, currCSS) {
			shouldUpdate = true
			cssCache = currCSS
		}

		if shouldUpdate {
			UpdateCSS()
		}

		time.Sleep(200 * time.Millisecond)
	}
}