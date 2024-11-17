package cmd

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"../apply"
	"../status/backup"
	"../status/spotify"
	"../utils"
)

// aplicação
func Apply() {
	backupVersion := backupSection.Key("version").MustString("")
	curBackupStatus := backupstatus.Get(spotifyPath, backupFolder, backupVersion)

	if curBackupStatus == backupstatus.EMPTY {
		utils.PrintError(`você não fez nenhum backup. execute "spot backup" antes de aplicar.`)
		
		os.Exit(1)
	} else if curBackupStatus == backupstatus.OUTDATED {
		if !quiet {
			utils.PrintWarning("a versão do spotify e a versão de backup são incompatíveis.")

			if !utils.ReadAnswer("continuar com a aplicação mesmo assim? [s/n] ", false) {
				os.Exit(1)
			}
		}
	}

	appFolder := filepath.Join(spotifyPath, "Apps")
	status := spotifystatus.Get(spotifyPath)

	if status != spotifystatus.APPLIED {
		os.RemoveAll(appFolder)

		utils.Copy(rawFolder, appFolder, true, nil)
	}

	replaceColors := settingSection.Key("replace_colors").MustInt(0) == 1
	injectCSS := settingSection.Key("inject_css").MustInt(0) == 1

	if replaceColors {
		utils.Copy(themedFolder, appFolder, true, nil)
	} else {
		utils.Copy(rawFolder, appFolder, true, nil)
	}

	themeName, err := settingSection.GetKey("current_theme")
	if err != nil {
		log.Fatal(err)
	}

	themeFolder := getThemeFolder(themeName.MustString("SpotDefault"))

	apply.UserCSS(
		appFolder,
		themeFolder,
		injectCSS,
		replaceColors
	)

	extentionList := featureSection.Key("extensions").Strings("|")

	apply.AdditionalOptions(appFolder, apply.Flag{
		ExperimentalFeatures: featureSection.Key("experimental_features").MustInt(0) == 1,
		FastUserSwitching:    featureSection.Key("fastUser_switching").MustInt(0) == 1,
		Home:                 featureSection.Key("home").MustInt(0) == 1,
		LyricAlwaysShow:      featureSection.Key("lyric_always_show").MustInt(0) == 1,
		LyricForceNoSync:     featureSection.Key("lyric_force_no_sync").MustInt(0) == 1,
		MadeForYouHub:        featureSection.Key("made_for_you_hub").MustInt(0) == 1,
		Radio:                featureSection.Key("radio").MustInt(0) == 1,
		SongPage:             featureSection.Key("song_page").MustInt(0) == 1,
		VisHighFramerate:     featureSection.Key("visualization_high_framerate").MustInt(0) == 1,
		
		Extension:            extentionList
	})

	pushExtensions(extentionList...)

	utils.PrintSuccess("o spotify progrediu!")
	utils.RestartSpotify(spotifyPath)
}

// updatecss
func UpdateCSS() {
	appFolder := filepath.Join(spotifyPath, "Apps")
	themeName, err := settingSection.GetKey("current_theme")

	if err != nil {
		log.Fatal(err)
	}

	themeFolder := getThemeFolder(themeName.MustString("SpotDefault"))

	apply.UserCSS(
		appFolder,
		themeFolder,
		settingSection.Key("inject_css").MustInt(0) == 1,
		settingSection.Key("replace_colors").MustInt(0) == 1,
	)

	utils.PrintSuccess(utils.PrependTime("css customizado está atualizado"))
}

// updateallextension
func UpdateAllExtension() {
	pushExtensions(featureSection.Key("extensions").Strings("|")...)

	utils.PrintSuccess(utils.PrependTime("All extensions are updated."))
}

func getExtensionPath(name string) (string, error) {
	extFilePath := filepath.Join(utils.GetExecutableDir(), "Extensions", name)

	if _, err := os.Stat(extFilePath); err == nil {
		return extFilePath, nil
	}

	extFilePath = filepath.Join(spotFolder, "Extensions", name)

	if _, err := os.Stat(extFilePath); err == nil {
		return extFilePath, nil
	}

	return "", errors.New("extensão não encontrada")
}

func pushExtensions(list ...string) {
	if len(list) > 0 {
		zlinkFolder := filepath.Join(spotifyPath, "Apps", "zlink")

		for _, v := range list {
			extPath, err := getExtensionPath(v)

			if err != nil {
				utils.PrintError(`extensão "` + v + `" não encontrada.`)

				continue
			}

			if err = utils.CopyFile(extPath, zlinkFolder); err != nil {
				utils.PrintError(err.Error())

				continue
			}
		}
	}
}