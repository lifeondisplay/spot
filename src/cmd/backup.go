package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/lifeondisplay/spot/src/backup"
	"github.com/lifeondisplay/spot/src/preprocess"
	backupstatus "github.com/lifeondisplay/spot/src/status/backup"
	"github.com/lifeondisplay/spot/src/utils"
)

// backup
func Backup() {
	backupVersion := backupSection.Key("version").MustString("")
	curBackupStatus := backupstatus.Get(prefsPath, backupFolder, backupVersion)

	if curBackupStatus != backupstatus.EMPTY {
		utils.PrintWarning("há backup disponível, mas limpe o backup atual primeiro!")
		ClearBackup()

		backupSection.Key("version").SetValue("")
		cfg.Write()
	}

	utils.PrintBold("fazendo backup de arquivos de aplicativos:")

	if err := backup.Start(prefsPath, backupFolder); err != nil {
		log.Fatal(err)
	}

	appList, err := ioutil.ReadDir(backupFolder)
	if err != nil {
		log.Fatal(err)
	}

	totalApp := len(appList)
	if totalApp > 0 {
		utils.PrintGreen("OK")
	} else {
		utils.PrintError("não foi possível fazer backup dos arquivos do aplicativo. reinstale o spotify e tente novamente.")

		os.Exit(1)
	}

	utils.PrintBold("Extracting:")
	tracker := utils.NewTracker(totalApp)

	if quiet {
		tracker.Quiet()
	}

	backup.Extract(backupFolder, rawFolder, tracker.Update)
	tracker.Finish()

	preprocSec := cfg.GetSection("pré-processos")

	tracker.Reset()

	utils.PrintBold("pré-processando:")

	preprocess.Start(
		rawFolder,
		preprocess.Flag{
			DisableSentry:  preprocSec.Key("disable_sentry").MustInt(0) == 1,
			DisableLogging: preprocSec.Key("disable_ui_logging").MustInt(0) == 1,
			RemoveRTL:      preprocSec.Key("remove_rtl_rule").MustInt(0) == 1,
			ExposeAPIs:     preprocSec.Key("expose_apis").MustInt(0) == 1,
		},
		tracker.Update,
	)

	tracker.Finish()

	utils.Copy(rawFolder, themedFolder, true, []string{".html", ".js", ".css"})

	tracker.Reset()

	preprocess.StartCSS(themedFolder, tracker.Update)
	tracker.Finish()

	backupSection.Key("version").SetValue(utils.GetSpotifyVersion(spotifyPath))
	cfg.Write()
	utils.PrintSuccess("está tudo pronto, você pode começar a aplicação agora!")
}

// clearbackup
func ClearBackup() {
	if !quiet {
		if !utils.ReadAnswer("antes de limpar o backup, certifique-se de ter restaurado ou reinstalado o spotify ao estado original. continuar? [s/n]: ", false) {
			os.Exit(1)
		}
	}

	os.RemoveAll(backupFolder)
	os.RemoveAll(rawFolder)
	os.RemoveAll(themedFolder)

	backupSection.Key("version").SetValue("")
	cfg.Write()
}

// restore
func Restore() {
	backupVersion := backupSection.Key("version").MustString("")
	curBackupStatus := backupstatus.Get(prefsPath, backupFolder, backupVersion)

	if curBackupStatus == backupstatus.EMPTY {
		utils.PrintError(`você não fez backup.`)

		os.Exit(1)
	} else if curBackupStatus == backupstatus.OUTDATED {
		if !quiet {
			utils.PrintWarning("a versão do spotify e a versão de backup são incompatíveis.")

			if !utils.ReadAnswer("continuar restaurando mesmo assim? [s/n] ", false) {
				os.Exit(1)
			}
		}
	}

	appFolder := filepath.Join(spotifyPath, "Apps")

	os.RemoveAll(appFolder)
	utils.Copy(backupFolder, appFolder, false, []string{".spa"})

	utils.PrintSuccess("o spotify foi restaurado.")
	
	RestartSpotify()
}