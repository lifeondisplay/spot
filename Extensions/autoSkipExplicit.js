// @ts-check

// pula músicas explícitas automaticamente.
// toggle no menu de perfil.

/// <reference path="../globals.d.ts" />

(function ChristianSpotify() {
    if (!Spot.LocalStorage) {
        setTimeout(ChristianSpotify, 200);

        return;
    }

    const BUTTON_TEXT = "modo puro";

    let isEnabled = Spot.LocalStorage.get("ChristianMode") === "1";

    const item = document.createElement("div");
    item.classList.add("MenuItem");

    if (isEnabled) {
        item.classList.add("MenuItemToggle--checked");
    }

    item.innerText = BUTTON_TEXT;

    item.onclick = () => {
        isEnabled = !isEnabled;

        Spot.LocalStorage.set("ChristianMode", isEnabled ? "1" : "0");

        if (isEnabled) {
            item.classList.add(
                "MenuItemToggle--checked",
                "MenuItem--is-active"
            );
        } else {
            item.classList.remove(
                "MenuItemToggle--checked",
                "MenuItem--is-active"
            );
        }
    };
    
    let menuEl = document.getElementById("PopoverMenu-container");

    // observando o menu de perfil
    let menuObserver = new MutationObserver(() => {
        const menuRoot = menuEl.querySelector(".Menu__root-items");
        
        if (menuRoot) {
            menuRoot.prepend(item);
        }
    });

    menuObserver.observe(menuEl, { childList: true });
    
    Spot.Player.addEventListener("songchange", () => {
        if (!Spot.Player.data)
            return;

        const isExplicit = isEnabled && Spot.Player.data.track.metadata.is_explicit === "true";
        
        if (isExplicit) {
            Spot.Player.next();
        }
    });
})