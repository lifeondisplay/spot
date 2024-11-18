package cmd

import (
	"log"

	"github.com/go-ini/ini"
	"github.com/lifeondisplay/spot/src/utils"
)

// setdevtool habilita/desabilita o modo de desenvolvedor no cliente do spotify
func SetDevTool(enable bool) {
	pref, err := ini.Load(prefsPath)
	if err != nil {
		log.Fatal(err)
	}

	rootSection, err := pref.GetSection("")
	if err != nil {
		log.Fatal(err)
	}

	devTool := rootSection.Key("app.enable-developer-mode")

	if enable {
		devTool.SetValue("true")
	} else {
		devTool.SetValue("false")
	}

	ini.PrettyFormat = false

	if pref.SaveTo(prefsPath) == nil {
		if enable {
			utils.PrintSuccess("devtool foi ativado!")
		} else {
			utils.PrintSuccess("devtool foi desativado!")
		}

		RestartSpotify()
	}
}