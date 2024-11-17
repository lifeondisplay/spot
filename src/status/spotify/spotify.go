package spotifystatus

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

// enum é um tipo de constantes de status do backup
type Enum int

const (
	// stock - spotify se encontra no seu estado original
	STOCK Enum = iota

	// invalid - pasta apps possui arquivos e diretórios misturados
	INVALID

	// applied - o spotify foi modificado
	APPLIED
)

// retorna o status da pasta de apps do spotify
func Get(spotifyPath string) Enum {
	appsFolder := filepath.Join(spotifyPath, "Apps")
	fileList, err := ioutil.ReadDir(appsFolder)
	if err != nil {
		log.Fatal(err)
	}

	spaCount := 0
	dirCount := 0
	for _, file := range fileList {
		if file.IsDir() {
			dirCount++
			continue
		}

		if strings.HasSuffix(file.Name(), ".spa") {
			spaCount++
		}
	}

	totalFiles := len(fileList)
	if spaCount == totalFiles {
		return STOCK
	}

	if dirCount == totalFiles {
		return APPLIED
	}
	
	return INVALID
}