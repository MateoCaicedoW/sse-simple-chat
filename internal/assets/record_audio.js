function Record() {
    let audioContext;
    let recorder;
    let audioChunks = [];
    let audioBlob;
    let audioUrl;
    let analyser;
    let source;
    
    const recordButton = document.getElementById('recordButton');
    const stopButton = document.getElementById('stopButton');
    const visualizer = document.getElementById('visualizer');
    const visualizerContainer = document.getElementById("visualizer-container")
    const message = document.getElementById("message")
    const canvasCtx = visualizer.getContext('2d');
    recordButton.addEventListener('click', startRecording);
    stopButton.addEventListener('click', stopRecording);
    
    async function startRecording() {
        message.classList.add("hidden")
        visualizerContainer.classList.remove("hidden")
        stopButton.classList.remove("hidden")
        recordButton.classList.add("hidden")
        audioChunks = [];
        audioContext = new (window.AudioContext || window.webkitAudioContext)();
        analyser = audioContext.createAnalyser();
        
        try {
            const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
            recorder = new MediaRecorder(stream);
            source = audioContext.createMediaStreamSource(stream);
            source.connect(analyser);
    
            recorder.addEventListener('dataavailable', event => {
                audioChunks.push(event.data);
            });
    
            recorder.addEventListener('stop', () => {
                audioBlob = new Blob(audioChunks, { type: 'audio/wav' });
                audioUrl = URL.createObjectURL(audioBlob);
                const formData = new FormData();
                formData.append('audio', audioBlob);
    
                fetch('/audio', {
                    method: 'POST',
                    body: formData
                }).then(response => response.text())
                .then(result => {
                    visualizerContainer.classList.add("hidden")
                    message.classList.remove("hidden")
                });
            });
    
            recorder.start();
            visualize();
        } catch (err) {
            console.error('Error accessing microphone:', err);
        }
    }
    
    function stopRecording() {
        recorder.stop();
        recordButton.classList.remove("hidden")
        stopButton.classList.add("hidden")
        
        // Stop all tracks in the MediaStream to release the microphone
        const tracks = source.mediaStream.getTracks();
        tracks.forEach(track => track.stop());
    }
    
    function visualize() {
        const WIDTH = visualizer.width;
        const HEIGHT = visualizer.height;
    
        analyser.fftSize = 2048;
        const bufferLength = analyser.frequencyBinCount;
        const dataArray = new Uint8Array(bufferLength);
    
        canvasCtx.clearRect(0, 0, WIDTH, HEIGHT);
    
        function draw() {
            const drawVisual = requestAnimationFrame(draw);
    
            analyser.getByteTimeDomainData(dataArray);
    
            canvasCtx.fillStyle = 'rgb(200, 200, 200)';
            canvasCtx.fillRect(0, 0, WIDTH, HEIGHT);
    
            canvasCtx.lineWidth = 2;
            canvasCtx.strokeStyle = 'rgb(0, 0, 0)';
    
            canvasCtx.beginPath();
    
            const sliceWidth = WIDTH * 1.0 / bufferLength;
            let x = 0;
    
            for (let i = 0; i < bufferLength; i++) {
                const v = dataArray[i] / 128.0;
                const y = v * HEIGHT / 2;
    
                if (i === 0) {
                    canvasCtx.moveTo(x, y);
                } else {
                    canvasCtx.lineTo(x, y);
                }
    
                x += sliceWidth;
            }
    
            canvasCtx.lineTo(WIDTH, HEIGHT / 2);
            canvasCtx.stroke();
        }
    
        draw();
    }
}

export default Record

