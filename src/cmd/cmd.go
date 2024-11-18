package cmd

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-ini/ini"
	"github.com/lifeondisplay/spot/src/utils"
)

var (
	spotFolder              = getSpotFolder()
	rawFolder, themedFolder = getExtractFolder()
	backupFolder            = getBackupFolder()
	quiet                   bool
	spotifyPath             string
	prefsPath               string
	cfg                     utils.Config
	settingSection          *ini.Section
	backupSection           *ini.Section
	featureSection          *ini.Section
)

// init
func Init(isQuiet bool) {
	quiet = isQuiet

	cfg = utils.ParseConfig(filepath.Join(spotFolder, "config.ini"))
	settingSection = cfg.GetSection("Setting")
	backupSection = cfg.GetSection("Backup")
	featureSection = cfg.GetSection("AdditionalOptions")

	spotifyPath = settingSection.Key("spotify_path").String()

	if len(spotifyPath) != 0 {
		if _, err := os.Stat(spotifyPath); err != nil {
			utils.PrintError(spotifyPath + ` não existe ou não é um path válido. defina manualmente "spotify_path" em config.ini para corrigir o diretório do spotify.`)
			
			os.Exit(1)
		}
	} else if spotifyPath = utils.FindAppPath(); len(spotifyPath) != 0 {
		settingSection.Key("spotify_path").SetValue(spotifyPath)

		cfg.Write()
	} else {
		utils.PrintError(`não foi possível detectar a localização do spotify. defina manualmente "spotify_path" em config.ini`)
		
		os.Exit(1)
	}

	prefsPath = settingSection.Key("prefs_path").String()

	if len(prefsPath) != 0 {
		if _, err := os.Stat(prefsPath); err != nil {
			utils.PrintError(prefsPath + ` não existe ou não é um path válido. defina manualmente "prefs_path" em config.ini para corrigir o path do arquivo "prefs".`)
			
			os.Exit(1)
		}
	} else if prefsPath = utils.FindPrefFilePath(); len(prefsPath) != 0 {
		settingSection.Key("prefs_path").SetValue(prefsPath)
		
		cfg.Write()
	} else {
		utils.PrintError(`não foi possível detectar a localização do arquivo "prefs" do spotify. defina manualmente "prefs_path" em config.ini`)

		os.Exit(1)
	}
}

// getconfigpath
func GetConfigPath() string {
	return cfg.GetPath()
}

// getspotifypath
func GetSpotifyPath() string {
	return spotifyPath
}

func getSpotFolder() string {
	result := "/"
	
	if runtime.GOOS == "windows" {
		result = filepath.Join(os.Getenv("USERPROFILE"), ".spot")
	} else if runtime.GOOS == "linux" {
		result = filepath.Join(os.Getenv("HOME"), ".spot")
	} else if runtime.GOOS == "darwin" {
		result = filepath.Join(os.Getenv("HOME"), "spot_data")
	}

	utils.CheckExistAndCreate(result)

	return result
}

func getBackupFolder() string {
	dir := filepath.Join(spotFolder, "Backup")
	utils.CheckExistAndCreate(dir)

	return dir
}

func getExtractFolder() (string, string) {
	dir := filepath.Join(spotFolder, "Extracted")
	utils.CheckExistAndCreate(dir)

	raw := filepath.Join(dir, "Raw")
	utils.CheckExistAndCreate(raw)

	themed := filepath.Join(dir, "Themed")
	utils.CheckExistAndCreate(themed)

	return raw, themed
}

func getThemeFolder(themeName string) string {
	folder := filepath.Join(spotFolder, "Themes", themeName)
	_, err := os.Stat(folder)
	if err == nil {
		return folder
	}

	folder = filepath.Join(utils.GetExecutableDir(), "Themes", themeName)
	_, err = os.Stat(folder)
	if err == nil {
		return folder
	}

	utils.PrintError(`tema "` + themeName + `" não encontrado`)
	os.Exit(1)

	return ""
}