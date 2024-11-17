// @ts-check

// nome: auto skip video
// pula vídeos automaticamente

/// <reference path="../globals.d.ts" />

(function SkipVideo() {
    Spot.Player.addEventListener("songchange", () => {
        const meta = Spot.Player.data.track.metadata;
        // os anúncios também são um tipo de mídia de vídeo, então é preciso excluí-los
        
        if (
            meta["media.type"] === "video" &&
            meta.is_advertisement !== "true"
        ) {
            Spot.Player.next();
        }
    });
})();