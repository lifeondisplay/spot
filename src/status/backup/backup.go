package backupstatus

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/lifeondisplay/spot/src/utils"
)

// enum é um tipo de constantes de status do backup
type Enum int

const (
	// empty - nenhum backup encontrado
	EMPTY Enum = iota

	// backuped - há backups disponíveis
	BACKUPED

	// outdated - o backup disponível possui uma versão diferente da versão do spotify
	OUTDATED
)

// retorna o status da pasta de backup
func Get(prefsPath, backupPath, backupVersion string) Enum {
	fileList, err := ioutil.ReadDir(backupPath)
	if err != nil {
		log.Fatal(err)
	}

	if len(fileList) == 0 {
		return EMPTY
	}

	spaCount := 0
	for _, file := range fileList {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".spa") {
			spaCount++
		}
	}

	if spaCount > 0 {
		spotifyVersion := utils.GetSpotifyVersion(prefsPath)
		if backupVersion != spotifyVersion {
			return OUTDATED
		}
		
		return BACKUPED
	}

	return EMPTY
}