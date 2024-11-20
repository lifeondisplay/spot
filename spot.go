package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/lifeondisplay/spot/src/cmd"
	"github.com/lifeondisplay/spot/src/utils"
	colorable "gopkg.in/mattn/go-colorable.v0"
)

const (
	version = "0.4.0"
)

func init() {
	if runtime.GOOS != "windows" && runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		utils.PrintError("sistema operacional não suportado.")

		os.Exit(1)
	}

	log.SetFlags(0)

	quiet := utils.FindFlag(os.Args, "-q", "--quite")

	if quiet {
		log.SetOutput(ioutil.Discard)
	} else {
		// suporte de output colorido para o windows
		log.SetOutput(colorable.NewColorableStdout())
	}

	cmd.Init(quiet)

	if len(os.Args) < 2 {
		utils.PrintInfo(`rode "spot -h" para lista de comandos.`)

		os.Exit(0)
	}

	for k, v := range os.Args {
		switch v {
			case "-c", "--config":
				fmt.Print(cmd.GetConfigPath())

				os.Exit(0)

			case "-h", "--help":
				kind := ""
				if len(os.Args) > k+1 {
					kind = os.Args[k+1]
				}

				if kind == "config" {
					helpConfig()
				} else {
					help()
				}

				os.Exit(0)

			case "-v", "--version":
				fmt.Print(version)

				os.Exit(0)
		}
	}
}

func main() {
	utils.PrintBold("spot v" + version)
	args := os.Args[1:]

	for _, argv := range args {
		switch argv {
			case "backup":
				cmd.Backup()

			case "clear":
				cmd.ClearBackup()

			case "apply":
				cmd.Apply()

			case "update":
				if utils.FindFlag(args, "-e", "--extension") {
					cmd.UpdateAllExtension()
				} else {
					cmd.UpdateCSS()
				}

			case "restore":
				cmd.Restore()

			case "enable-devtool":
				cmd.SetDevTool(true)

			case "disable-devtool":
				cmd.SetDevTool(false)

			case "watch":
				if utils.FindFlag(args, "-e", "--extension") {
					cmd.WatchExtensions()
				} else {
					cmd.Watch()
				}

			case "restart":
				cmd.RestartSpotify()

			default:
				if argv[0] != '-' {
					utils.PrintError(`comando "` + argv + `" não encontrado.`)
					utils.PrintInfo(`rode "spot -h" para uma lista dos comandos válidos.`)

					os.Exit(1)
				}
		}
	}
}

func help() {
	fmt.Println("spot v" + version)
	fmt.Print(`USO
spot [<flag>] <comando>

DESCRIÇÃO
personaliza a interface e a funcionalidade do cliente spotify

COMANDOS
backup              inicia o backup e o pré-processamento dos arquivos do aplicativo.
apply               aplica a customização.
update              atualiza o css.
restore             restaura o spotify ao estado original.
clear               limpa os arquivos de backup atuais.
enable-devtool      ativa as ferramentas de desenvolvedor do spotify. (console, inspect, ...),
                    pressione ctrl + shift + i no cliente para começar a usar.
disable-devtool     desativa as ferramentas de desenvolvedor do spotify.
watch               entra no modo de espectador. automaticamente atualiza o css quando o
                    arquivo color.ini ou user.css for alterado.
restart				reinicia o cliente spotify.

FLAGS
-q, --quiet         modo silencioso (sem output).
-e, --extension     utilize com "update" ou "watch" para focar nas extensões.
-c, --config        retorna o caminho do arquivo de configuração e saia.
-h, --help          gera este texto de ajuda e saia.
-v, --version       retorna o número da versão e saia.

para informação de configuração, rode "spot -h config".
`)
}

func helpConfig() {
	fmt.Print(`SIGNIFICADO DA CONFIGURAÇÃO
[Setting]
spotify_path
	path para o diretório spotify

prefs_path
	path para o arquivo "prefs" do spotify

current_theme
	nome da pasta do seu tema

inject_css
	se o css personalizado de user.css na pasta do tema é aplicado

replace_colors
    se as cores personalizadas são aplicadas

[Preprocesses]
disable_sentry
    impede que o sentry envie log/erro/aviso do console aos desenvolvedores do spotify.
    ative se não quiser chamar a atenção deles ao desenvolver extensão ou aplicativo.

disable_ui_logging
    vários elementos registram cada clique e rolagem do usuário.
    ative para interromper o registro e melhorar a experiência do usuário.

remove_rtl_rule
    para oferecer suporte ao árabe e outros idiomas da direita para a esquerda, o
	spotify adicionou muitas regras css que estão obsoletas para usuários da
	esquerda para a direita.
    ative para remover todos eles e melhorar a velocidade de renderização.

expose_apis
    vaza algumas apis, funções e objetos do spotify para o objeto global
	do spot que são úteis para fazer extensões para estender a
	funcionalidade do spotify.

[AdditionalOptions]
experimental_features
    permite o acesso aos recursos experimentais do spotify. abra-o no
	menu do perfil (canto superior direito).

fastUser_switching
    permite a alteração de conta imediatamente. abra-o no menu do perfil.

home
    habilita a página inicial. acesse-o na barra lateral esquerda.

lyric_always_show
    força o botão de letras para mostrar o tempo todo na barra do player.
    útil para quem deseja assistir a página de visualização.

lyric_force_no_sync
    força a exibição de todas as letras.

made_for_you_hub
    ativa a página feito para você. acesse-o na barra lateral esquerda.

radio
    habilita a página de rádio. acesse-o na barra lateral esquerda.

song_page
    clicar no nome da música na barra do player acessará a página da
	música (em vez da página do álbum) para descobrir as listas de
	reprodução nas quais ela aparece.

visualization_high_framerate
    força a visualização no aplicativo de letras para renderizar em 60fps.
`)
}