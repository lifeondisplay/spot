package apply

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"../utils"
	"github.com/go-ini/ini"
)

// a flag ativa/desativa o recurso adicional
type Flag struct {
	ExperimentalFeatures bool
	FastUserSwitching    bool
	Home                 bool
	LyricAlwaysShow      bool
	LyricForceNoSync     bool
	MadeForYouHub        bool
	Radio                bool
	SongPage             bool
	VisHighFramerate     bool
	Extension            []string
}

// additionaloptions
func AdditionalOptions(appsFolderPath string, flags Flag) {
	appList, err := ioutil.ReadDir(appsFolderPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, app := range appList {
		appName := app.Name()
		appPath := filepath.Join(appsFolderPath, appName)

		err = filepath.Walk(appPath, func(path string, info os.FileInfo, err error) error {
			fileName := info.Name()
			extension := filepath.Ext(fileName)

			switch extension {
				case ".js":
					if appName == "lyrics" {
						lyricsMod(path, flags)
					} else if appName == "zlink" {
						zlinkMod(path, flags)
					}
				case ".css":
				case ".html":
					if appName == "zlink" && len(flags.Extension) > 0 {
						utils.ModifyFile(path, func(content string) string {
							for _, v := range flags.Extension {
								content = utils.Replace(content, `</body>`, `<script class="spot-extension" src="`+v+`"></script>${0}`)
							}

							return content
						})
					}
			}

			return nil
		})
	}
}

// usercss
func UserCSS(appsFolderPath, themeFolder string, injectCSS, customColor bool) {
	var userCSS string

	if customColor {
		userCSS += getColorCSS(themeFolder)
	} else {
		userCSS += getColorCSS("")
	}

	if injectCSS {
		userCSS += getUserCSS(themeFolder)
	}

	userCSSDestPath := filepath.Join(appsFolderPath, "zlink", "css", "user.css")
	ioutil.WriteFile(userCSSDestPath, []byte(userCSS), 0700)
	
	userCSSDestPath = filepath.Join(appsFolderPath, "login", "css", "user.css")
	ioutil.WriteFile(userCSSDestPath, []byte(userCSS), 0700)
}

func lyricsMod(jsPath string, flags Flag) {
	utils.ModifyFile(jsPath, func(content string) string {
		if flags.VisHighFramerate {
			content = utils.Replace(content, `[\w_]+\.highVisualizationFrameRate\s?=`, `${0}true||`)
		}

		if flags.LyricForceNoSync {
			content = utils.Replace(content, `[\w_]+\.forceNoSyncLyrics\s?=`, `${0}true||`)
		}

		return content
	})
}

func zlinkMod(jsPath string, flags Flag) {
	utils.ModifyFile(jsPath, func(content string) string {
		if flags.ExperimentalFeatures {
			content = utils.Replace(content, `[\w_]+(&&[\w_]+\.default.createElement\([\w_]+\.default,\{name:"experiments)`, `true${1}`)
		}

		if flags.FastUserSwitching {
			content = utils.Replace(content, `[\w_]+(&&[\w_]+\.default.createElement\([\w_]+\.default,\{name:"switch\-user)`, `true${1}`)
		}

		if flags.Home {
			content = utils.Replace(content, `this\._initialState\.isHomeEnabled`, "true")
			content = utils.Replace(content, `[\w_]+(&&[\w_]+\.default\.createElement\([\w_]+\.default,{isActive:\/\^spotify:app:home\/)`, "true${1}")
			content = utils.Replace(content, `[\w_]+\.isHomeEnabled`, "true")
		}

		if flags.LyricAlwaysShow {
			content = utils.Replace(content, `(lyricsEnabled\()[\w_]+&&\(.+?\)`, `${1}true`)
		}

		if flags.MadeForYouHub {
			content = utils.Replace(content, `[\w_]+(&&[\w_]+\.default.createElement\([\w_]+\.default,\{isActive:/\^spotify:app:made\-for\-you)`, `true${1}`)
		}

		if flags.Radio {
			content = utils.Replace(content, `radioIsVisible=`, `${0}true||`)
		}

		if flags.SongPage {
			content = utils.Replace(content, `window\.initialState\.isSongPageEnabled`, `true`)
		}

		return content
	})
}

func getUserCSS(themeFolder string) string {
	cssFilePath := filepath.Join(themeFolder, "user.css")
	_, err := os.Stat(cssFilePath)

	if err != nil {
		return ""
	}

	content, err := ioutil.ReadFile(cssFilePath)
	if err != nil {
		return ""
	}

	return string(content)
}

// lista de nomes de cores e seus valores padrão
var baseColorList = map[string]string{
	"main_fg":                               "ffffff",
	"secondary_fg":                          "c0c0c0",
	"main_bg":                               "282828",
	"sidebar_and_player_bg":                 "000000",
	"cover_overlay_and_shadow":              "000000",
	"indicator_fg_and_button_bg":            "1db954",
	"pressing_fg":                           "cdcdcd",
	"slider_bg":                             "404040",
	"sidebar_indicator_and_hover_button_bg": "1ed660",
	"scrollbar_fg_and_selected_row_bg":      "333333",
	"pressing_button_fg":                    "cccccc",
	"pressing_button_bg":                    "179443",
	"selected_button":                       "18ac4d",
	"miscellaneous_bg":                      "4687d6",
	"miscellaneous_hover_bg":                "2e77d0",
	"preserve_1":                            "ffffff"
}

func getColorCSS(themeFolder string) string {
	var colorCfg *ini.File
	var err error

	if len(themeFolder) == 0 {
		colorCfg = ini.Empty()
	} else {
		colorFilePath := filepath.Join(themeFolder, "color.ini")
		if colorCfg, err = ini.Load(colorFilePath); err != nil {
			colorCfg = ini.Empty()
		}
	}

	base := colorCfg.Section("Base")

	var variableList string

	for k, v := range baseColorList {
		parsed := utils.ParseColor(base.Key(k).MustString(v))
		variableList += fmt.Sprintf(`
    --modspotify_%s: #%s;
    --modspotify_rgb_%s: %s;`,
			k, parsed.Hex(),
			k, parsed.RGB())
	}

	more, err := colorCfg.GetSection("More")

	if err == nil {
		for _, v := range more.KeyStrings() {
			parsed := utils.ParseColor(more.Key(v).MustString("ffffff"))
			variableList += fmt.Sprintf(`
    --modspotify_more_%s: #%s;
    --modspotify_more_rgb_%s: %s;`,
				v, parsed.Hex(),
				v, parsed.RGB())
		}
	}

	return fmt.Sprintf(":root {%s\n}\n", variableList)
}