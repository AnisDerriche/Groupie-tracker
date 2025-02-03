document.querySelectorAll('.son').forEach(son => {
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
