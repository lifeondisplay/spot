package preprocess

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"../utils"
)

// a flag habilita/desabilita os pré-processos a serem aplicados
type Flag struct {
	DisableSentry  bool
	DisableLogging bool
	RemoveRTL      bool
	ExposeAPIs     bool
}

// inicia o pré-processamento ativo de aplicativos em extractapppath
//
// disablesentry: impede que o sentry envie log/erro/aviso do console para desenvolvedores do spotify
// disablelogging: interrompe vários elementos para registrar a interação do usuário
// removertl: remove todas as regras css da direita para a esquerda para simplificar os arquivos css
func Start(extractedAppsPath string, flags Flag, callback func(appName string, err error)) {
	appList, err := ioutil.ReadDir(extractedAppsPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, app := range appList {
		appName := app.Name()
		appPath := filepath.Join(extractedAppsPath, appName)

		err = filepath.Walk(appPath, func(path string, info os.FileInfo, err error) error {
			fileName := info.Name()
			extension := filepath.Ext(fileName)

			switch extension {
				case ".js":
					utils.ModifyFile(path, func(content string) string {
						if flags.DisableSentry {
							content = disableSentry(content)
						}

						if flags.DisableLogging {
							content = disableLogging(content, appName)
						}

						if appName == "zlink" && flags.ExposeAPIs {
							content = exposeAPIs(content)
						}

						return content
					})

				case ".css":
					if fileName == "glue.css" && appName != "zlink" && appName != "login" {
						os.Remove(path)
						return nil
					}

					if flags.RemoveRTL {
						utils.ModifyFile(path, removeRTL)
					}

				case ".html":
					utils.ModifyFile(path, func(content string) string {
						if appName != "zlink" && appName != "login" {
							content = utils.Replace(content, `css/glue\.css`, "https://zlink.app.spotify.com/css/glue.css")
							content = utils.Replace(content, `</head>`, `<link rel="stylesheet" class="userCSS" href="https://zlink.app.spotify.com/css/user.css"></head>`)
						} else {
							content = utils.Replace(content, `</head>`, `<link rel="stylesheet" class="userCSS" href="css/user.css"></head>`)
						}

						if appName == "zlink" && flags.ExposeAPIs {
							content = utils.Replace(content, `(<script src="init\.js"></script>)`, `${1}<script type="text/javascript" src="/spotWrapper.js"></script>`)
						}

						return content
					})
			}

			return nil
		})

		if appName == "zlink" && flags.ExposeAPIs {
			utils.RunCopy(".", appPath, false, []string{"spotWrapper.js"})
		}

		if err != nil {
			callback("", err)
		} else {
			callback(appName, nil)
		}
	}
}

// startcss modifica todos os arquivos css em extractappspath para
// alterar todos os valores de cores com variáveis ​​css
func StartCSS(extractedAppsPath string, callback func(appName string, err error)) {
	appList, err := ioutil.ReadDir(extractedAppsPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, app := range appList {
		appName := app.Name()
		appPath := filepath.Join(extractedAppsPath, appName)

		err = filepath.Walk(appPath, func(path string, info os.FileInfo, err error) error {
			if filepath.Ext(info.Name()) == ".css" {
				utils.ModifyFile(path, func(content string) string {
					content = utils.Replace(content, "#1ed660", "var(--modspotify_sidebar_indicator_and_hover_button_bg)")
					content = utils.Replace(content, "#1ed760", "var(--modspotify_sidebar_indicator_and_hover_button_bg)")
					content = utils.Replace(content, "#1db954", "var(--modspotify_indicator_fg_and_button_bg)")
					content = utils.Replace(content, "#1df369", "var(--modspotify_indicator_fg_and_button_bg)")
					content = utils.Replace(content, "#1df269", "var(--modspotify_indicator_fg_and_button_bg)")
					content = utils.Replace(content, "#1cd85e", "var(--modspotify_indicator_fg_and_button_bg)")
					content = utils.Replace(content, "#1bd85e", "var(--modspotify_indicator_fg_and_button_bg)")
					content = utils.Replace(content, "#18ac4d", "var(--modspotify_selected_button)")
					content = utils.Replace(content, "#18ab4d", "var(--modspotify_selected_button)")
					content = utils.Replace(content, "#179443", "var(--modspotify_pressing_button_bg)")
					content = utils.Replace(content, "#14833b", "var(--modspotify_pressing_button_bg)")
					content = utils.Replace(content, "#282828", "var(--modspotify_main_bg)")
					content = utils.Replace(content, "#121212", "var(--modspotify_main_bg)")
					content = utils.Replace(content, "#999999", "var(--modspotify_main_bg)")
					content = utils.Replace(content, "#606060", "var(--modspotify_main_bg)")
					content = utils.Replace(content, `rgba\(18,\s?18,\s?18,\s?([\d\.]+)\)`, "rgba(var(--modspotify_rgb_sidebar_and_player_bg),${1})")
					content = utils.Replace(content, "#181818", "var(--modspotify_sidebar_and_player_bg)")
					content = utils.Replace(content, `rgba\(18,\s?19,\s?20,\s?([\d\.]+)\)`, "rgba(var(--modspotify_rgb_sidebar_and_player_bg),${1})")
					content = utils.Replace(content, "#000000", "var(--modspotify_sidebar_and_player_bg)")
					content = utils.Replace(content, "#333333", "var(--modspotify_scrollbar_fg_and_selected_row_bg)")
					content = utils.Replace(content, "#3f3f3f", "var(--modspotify_scrollbar_fg_and_selected_row_bg)")
					content = utils.Replace(content, "#535353", "var(--modspotify_scrollbar_fg_and_selected_row_bg)")
					content = utils.Replace(content, "#404040", "var(--modspotify_slider_bg)")
					content = utils.Replace(content, `rgba\(80,\s?55,\s?80,\s?([\d\.]+)\)`, "rgba(var(--modspotify_rgb_sidebar_and_player_bg),${1})")
					content = utils.Replace(content, `rgba\(40,\s?40,\s?40,\s?([\d\.]+)\)`, "rgba(var(--modspotify_rgb_sidebar_and_player_bg),${1})")
					content = utils.Replace(content, `rgba\(40,\s?40,\s?40,\s?([\d\.]+)\)`, "rgba(var(--modspotify_rgb_sidebar_and_player_bg),${1})")
					content = utils.Replace(content, `rgba\(24,\s?24,\s?24,\s?([\d\.]+)\)`, "rgba(var(--modspotify_rgb_sidebar_and_player_bg),${1})")
					content = utils.Replace(content, `rgba\(18,\s?19,\s?20,\s?([\d\.]+)\)`, "rgba(var(--modspotify_rgb_sidebar_and_player_bg),${1})")
					content = utils.Replace(content, "#000011", "var(--modspotify_sidebar_and_player_bg)")
					content = utils.Replace(content, "#0a1a2d", "var(--modspotify_sidebar_and_player_bg)")
					content = utils.Replace(content, "#ffffff", "var(--modspotify_main_fg)")
					content = utils.Replace(content, "#f8f8f7", "var(--modspotify_pressing_fg)")
					content = utils.Replace(content, "#fcfcfc", "var(--modspotify_pressing_fg)")
					content = utils.Replace(content, "#d9d9d9", "var(--modspotify_pressing_fg)")
					content = utils.Replace(content, "#cdcdcd", "var(--modspotify_pressing_fg)")
					content = utils.Replace(content, "#e6e6e6", "var(--modspotify_pressing_fg)")
					content = utils.Replace(content, "#e5e5e5", "var(--modspotify_pressing_fg)")
					content = utils.Replace(content, "#adafb2", "var(--modspotify_secondary_fg)")
					content = utils.Replace(content, "#c8c8c8", "var(--modspotify_secondary_fg)")
					content = utils.Replace(content, "#a0a0a0", "var(--modspotify_secondary_fg)")
					content = utils.Replace(content, "#bec0bb", "var(--modspotify_secondary_fg)")
					content = utils.Replace(content, "#bababa", "var(--modspotify_secondary_fg)")
					content = utils.Replace(content, "#b3b3b3", "var(--modspotify_secondary_fg)")
					content = utils.Replace(content, "#c0c0c0", "var(--modspotify_secondary_fg)")
					content = utils.Replace(content, `rgba\(179,\s?179,\s?179,\s?([\d\.]+)\)`, "rgba(var(--modspotify_rgb_secondary_fg),${1})")
					content = utils.Replace(content, "#cccccc", "var(--modspotify_pressing_button_fg)")
					content = utils.Replace(content, "#ededed", "var(--modspotify_pressing_button_fg)")
					content = utils.Replace(content, "#4687d6", "var(--modspotify_miscellaneous_bg)")
					content = utils.Replace(content, `rgba\(70,\s?135,\s?214,\s?([\d\.]+)\)`, "rgba(var(--modspotify_rgb_miscellaneous_bg),${1})")
					content = utils.Replace(content, "#2e77d0", "var(--modspotify_miscellaneous_hover_bg)")
					content = utils.Replace(content, `rgba\(51,\s?153,\s?255,\s?([\d\.]+)\)`, "rgba(var(--modspotify_rgb_miscellaneous_hover_bg),${1})")
					content = utils.Replace(content, `rgba\(30,\s?50,\s?100,\s?([\d\.]+)\)`, "rgba(var(--modspotify_rgb_miscellaneous_hover_bg),${1})")
					content = utils.Replace(content, `rgba\(24,\s?24,\s?24,\s?([\d\.]+)\)`, "rgba(var(--modspotify_rgb_sidebar_and_player_bg),${1})")
					content = utils.Replace(content, `rgba\(25,\s?20,\s?20,\s?([\d\.]+)\)`, "rgba(var(--modspotify_rgb_sidebar_and_player_bg),${1})")
					content = utils.Replace(content, `rgba\(160,\s?160,\s?160,\s?([\d\.]+)\)`, "rgba(var(--modspotify_rgb_pressing_button_fg),${1})")
					content = utils.Replace(content, `rgba\(255,\s?255,\s?255,\s?([\d\.]+)\)`, "rgba(var(--modspotify_rgb_pressing_button_fg),${1})")
					content = utils.Replace(content, "#ddd", "var(--modspotify_pressing_button_fg)")
					content = utils.Replace(content, "#000", "var(--modspotify_sidebar_and_player_bg)")
					content = utils.Replace(content, "#333", "var(--modspotify_scrollbar_fg_and_selected_row_bg)")
					content = utils.Replace(content, "#444", "var(--modspotify_slider_bg)")
					content = utils.Replace(content, "#fff", "var(--modspotify_main_fg)")
					content = utils.Replace(content, "black;", " var(--modspotify_sidebar_and_player_bg)")
					content = utils.Replace(content, "gray;", " var(--modspotify_main_bg)")
					content = utils.Replace(content, "lightgray;", " var(--modspotify_pressing_button_fg)")
					content = utils.Replace(content, "white;", " var(--modspotify_main_fg)")
					content = utils.Replace(content, `rgba\(0,\s?0,\s?0,\s?([\d\.]+)\)`, "rgba(var(--modspotify_rgb_cover_overlay_and_shadow),${1})")
					content = utils.Replace(content, "#fff", "var(--modspotify_main_fg)")
					content = utils.Replace(content, "#000", "var(--modspotify_sidebar_and_player_bg)")
					return content
				})
			}

			return nil
		})

		if err != nil {
			callback("", err)
		} else {
			callback(appName, nil)
		}
	}
}

func disableSentry(input string) string {
	input = utils.Replace(input, `sentry\.install\(\)[,;]`, "")
	input = utils.Replace(input, `"https://\w+@sentry.io/\d+"`, `"https://NO@TELEMETRY.IS/BAD"`)

	return input
}

func disableLogging(input, appName string) string {
	input = utils.Replace(input, `data\-log\-click="[\w\-]+"`, "")
	input = utils.Replace(input, `data\-log\-context="[\w\-]+"`, "")

	switch appName {
		case "browse", "collection", "genre", "hub":
			input = utils.Replace(input, `logUIInteraction5\([\w_]+,\s?[\w_]+\)\s?\{`, "${0}return;")
			input = utils.Replace(input, `logUIImpression5\([\w_]+,\s?[\w_]+\)\s?\{`, "${0}return;")
			input = utils.Replace(input, `_logUIInteraction5\([\w_]+\)\s?\{`, "${0}return;")
			input = utils.Replace(input, `_logUIImpression5\([\w_]+\)\s?\{`, "${0}return;")
			input = utils.Replace(input, `this\._documentFragment\.query\(['"]\[data\-log\-click\]['"]\)`, "return;${0}")
			input = utils.Replace(input, `_onClickDataLogClick\([\w_]+\)\s?\{`, "${0}return;")
			input = utils.Replace(input, `_setUpStandardImpressionLogging\(\)\s?\{`, "${0}return;")

		case "zlink":
			input = utils.Replace(input, `prototype\._logUIInteraction5=function\(.+?\)\{`, "${0}return;")

		case "lyrics":
			input = utils.Replace(input, `\.prototype\.log.+?\{`, "${0}return;")

		case "playlist":
			input = utils.Replace(input, `logPlaylistImpression=function\(.+?\)\s?\{`, "${0}return;")
			input = utils.Replace(input, `logEndOfListImpression=function\(.+?\)\s?\{`, "${0}return;")
			input = utils.Replace(input, `logListQuickJump=function\(.+?\)\s?\{`, "${0}return;")
			input = utils.Replace(input, `logListItemSelected=function\(.+?\)\{`, "${0}return;")
			input = utils.Replace(input, `logFeedbackInteraction=function\(.+?\)\s?\{`, "${0}return;")
			// For ver 1.80
			input = utils.Replace(input, `(exports\.logPlaylistImpression =) \w+`, "${1}void")
			input = utils.Replace(input, `(exports\.logEndOfListImpression =) \w+`, "${1}void")
			input = utils.Replace(input, `(exports\.logListQuickJump =) \w+`, "${1}void")
			input = utils.Replace(input, `(exports\.logListItemSelected =) \w+`, "${1}void")
			input = utils.Replace(input, `(exports\.logFeedbackInteraction =) \w+`, "${1}void")
	}

	return input
}

func removeRTL(input string) string {
	input = utils.Replace(input, `,\s?\[dir=rtl\].+?(\{.+?\})`, "$1")
	input = utils.Replace(input, `\[dir=rtl\].+?\{.+?\}`, "")
	input = utils.Replace(input, `\[dir=ltr\]`, "")
	input = utils.Replace(input, `\[dir\]`, "")
	input = utils.Replace(input, `,\s?\[lang=ar\].+?(\{.+?\})`, "$1")
	input = utils.Replace(input, `\[lang=ar\].+?\{.+?\}`, "")
	input = utils.Replace(input, `\[dir=\"rtl\"\].+?\{.+?\}`, "")
	input = utils.Replace(input, `html\[dir=\"rtl\"\].+?\{.+?\}`, "")

	return input
}

func findSymbol(debugInfo, content string, clues []string) []string {
	for _, v := range clues {
		re := regexp.MustCompile(v)
		found := re.FindStringSubmatch(content)

		if found != nil {
			return found[1:]
		}
	}

	utils.PrintColor("red", false, "\nCannot find symbol for "+debugInfo)

	return nil
}

func exposeAPIs(input string) string {
	playerUI := findSymbol("playerUI", input, []string{
		`([\w_]+)\.prototype\.updateProgressBarLabels`,
		`([\w_]+)\.prototype\._onConnectionStateChange`})

	if playerUI != nil {
		input = utils.Replace(
			input,
			playerUI[0]+`\.prototype\.setup=function\(\)\{`,
			"${0}"+spotPlayerJS)

		// registra o evento de mudança de progresso
		input = utils.Replace(
			input,
			playerUI[0]+`\.prototype\._onProgressBarProgress=function.+?\{`,
			`${0}Spot.Player.dispatchEvent&&Spot.Player.dispatchEvent(new Event("onprogress"));`)
	}

	// vazamento de metadados da faixa, estado do player e lista de reprodução atual para o spot.player.data
	input = utils.Replace(
		input,
		`(const [\w_]+=([\w_]+)\.track\.metadata;)`,
		`${1}Spot.Player.data=${2};`)

	// encontra o símbolo event dispatcher (eventsymbol[0]) e event creator (eventsymbol[1])
	eventSymbols := findSymbol("EventDispatcher and Event Creatir", input, []string{
		`([\w_]+)\.default\.dispatchEvent\(new ([\w_]+)\.default\([\w_]+\.default\.NAVIGATION_OPEN_URI`,
		`([\w_]+)\.default\.dispatchEvent\(new ([\w_]+)\.default\("show\-notification\-bubble"`})

	showNotification := ""
	if eventSymbols != nil {
		showNotification = fmt.Sprintf(
			`Spot.showNotification = (text) => {%s.default.dispatchEvent(new %s.default("show-notification-bubble", {i18n: text}))};`,
			eventSymbols[0],
			eventSymbols[1])
	}

	// vazamento de localstorage e shownotification
	input = utils.Replace(
		input,
		`(const [\w_]+=([\w_]+)\.default\.get\([\w_]+\);)`,
		`${1}Spot.LocalStorage=${2}.default;`+showNotification)

	// encontra os símbolos player (playercosmossymbols[0]) e cosmos api (playercosmossymbols[1])
	playerCosmosSymbols := findSymbol("player and cosmos in PlayerHelper", input, []string{
		`this\._player=new ([\w_]+)\(([\w_]+)\.resolver,"spotify:app:zlink"`,
		`return new ([\w_]+)\(([\w_]+)\.resolver,"spotify:app:zlink","zlink"`,
	})

	if playerCosmosSymbols != nil {
		// inscreve-se na fila e defina os dados para spot.queue
		input = utils.Replace(
			input,
			`([\w_]+.prototype._player=null)`,
			fmt.Sprintf(
				`;new %s(%s.resolver,"spotify:internal:queue","queue","1.0.0").subscribeToQueue((e,r)=>{if(e){console.log(e);return;}Spot.Queue=r.getJSONBody();});${1}`,
				playerCosmosSymbols[0],
				playerCosmosSymbols[1]))
	}

	// vazamento dos métodos addtoqueue e removefromqueue
	input = utils.Replace(
		input,
		`(const [\w_]+=function\([\w_]+,[\w_]+\)\{)this\._bridge`,
		"${1}"+spotQueueJS+"this._bridge")

	// registra o evento de mudança de estado de reprodução/pausa
	input = utils.Replace(
		input,
		`this\.playing\([\w_]+\.is_playing&&![\w_]+\.is_paused\).+?;`,
		`${0}(this.playing()!==this._isPlaying)&&(this._isPlaying=this.playing(),Spot.Player.dispatchEvent&&Spot.Player.dispatchEvent(new Event("onplaypause")));`)

	// registra evento de mudança de música
	input = utils.Replace(
		input,
		`this\._uri=[\w_]+\.uri,this\._trackMetadata=[\w_]+\.metadata`,
		`${0},Spot.Player.dispatchEvent&&Spot.Player.dispatchEvent(new Event("songchange"))`)

	// vazamento de playbackcontrol para spot.playbackcontrol
	input = utils.Replace(
		input,
		`,(([\w_]+)\.playFromPlaylistResolver=)`,
		`;Spot.PlaybackControl = ${2};${1}`)

	// desativa a restrição de função de exposição
	input = utils.Replace(
		input,
		`(expose=function.+?)[\w_]+\.__spotify&&[\w_]+\.__spotify\.developer_mode&&`,
		"${1}")

	return input
}

const spotPlayerJS = `
Spot.Player.seek=(p)=>{if(p<=1)p=Math.round(p*(Spot.Player.data?Spot.Player.data.track.metadata.duration:0));this.seek(p)};
Spot.Player.getProgressMs=()=>this.progressbar.getRealValue();
Spot.Player.getProgressPercent=()=>this.progressbar.getPercentage();
Spot.Player.getDuration=()=>this.progressbar.getMaxValue();
Spot.Player.skipForward=(a=15e3)=>Spot.Player.seek(Spot.Player.getProgressMs()+a);
Spot.Player.skipBack=(a=15e3)=>Spot.Player.seek(Spot.Player.getProgressMs()-a);
Spot.Player.setVolume=(v)=>this.changeVolume(v, false);
Spot.Player.increaseVolume=()=>this.increaseVolume();
Spot.Player.decreaseVolume=()=>this.decreaseVolume();
Spot.Player.getVolume=()=>this.volumebar.getValue();
Spot.Player.next=()=>this._doSkipToNext();
Spot.Player.back=()=>this._doSkipToPrevious();
Spot.Player.togglePlay=()=>this._doTogglePlay();
Spot.Player.play=()=>{!this.playing() && this._doTogglePlay();};
Spot.Player.pause=()=>{this.playing() && this._doTogglePlay();};
Spot.Player.isPlaying=()=>this.progressbar.isPlaying();
Spot.Player.toggleShuffle=()=>this.toggleShuffle();
Spot.Player.getShuffle=()=>this.shuffle();
Spot.Player.setShuffle=(b)=>{this.shuffle(b)};
Spot.Player.toggleRepeat=()=>this.toggleRepeat();
Spot.Player.getRepeat=()=>this.repeat();
Spot.Player.setRepeat=(r)=>{this.repeat(r)};
Spot.Player.getMute=()=>this.mute();
Spot.Player.toggleMute=()=>this._doToggleMute();
Spot.Player.setMute=(b)=>{this.volumeEnabled()&&this.changeVolume(this._unmutedVolume,b)};
Spot.Player.thumbUp=()=>this.thumbUp();
Spot.Player.getThumbUp=()=>this.trackThumbedUp();
Spot.Player.thumbDown=()=>this.thumbDown();
Spot.Player.getThumbDown=()=>this.trackThumbedDown();
Spot.Player.formatTime=(ms)=>this._formatTime(ms);
Spot.Player.eventListeners={};
Spot.Player.addEventListener= (type, callback) => {
	if (!(type in Spot.Player.eventListeners)) {
		Spot.Player.eventListeners[type] = [];
	}
	Spot.Player.eventListeners[type].push(callback)
};
Spot.Player.removeEventListener = (type, callback) => {
    if (!(type in Spot.Player.eventListeners)) {
        return;
    }
    var stack = Spot.Player.eventListeners[type];
    for (var i = 0, l = stack.length; i < l; i += 1) {
        if (stack[i] === callback) {
            stack.splice(i, 1);
            return;
        }
    }
};
Spot.Player.dispatchEvent = (event) => {
    if (!(event.type in Spot.Player.eventListeners)) {
        return true;
    }
    var stack = Spot.Player.eventListeners[event.type];
    for (var i = 0, l = stack.length; i < l; i += 1) {
        stack[i](event);
    }
    return !event.defaultPrevented;
};
`

const spotQueueJS = `
Spot.addToQueue = (uri,callback) => {
	uri = Spot.LibURI.from(uri);
	if (uri.type === Spot.LibURI.Type.ALBUM) {
		this.getAlbumTracks(uri, (err,tracks) => {
			if (err) {
				console.log("Spot.addToQueue", err);
				return;
			}

			this.queueTracks(tracks, callback)
		})
	} else if (uri.type === Spot.LibURI.Type.TRACK || uri.type === Spot.LibURI.Type.EPISODE) {
		this.queueTracks([uri], callback);
	} else {
		console.log("Spot.addToQueue: Only Track, Album, Episode URIs are accepted");
		return;
	}
};

Spot.removeFromQueue = (uri, callback) => {
    if (Spot.Queue) {
        let indices = [];
        const uriObj = Spot.LibURI.from(uri);
        if (uriObj.type === Spot.LibURI.Type.ALBUM) {
            this.getAlbumTracks(uriObj, (err, tracks) => {
                if (err) {
                    console.log(err);
                    return;
                }
                tracks.forEach((trackUri) => {
					Spot.Queue.next_tracks.forEach((nt, index) => {
						trackUri == nt.uri && indices.push(index);
					})
				})
            })
        } else if (uriObj.type === Spot.LibURI.Type.TRACK || uriObj.type === Spot.LibURI.Type.EPISODE) {
            Spot.Queue.next_tracks.forEach((track, index) => {
				track.uri == uri && indices.push(index)
			})
        } else {
			console.log("Spot.removeFromQueue: Only Album, Track and Episode URIs are accepted")
			return;
		}

        indices = indices.reduce((a, b) => {
            if (a.indexOf(b) < 0) {
                a.push(b)
            }
            return a
        }, []);
        this.removeTracksFromQueue(indices, callback)
    }
};
`