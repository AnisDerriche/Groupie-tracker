document.querySelectorAll('/son').forEach(son => {
    const image = son.querySelector('img'); 
    const audio = son.querySelector('audio');

    image.addEventListener('click', () => {
        if (audio.paused) {
            // Arrêter les autres sons en cours
            document.querySelectorAll('audio').forEach(a => a.pause());
            audio.play();
        } else {
            audio.pause();
        }
    });
});

document.querySelectorAll('.gif').forEach(img => {
    img.addEventListener('click', () => {
        if (img.dataset.active === "true") {
            img.src = "son2.png"; // Remplace par l'image statique
            img.dataset.active = "false";
        } else {
            img.src = img.dataset.gif; // Remet le GIF
            img.dataset.active = "true";
        }
    });
});