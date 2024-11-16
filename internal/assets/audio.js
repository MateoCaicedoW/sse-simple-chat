class AudioPlayer {
    constructor(container) {
        this.container = container;
        this.audio = container.querySelector('audio');
        this.playPauseBtn = container.querySelector('.play-pause-btn');
        this.progress = container.querySelector('.progress');
        this.timeDisplay = container.querySelector('.time');

        this.isPlaying = false;

        this.bindEvents();
    }

    bindEvents() {
        this.playPauseBtn.addEventListener('click', () => this.togglePlayPause());
        this.audio.addEventListener('timeupdate', () => this.updateProgress());
        this.audio.addEventListener('ended', () => this.onEnded());
    }

    togglePlayPause() {
        if (this.isPlaying) {
            this.pause();
        } else {
            this.play();
        }
    }

    play() {
        this.audio.play();
        this.playPauseBtn.innerHTML = '❚❚';
        this.playPauseBtn.setAttribute('aria-label', 'Pause audio');
        this.isPlaying = true;
    }

    pause() {
        this.audio.pause();
        this.playPauseBtn.innerHTML = '▶';
        this.playPauseBtn.setAttribute('aria-label', 'Play audio');
        this.isPlaying = false;
    }

    updateProgress() {
        const percent = (this.audio.currentTime / this.audio.duration) * 100;
        this.progress.style.width = `${percent}%`;
        this.timeDisplay.textContent = this.formatTime(this.audio.currentTime);
    }

    onEnded() {
        this.pause();
        this.progress.style.width = '0%';
        this.timeDisplay.textContent = this.formatTime(0);
    }

    formatTime(seconds) {
        const minutes = Math.floor(seconds / 60);
        const remainingSeconds = Math.floor(seconds % 60);
        return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
    }
}


export default AudioPlayer