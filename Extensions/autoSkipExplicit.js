// @ts-check

// pula músicas explícitas automaticamente

/// <reference path="../globals.d.ts" />

(function ChristianSpotify() {
    if (!Spot.LocalStorage) {
        setTimeout(ChristianSpotify, 1000);

        return;
    }

    let isEnabled = Spot.LocalStorage.get("ChristianMode") === "1";

    new Spot.Menu.Item("Christian mode", isEnabled, (self) => {
        isEnabled = !isEnabled;

        Spot.LocalStorage.set("ChristianMode", isEnabled ? "1" : "0");

        self.setState(isEnabled);
    }).register();

    Spot.Player.addEventListener("songchange", () => {
        if (!isEnabled)
            return;

        const data = Spot.Player.data || Spot.Queue;

        if (!data)
            return;

        const isExplicit = data.track.metadata.is_explicit;
        
        if (isExplicit === "true") {
            Spot.Player.next();
        }
    });
})();