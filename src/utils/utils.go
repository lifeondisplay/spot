package utils

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/go-ini/ini"
)

// readanswer imprime um formulário sim/não com string de `info` e
// retorna o valor booleano com base na entrada do usuário (y ou n)
// ou retorna `defaultanswer` se a entrada for omitida
//
// se o input não for nenhuma delas, imprima o formulário novamente
func ReadAnswer(info string, defaultAnswer bool) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(info)

	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\r", "", 1)
	text = strings.Replace(text, "\n", "", 1)

	if len(text) == 0 {
		return defaultAnswer
	} else if text == "y" || text == "Y" {
		return true
	} else if text == "n" || text == "N" {
		return false
	}

	return ReadAnswer(info, defaultAnswer)
}

// checa a existência de uma pasta ou então faz essa pasta caso não exista
func CheckExistAndCreate(dir string) {
	_, err := os.Stat(dir)

	if err != nil {
		os.Mkdir(dir, 0644)
	}
}

// descompacta um zip
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)

	if err != nil {
		return err
	}

	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()

		if err != nil {
			return err
		}

		defer rc.Close()
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			var fdir string

			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, f.Mode())

			if err != nil {
				log.Fatal(err)
				return err
			}

			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			defer f.Close()

			_, err = io.Copy(f, rc)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// utiliza regexp para encontrar qualquer coincidência do `input` com `regexpterm`
// e substitui por `replaceterm` e então retorna uma nova string
func Replace(input string, regexpTerm string, replaceTerm string) string {
	re := regexp.MustCompile(regexpTerm)

	return re.ReplaceAllString(input, replaceTerm)
}

// abre um arquivo, altera o conteúdo desse arquivo executando o
// callback `repl` e escreve esse novo conteúdo
func ModifyFile(path string, repl func(string) string) {
	raw, err := ioutil.ReadFile(path)

	if err != nil {
		log.Print(err)

		return
	}

	content := repl(string(raw))

	ioutil.WriteFile(path, []byte(content), 0644)
}

// encontra o path do arquivo `prefs` baseado no sistema operacional e retorna
// uma referência de `ini.file`
func GetPrefsCfg(spotifyPath string) (*ini.File, string, error) {
	var path string

	if runtime.GOOS == "windows" {
		path = filepath.Join(spotifyPath, "prefs")
	} else if runtime.GOOS == "linux" {
		path = filepath.Join(os.Getenv("HOME"), ".config", "spotify", "prefs")
	} else if runtime.GOOS == "darwin" {
		path = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Spotify", "prefs")
	} else {
		return nil, "", errors.New("Unsupported OS")
	}

	cfg, err := ini.Load(path)

	if err != nil {
		cfg = ini.Empty()
	}

	return cfg, path, nil
}

// getspotifyversion
func GetSpotifyVersion(spotifyPath string) string {
	pref, _, err := GetPrefsCfg(spotifyPath)
	if err != nil {
		log.Fatal(err)
	}

	rootSection, err := pref.GetSection("")
	if err != nil {
		log.Fatal(err)
	}

	version := rootSection.Key("app.last-launched-version")

	return version.MustString("")
}

// copia todos os arquivos ou utiliza `filters` para copiar
// certos arquivos que coincidem com as patterns
func RunCopy(from, to string, filters []string) error {
	var cmd *exec.Cmd
	var paraList = []string{from, to}

	if runtime.GOOS == "windows" {
		roboCopy := filepath.Join(os.Getenv("windir"), "System32\\robocopy.exe")

		paraList = append(paraList, "/E")

		if filters != nil && len(filters) > 0 {
			paraList = append(paraList, filters...)
		}

		cmd = exec.Command(roboCopy, paraList...)
	} else if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		var optionList = []string{"--recursive"}

		if filters != nil && len(filters) > 0 {
			optionList = append(optionList, "--include", "*/")

			for _, v := range filters {
				optionList = append(optionList, "--include", v)
			}

			optionList = append(optionList, "--exclude", "*")
		}

		optionList = append(optionList, from+"/", to)

		cmd = exec.Command("rsync", optionList...)
	} else {
		return errors.New("sistema operacional não suportado")
	}

	cmd.Run()

	return nil
}

// retorna o diretório do processo atual
func GetExecutableDir() string {
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		exe, err := os.Executable()

		if err != nil {
			log.Fatal(err)
		}
		
		return filepath.Dir(exe)
	}

	log.Fatal("sistema operacional não suportado")

	return ""
}

// retorna o diretório jshelper no diretório do executável
func GetJsHelperDir() string {
	return filepath.Join(GetExecutableDir(), "jsHelper")
}

// restartspotify
func RestartSpotify(spotifyPath string) {
	if runtime.GOOS == "windows" {
		exec.Command("taskkill", "/F", "/IM", "spotify.exe").Run()
		
		exec.Command(filepath.Join(spotifyPath, "spotify.exe")).Start()
	} else if runtime.GOOS == "linux" {
		exec.Command("pkill", "spotify").Run()

		exec.Command(filepath.Join(spotifyPath, "spotify")).Start()
	} else if runtime.GOOS == "darwin" {
		exec.Command("pkill", "Spotify").Run()
		
		exec.Command("open", "/Application/Spotify.app").Start()
	}
}