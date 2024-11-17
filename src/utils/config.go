package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-ini/ini"
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

	setting, err := cfg.NewSection("Setting")
	if err != nil {
		panic(err)
	}

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

	setting.NewKey("spotify_path", defaultSpotifyPath)
	setting.NewKey("current_theme", "SpotDefault")
	setting.NewKey("inject_css", "1")
	setting.NewKey("replace_colors", "1")

	preProc, err := cfg.NewSection("Preprocesses")
	if err != nil {
		panic(err)
	}

	preProc.NewKey("disable_sentry", "1")
	preProc.NewKey("disable_ui_logging", "1")
	preProc.NewKey("remove_rtl_rule", "1")
	preProc.NewKey("expose_apis", "1")

	feature, err := cfg.NewSection("AdditionalOptions")
	if err != nil {
		panic(err)
	}

	feature.NewKey("experimental_features", "0")
	feature.NewKey("fastUser_switching", "0")
	feature.NewKey("home", "0")
	feature.NewKey("lyric_always_show", "0")
	feature.NewKey("lyric_force_no_sync", "0")
	feature.NewKey("made_for_you_hub", "0")
	feature.NewKey("radio", "0")
	feature.NewKey("song_page", "0")
	feature.NewKey("visualization_high_framerate", "0")

	version, err := cfg.NewSection("Backup")
	if err != nil {
		panic(err)
	}

	version.Comment = "NÃO ALTERAR!"
	version.NewKey("version", "")

	return cfg
}