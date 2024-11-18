package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-ini/ini"
)

var (
	configLayout = map[string]map[string]string {
		"Setting": map[string]string {
			"spotify_path":   "",
			"current_theme":  "SpotDefault",
			"inject_css":     "1",
			"replace_colors": "1"
		},

		"Preprocesses": map[string]string{
			"disable_sentry":     "1",
			"disable_ui_logging": "1",
			"remove_rtl_rule":    "1",
			"expose_apis":        "1"
		},

		"AdditionalOptions": map[string]string{
			"experimental_features":        "0",
			"fastUser_switching":           "0",
			"home":                         "0",
			"lyric_always_show":            "0",
			"lyric_force_no_sync":          "0",
			"made_for_you_hub":             "0",
			"radio":                        "0",
			"song_page":                    "0",
			"visualization_high_framerate": "0",
			"extensions":                   ""
		}
	}
)

type config struct {
	path    string
	content *ini.File
}

// config
type Config interface {
	Write()

	GetSection(string) *ini.Section
	GetPath() string
}

// parseConfig lê o conteúdo do arquivo de configuração, retorna
// a configuração padrão se o arquivo não existir
func ParseConfig(configPath string) Config {
	cfg, err := ini.LoadSources(
		ini.LoadOptions{
			IgnoreContinuation: true,
		}, configPath)

	if err != nil {
		defaultConfig := config{
			path:    configPath,
			content: getDefaultConfig()
		}

		defaultConfig.Write()

		PrintSuccess("config.ini padrão gerado.")

		return defaultConfig
	}

	needRewrite := false

	for sectionName, keyList := range configLayout {
		section, err := cfg.GetSection(sectionName)
		if err != nil {
			section, _ = cfg.NewSection(sectionName)
			needRewrite = true
		}

		for keyName, defaultValue := range keyList {
			if _, err := section.GetKey(keyName); err != nil {
				section.NewKey(keyName, defaultValue)
				needRewrite = true
			}
		}
	}

	if needRewrite {
		PrintSuccess("configuração atualizada.")

		cfg.SaveTo(configPath)
	}

	return config {
		path:    configPath,
		content: cfg
	}
}

// escreve o conteúdo do arquivo de configuração
func (c config) Write() {
	c.content.SaveTo(c.path)
}

func (c config) GetSection(name string) *ini.Section {
	sec, err := c.content.GetSection(name)

	if err != nil {
		Fatal(err)
	}
	
	return sec
}

func (c config) GetPath() string {
	return c.path
}

func getDefaultConfig() *ini.File {
	var cfg = ini.Empty()
	var defaultSpotifyPath string
	
	if runtime.GOOS == "windows" {
		defaultSpotifyPath = filepath.Join(os.Getenv("APPDATA"), "Spotify")
		_, err := os.Stat(defaultSpotifyPath)
		if err != nil {
			PrintError(`o conteúdo do spotify não é encontrado no local "` + defaultSpotifyPath + `"`)
			PrintInfo(`certifique-se de estar usando a versão normal do spotify, e não a versão da windows store.`)
		}
	} else if runtime.GOOS == "linux" {
		path, err := exec.Command("whereis", "spotify").Output()

		if err == nil {
			pathString := strings.Replace(string(path), "spotify: ", "", 1)

			pathString = strings.Replace(pathString, "\n", "", -1)
			pathList := strings.Split(pathString, " ")

			for _, v := range pathList {
				_, err := os.Stat(filepath.Join(v, "Apps"))

				if err == nil {
					defaultSpotifyPath = v
					break
				}
			}
		}

		if len(defaultSpotifyPath) == 0 {
			PrintWarning("não foi possível detectar a localização do spotify.")
		}
	} else if runtime.GOOS == "darwin" {
		defaultSpotifyPath = filepath.Join("/Applications", "Spotify.app", "Contents", "Resources")

		_, err := os.Stat(defaultSpotifyPath)
		if err != nil {
			PrintError(`o conteúdo do spotify não é encontrado no local "` + defaultSpotifyPath + `"`)
		}
	}

	configLayout["Setting"]["spotify_path"] = defaultSpotifyPath

	for sectionName, keyList := range configLayout {
		section, err := cfg.NewSection(sectionName)
		if err != nil {
			panic(err)
		}

		for keyName, defaultValue := range keyList {
			section.NewKey(keyName, defaultValue)
		}
	}

	version, err := cfg.NewSection("Backup")
	if err != nil {
		panic(err)
	}

	version.Comment = "NÃO ALTERAR!"
	version.NewKey("version", "")

	return cfg
}