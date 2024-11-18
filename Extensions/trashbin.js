// @ts-check

// nome: trashbin
// lança músicas para a trashbin para nunca mais ouví-las novamente

/// <reference path="../globals.d.ts" />

(function TrashBin() {
    /**
     * por padrão, a lista de músicas inúteis é salva em spot.localstorage, mas
     * tudo será limpo se o spotify for desinstalado. então, em vez de coletar
     * as músicas inúteis novamente, você pode usar o serviço jsonbin para
     * armazenar sua lista, que é totalmente gratuito e rápido.
     * 
     * https://jsonbin.io/
     * 
     * e e crie um jsonbin com o objeto padrão:

{
    "trashSongList": {},
    "trashArtistList": {}
}

     * em clique em criar. após isso, ele irá gerar um url de acesso, clique em
     * copiar e cole-o na constante jsonbinurl abaixo. o url deve ficar assim:

        https://api.jsonbin.io/b/XXXXXXXXXXXXXXXXXXXX

     * salve esse arquivo, rode o comando "apply" no spot para fazer a alteração
     */

    const jsonBinURL = "";

    if (!Spot.Player.data || (!jsonBinURL && !Spot.LocalStorage)) {
        setTimeout(TrashBin, 1000);

        return;
    }

    let trashSongList = {};
    let trashArtistList = {};

    let userHitBack = false;

    /** @type {HTMLElement} */
    let trashIcon;

    let banSong = () => {};

    const THROW_TEXT = "jogar na lixeira";
    const UNTHROW_TEXT = "retirar da lixeira";

    function createTrashArtistButton() {
        const div = document.createElement("div");

        div.classList.add("glue-page-header__button", "throw-artist");

        const button = document.createElement("button");

        button.classList.add(
            "button",
            "button-icon-with-stroke",
            "spoticon-browse-active-16"
        );

        button.setAttribute("data-tooltip", THROW_TEXT);

        div.appendChild(button);

        return div;
    }

    function createTrashTrackButton() {
        const button = document.createElement("button");

        button.classList.add(
            "button",
            "button-icon-only",
            "spoticon-browse-active-16"
        );

        button.setAttribute("data-tooltip-text", THROW_TEXT);

        button.style.position = "absolute";
        button.style.right = "24px";
        button.style.top = "-6px";
        button.style.transform = "scaleX(0.75)";

        return button;
    }

    // busca as faixas de lixo armazenadas e lista de artistas
    if (jsonBinURL) {
        fetch(`${jsonBinURL}/latest`)
            .then((res) => res.json())
            .then((data) => {
                trashSongList = data["trashSongList"];
                trashArtistList = data["trashArtistList"];

                if (!trashSongList || !trashArtistList) {
                    trashSongList = trashSongList || {};
                    trashArtistList = trashArtistList || {};

                    putDataOnline();
                }
            })
            .catch(console.log);
    } else {
        trashSongList = JSON.parse(Spot.LocalStorage.get("TrashSongList")) || {};
        trashArtistList = JSON.parse(Spot.LocalStorage.get("TrashArtistList")) || {};

        putDataLocal();
    }

    trashIcon = createTrashTrackButton();

    trashIcon.onclick = () => {
        banSong();

        const uri = Spot.Player.data.track.uri;
        const isBanned = !trashSongList[uri];

        if (isBanned) {
            trashSongList[uri] = true;

            Spot.Player.next();
        } else {
            delete trashSongList[uri];
        }

        updateTrackIconState(isBanned);

        storeList();
    };

    document.querySelector(".track-text-item").appendChild(trashIcon);

    // rastreamento quando os usuários clicam no botão anterior
    // ao fazer isso, o usuário pode retornar à música lançada
    // para retirá-la da lixeira
    document.getElementById("player-button-previous").addEventListener("click", () => (userHitBack = true));

    updateIconPosition();
    updateTrackIconState(trashSongList[Spot.Player.data.track.uri]);

    Spot.Player.addEventListener("songchange", watchChange);

    Spot.Player.addEventListener("appchange", ({ data: data }) => {
        if (data.isEmbeddedApp === true)
            return;

        if (data.id !== "artist")
            return;

        if (data.container.contentDocument.querySelector(".throw-artist"))
            return;

        const headers = data.container.contentDocument.querySelectorAll(
            ".glue-page-header__buttons"
        );

        if (headers.length < 1)
            return;

        const uri = `spotify:artist:${data.uri.split(":")[3]}`;

        headers.forEach((h) => {
            const button = createTrashArtistButton();

            button.onclick = () => {
                const isBanned = !trashArtistList[uri];

                if (isBanned) {
                    trashArtistList[uri] = true;
                } else {
                    delete trashArtistList[uri];
                }

                storeList();

                updateArtistIconState(button, isBanned);
            };

            h.appendChild(button);

            updateArtistIconState(button, trashArtistList[uri] !== undefined);
        });
    });

    function watchChange() {
        const data = Spot.Player.data || Spot.Queue;

        if (!data)
            return;

        const isBanned = trashSongList[data.track.uri];

        updateIconPosition();
        updateTrackIconState(isBanned);

        if (userHitBack) {
            userHitBack = false;

            return;
        }

        if (isBanned) {
            Spot.Player.next();

            return;
        }

        let uriIndex = 0;
        let artistUri = data.track.metadata["artist_uri"];

        while (artistUri) {
            if (trashArtistList[artistUri]) {
                Spot.Player.next();

                return;
            }

            uriIndex++;

            artistUri = data.track.metadata["artist_uri:" + uriIndex];
        }
    }

    // altera a posição do ícone da lixeira com base no contexto da playlist
    //
    // em playlists normais, track-text-item tem um ícone e seu padding-left
    // é 32px, apenas o suficiente para um ícone. ao anexar a classe de dois
    // ícones, seu padding-left é expandido para 64px
    //
    // na playlist das descobertas da semana, track-text-item tem dois ícones:
    // coração e banimento.
    //
    // a funcionalidade de banimento é semelhante à nossa, então, em vez de
    // preencher aquela pequena zona com 3 ícones, eu escondo o botão banir
    // do spotify e o substituo pelo ícone de lixeira. no entanto, ainda ativo
    // o menu de contexto do botão banir sempre que o usuário clica no ícone
    // da lixeira

    function updateIconPosition() {
        const trackContainer = document.querySelector(".track-text-item");

        if (!trackContainer.classList.contains("two-icons")) {
            trackContainer.classList.add("two-icons");

            trashIcon.style.right = "24px";

            return;
        }

        /** @type {HTMLElement} */
        const banButton = document.querySelector(
            ".track-text-item .nowplaying-ban-button"
        );

        if (banButton.style.display !== "none") {
            banButton.style.visibility = "hidden";

            trashIcon.style.right = "0px";

            banSong = banButton.click.bind(banButton);
        } else {
            banSong = () => {};
        }
    }

    /**
     * @param {boolean} isBanned
     */
    function updateTrackIconState(isBanned) {
        if (
            Spot.Player.data.track.metadata["is_advertisement"] === "true"
        ) {
            trashIcon.setAttribute("disabled", "true");

            return;
        }

        trashIcon.removeAttribute("disabled");

        if (isBanned) {
            trashIcon.classList.add("active");
            trashIcon.setAttribute("data-tooltip-text", UNTHROW_TEXT);
        } else {
            trashIcon.classList.remove("active");
            trashIcon.setAttribute("data-tooltip-text", THROW_TEXT);
        }
    }

    /**
     * @param {HTMLElement} button
     * @param {boolean} isBanned
     */
    function updateArtistIconState(button, isBanned) {
        const inner = button.querySelector("button");
        
        if (isBanned) {
            inner.classList.add("contextmenu-active");
            inner.setAttribute("data-tooltip", UNTHROW_TEXT);
        } else {
            inner.classList.remove("contextmenu-active");
            inner.setAttribute("data-tooltip", THROW_TEXT);
        }
    }

    function storeList() {
        if (jsonBinURL) {
            putDataOnline();
        } else {
            putDataLocal();
        }
    }

    function putDataOnline() {
        fetch(`${jsonBinURL}`, {
            method: "PUT",

            headers: {
                "Content-Type": "application/json"
            },

            body: JSON.stringify({
                trashSongList,
                trashArtistList
            }),
        }).catch(console.log);
    }

    function putDataLocal() {
        Spot.LocalStorage.set("TrashSongList", JSON.stringify(trashSongList));
        Spot.LocalStorage.set("TrashArtistList", JSON.stringify(trashArtistList));
    }
})();