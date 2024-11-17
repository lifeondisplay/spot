// @ts-check

// nome: shuffle+
// shuffle de verdade sem bias

/// <reference path="../globals.d.ts" />

(function ShufflePlus() {
    if (!Spot.CosmosAPI) {
        setTimeout(ShufflePlus, 1000);

        return;
    }

    // texto de notificação quando a fila é embaralhada com sucesso
    /** @param {number} count */
    const NOTIFICATION_TEXT = (count) => `embaralhou ${count} itens!`;

    // se a queue do shuffer deve ser exibida
    const showShuffleQueueButton = true;

    // botões de reprodução aleatória de contexto
    function createContextButton() {
        const b = document.createElement("button");

        b.classList.add("button", "button-green");
        b.innerText = "Shuffle Context";

        b.setAttribute(
            "data-tooltip",
            "detecta o contexto de jogo atual e embaralhe todos os seus itens."
        );

        b.onclick = () => {
            const contextURI = Spot.Player.data.context_uri;
            const uriObj = Spot.LibURI.fromString(contextURI);

            switch (uriObj.type) {
                case Spot.LibURI.Type.SHOW:
                    showShuffle(uriObj.getBase62Id());

                    break;

                case Spot.LibURI.Type.PLAYLIST:
                    playlistShuffle(contextURI);

                    break;

                case Spot.LibURI.Type.FOLDER:
                    folderShuffle(contextURI);

                    break;

                case Spot.LibURI.Type.ALBUM:
                    albumShuffle(contextURI);

                    break;

                case Spot.LibURI.Type.COLLECTION:
                    collectionShuffle();

                    break;

                default:
                    Spot.showNotification(
                        `tipo de uri de contexto não compatível: ${uriObj.type}`
                    );
            }
        };

        return b;
    }

    // botões aleatórios da fila
    function createQueueButton() {
        const b = document.createElement("button");

        b.classList.add("button", "button-green");
        b.innerText = "fila aleatória";

        b.setAttribute(
            "data-tooltip",
            "embaralha os primeiros 80 itens ou menos que estão visíveis na página fila. útil apenas para fila de contexto misto."
        );

        b.onclick = () => {
            /** @type {Array} */
            let replace = Spot.Queue.next_tracks;

            let delimiterIndex = replace.findIndex(
                (value) => value.uri === "spotify:delimiter"
            );

            if (delimiterIndex !== -1) {
                replace.splice(delimiterIndex);
            }

            setQueue(shuffle(replace));
        };

        return b;
    }

    Spot.Player.addEventListener("appchange", ({ data: data }) => {
        if (data.isEmbeddedApp === true)
            return;

        if (data.id !== "queue")
            return;

        const headers = data.container.contentDocument.querySelectorAll(
            ".glue-page-header__buttons"
        );

        for (const e of headers) {
            if (e.hasAttribute("shuffleplus")) {
                continue;
            }

            e.setAttribute("shuffleplus", "1");

            e.append(createContextButton());

            if (showShuffleQueueButton) {
                e.append(createQueueButton());
            }
        }
    });

    /**
     * @param {string} uri
     * @returns {Promise<Array<{ uri: string }>>}
     */
    function requestPlaylist(uri) {
        return new Promise((resolve, reject) => {
            Spot.BridgeAPI.cosmosJSON(
                {
                    method: "GET",

                    uri: `sp://core-playlist/v1/playlist/${uri}/rows`,

                    body: {
                        policy: {
                            link: true
                        }
                    }
                },

                (error, res) => {
                    if (error) {
                        reject(error);

                        return;
                    }

                    let replace = res.rows;

                    replace = replace.map((item) => ({
                        uri: item.link
                    }));

                    resolve(replace);
                }
            );
        });
    }

    /**
     * @param {string} uri
     */
    function playlistShuffle(uri) {
        requestPlaylist(uri)
            .then((res) => setQueue(shuffle(res)))
            .catch((error) => console.error("lista de reprodução aleatória:", error));
    }

    /**
     * @param {string} uri
     */
    function folderShuffle(uri) {
        Spot.BridgeAPI.cosmosJSON(
            {
                method: "GET",

                uri: `sp://core-playlist/v1/rootlist`,

                body: {
                    policy: {
                        folder: {
                            rows: true,
                            link: true
                        }
                    }
                }
            },

            (error, res) => {
                if (error) {
                    console.error("pasta shuffle:", error);

                    return;
                }

                const requestFolder = res.rows.filter(
                    (item) => item.link === uri
                );

                if (requestFolder === 0) {
                    console.error("pasta shuffle: não é possível encontrar a pasta");
                    
                    return;
                }

                const requestPlaylists = requestFolder[0].rows.map((item) =>
                    requestPlaylist(item.link)
                );

                Promise.all(requestPlaylists)
                    .then((playlists) => {
                        const trackList = [];

                        playlists.forEach((p) => {
                            trackList.push(...p);
                        });

                        setQueue(shuffle(trackList));
                    }).catch((error) => console.error("pasta shuffle:", error));
            }
        );
    }

    function collectionShuffle() {
        Spot.BridgeAPI.cosmosJSON(
            {
                method: "GET",

                uri: "sp://core-collection/unstable/@/list/tracks/all",

                body: {
                    policy: {
                        list: {
                            link: true
                        }
                    }
                }
            },

            (error, res) => {
                if (error) {
                    console.log("collectionShuffle", error);

                    return;
                }

                let replace = res.items;

                replace = replace.map((item) => ({
                    uri: item.link,
                }));

                setQueue(shuffle(replace));
            }
        );
    }

    /**
     * @param {string} uri
     */
    function albumShuffle(uri) {
        const arg = [uri, 0, -1];
        Spot.BridgeAPI.request(
            "album_tracks_snapshot",

            arg,

            (error, res) => {
                if (error) {
                    console.error("álbum aleatorizado: ", error);

                    return;
                }

                let replace = res.array;

                replace = replace.map((item) => ({
                    uri: item
                }));

                setQueue(shuffle(replace));
            }
        );
    }

    /**
     * @param {string} uriBase62
     */
    function showShuffle(uriBase62) {
        Spot.CosmosAPI.resolver.get(
            {
                url: `sp://core-show/unstable/show/${uriBase62}`
            },

            (error, res) => {
                if (error) {
                    console.error("mostrar shuffle:", error);

                    return;
                }

                let replace = res.getJSONBody().items;

                replace = replace.map((item) => ({
                    uri: item.link,
                }));

                setQueue(shuffle(replace));
            }
        );
    }

    /**
     * @param {Array<{ uri: string }>} array lista de itens para embaralhar
     * @returns {Array<{ uri: string }>} array embaralhado
     *
     * de: https://bost.ocks.org/mike/shuffle/
     */
    function shuffle(array) {
        let counter = array.length;

        // embora existam elementos na matriz
        while (counter > 0) {
            // escolhe um índice aleatório
            let index = Math.floor(Math.random() * counter);

            // diminui o contador em 1
            counter--;

            // e troca o último elemento por ele
            let temp = array[counter];

            array[counter] = array[index];
            array[index] = temp;
        }

        return array;
    }

    /**
     * @param {Array<{ uri: string }>} state
     */
    function setQueue(state) {
        const count = state.length;

        state.push({ uri: "spotify:delimiter" });

        const currentQueue = Spot.Queue;
        
        currentQueue.next_tracks = state;

        const stringified = JSON.stringify(currentQueue);

        state.length = 0; // array nivelado

        const request = new Spot.CosmosAPI.Request(
            "PUT",
            "sp://player/v2/main/queue",
            null,
            stringified
        );

        Spot.CosmosAPI.resolver.resolve(request, (error, _) => {
            if (error) {
                console.log(error);

                return;
            }

            Spot.showNotification(NOTIFICATION_TEXT(count));
        });
    }
})();