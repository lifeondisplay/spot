// @ts-check

// nome: modo dj
// apenas o modo queue, escondendo todos os controles

/// <reference path="../globals.d.ts" />

(function DJMode() {
    if (!Spot.LocalStorage || !Spot.addToQueue || !Spot.URI) {
        setTimeout(DJMode, 1000);
        
        return;
    }

    /**
     * @type {{enabled: boolean; hideControls: boolean}}
     */
    let setting = JSON.parse(Spot.LocalStorage.get("DJMode"));

    if (!setting || typeof setting !== "object") {
        setting = {
            enabled: false,
            hideControls: false
        };

        Spot.LocalStorage.set("DJMode", JSON.stringify(setting));
    }

    const enableToggle = new Spot.Menu.Item(
        "Enabled",
        setting.enabled,
        () => {
            setting.enabled = !setting.enabled;
            Spot.LocalStorage.set("DJMode", JSON.stringify(setting));
            document.location.reload();
        } 
    );

    const hideToggle = new Spot.Menu.Item(
        "Hide controls",
        setting.enabled && setting.hideControls,
        () => {
            setting.hideControls = !setting.hideControls;
            showHideControl(setting.hideControls);

            hideToggle.setState(setting.hideControls);
            Spot.LocalStorage.set("DJMode", JSON.stringify(setting));
        }
    );

    new Spot.Menu.SubMenu("DJ Mode", [
        enableToggle,
        hideToggle
    ]).register();

    if (!setting.enabled) {
        // não fazer nada quando o modo dj estiver desligado
        return;
    }

    /** @type {HTMLElement} */
    const playerControl = document.querySelector(".player-controls-container");

    /** @type {HTMLElement} */
    const extraControl = document.querySelector(".extra-controls-container");

    /** @type {HTMLElement} */
    const nowPlayingAddButton = document.querySelector(
        ".view-player .nowplaying-add-button"
    );

    const IFRAME_HIDE_ELEMENT_LIST =
        [
            '[data-ta-id="card-button-play"]',
            '[data-ta-id="card-button-add"]',
            '[data-ta-id="card-button-context-menu"]',
            "div.glue-page-header__buttons",
            "th.tl-more",
            ".tl-cell.tl-more",
            "th.tl-save",
            ".tl-cell.tl-save",
            "th.tl-feedback",
            ".tl-cell.tl-feedback",
            "th.tl-more",
            ".tl-cell.tl-more",
        ].join(",") + "{display: none !important}";

    const EMBEDDED_HIDE_ELEMENT_LIST =
        [
            "div.Header__buttons",
            '[data-ta-id="play-button"]',
            '[data-ta-id="card-button-add"]',
            '[data-ta-id="card-button-context-menu"]',
            '[data-ta-id="ta-table-cell-add"]',
            '[data-ta-id="ta-table-cell-more"]',
            'th[aria-label=""]',
        ].join(",") + "{display: none !important}";

    /**
     * @param {boolean} hide
     */
    function showHideControl(hide) {
        if (hide) {
            playerControl.style.display = "none";
            extraControl.style.display = "none";
            nowPlayingAddButton.style.display = "none";
        } else {
            playerControl.style.display = "";
            extraControl.style.display = "";
            nowPlayingAddButton.style.display = "";
        }
    }

    showHideControl(setting.hideControls);

    let interval;

    Spot.Player.addEventListener("appchange", ({ data: data }) => {
        interval && clearInterval(interval);

        if (data.isEmbeddedApp === true) {
            interval = setInterval(() => applyEmbedded(data.container), 500);
        } else {
            interval = setInterval(() => applyIframe(data.container.contentDocument), 500);
        }
    });

    /**
     * @param {string} uri
     * @returns {boolean} se o uri de entrada é compatível
     */
    function isValidURI(uri) {
        const uriType = Spot.URI.fromString(uri).type;

        if (
            uriType === Spot.URI.Type.ALBUM ||
            uriType === Spot.URI.Type.TRACK ||
            uriType === Spot.URI.Type.EPISODE
        ) {
            return true;
        }

        return false;
    }

    /**
     * @param {string} uri
     */
    const clickFunc = (uri) => () =>
        Spot.addToQueue(uri).then(() => {
            Spot.BridgeAPI.request("track_metadata", [uri], (_, track) => {
                Spot.showNotification(
                    `${track.name} - ${track.artists[0].name} adicionado à lista de reprodução.`
                );
            });
    });

    /**
     * @param {Document} doc
     */
    function applyIframe(doc) {
        doc.querySelectorAll(
            ".tl-cell.tl-play, .tl-cell.tl-number, .tl-cell.tl-type"
        ).forEach((cell) => {
            const playButton = cell.querySelector("button");

            if (playButton.hasAttribute("djmode-injected")) {
                return;
            }

            const songURI = cell.parentElement.getAttribute("data-uri");

            if (!isValidURI(songURI)) {
                return;
            }

            // remove todas as intenções de interação padrão
            playButton.setAttribute("data-button", "");
            playButton.setAttribute("data-ta-id", "");
            playButton.setAttribute("data-interaction-target", "");
            playButton.setAttribute("data-interaction-intent", "");
            playButton.setAttribute("data-log-click", "");
            playButton.setAttribute("djmode-injected", "true");
            playButton.onclick = clickFunc(songURI);
        });
        if (setting.hideControls) {
            addCSS(doc, "IframeDJModeHideControl", IFRAME_HIDE_ELEMENT_LIST);
        } else {
            removeCSS(doc, "IframeDJModeHideControl");
        }
    }

    /**
     *
     * @param {HTMLElement} container
     */
    function applyEmbedded(container) {
        const cellList = container.querySelectorAll(".TableCellTrackNumber");

        cellList.forEach((/** @type {HTMLElement} */ cell) => {
            const songURI = cell.parentElement.getAttribute("data-ta-uri");

            if (!isValidURI(songURI)) {
                return;
            }

            cell.onmouseover = () => {
                const playButton = cell.querySelector(
                    ".TableCellTrackNumber__button-wrapper"
                );

                if (playButton.hasAttribute("djmode-injected")) {
                    return;
                }

                const newButton = document.createElement("button");

                newButton.setAttribute("type", "button");

                newButton.classList.add(
                    "button",
                    "button-icon-with-stroke",
                    "button-play"
                );

                playButton?.removeChild(playButton.querySelector("button"));
                playButton.appendChild(newButton);
                playButton.setAttribute("djmode-injected", "true");

                newButton.onclick = clickFunc(songURI);
            };
        });

        if (setting.hideControls) {
            addCSS(document, "EmbeddedDJModeHideControl", EMBEDDED_HIDE_ELEMENT_LIST);
        } else {
            removeCSS(document, "EmbeddedDJModeHideControl");
        }
    }

    /**
     * @param {Document} doc
     * @param {string} id
     * @param {string} text
     */
    function addCSS(doc, id, text) {
        if (!doc.querySelector("head #" + id)) {
            const style = document.createElement("style");

            style.id = id;
            style.innerText = text;

            doc.querySelector("head").append(style);
        }
    }

    /**
     * @param {Document} doc
     * @param {string} id
     */
    function removeCSS(doc, id) {
        const found = doc.querySelector("head #" + id);

        if (found) {
            found.remove();
        }
    }
})();