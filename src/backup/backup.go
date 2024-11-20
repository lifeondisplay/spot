package backup

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/lifeondisplay/spot/src/utils"
)

// fazer backup da pasta spotify apps para backuppath
func Start(spotifyPath, backupPath string) error {
	appsFolder := filepath.Join(spotifyPath, "Apps")

	return utils.Copy(appsFolder, backupPath, false, []string{".spa"})
}

// extrai todos os arquivos spa de backuppath para extractpath
// e chama `callback` em cada aplicativo extraído com sucesso
func Extract(backupPath, extractPath string, callback func(finishedApp string, err error)) {
	filepath.Walk(backupPath, func(appPath string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".spa") {
			appName := strings.Replace(info.Name(), ".spa", "", 1)
			appExtractToFolder := filepath.Join(extractPath, appName)

			err := utils.Unzip(appPath, appExtractToFolder)

			if err != nil {
				callback("", err)
			} else {
				callback(appName, nil)
			}
		}
		
		return nil
	}
}