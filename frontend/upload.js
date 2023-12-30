// upload.js
document.getElementById('upload-form').addEventListener('submit', function(e) {
    e.preventDefault();
    const file = document.getElementById('file-input').files[0];
    if (file) {
        uploadFile(file);
    } else {
        alert('Please select a file to upload.');
    }
});

function uploadFile(file) {
    // Implement the logic to request a pre-signed URL and upload the file
    // ...
    // On successful upload:
    updateVideoPlayer(file.name);
}

function updateVideoPlayer(fileName) {
    // Assuming you get a public URL for the video after upload
    // This could be a pre-signed GET URL or a public URL if your S3 objects are public
    const videoPlayer = document.getElementById('video-player');
    videoPlayer.src = 'http://your-s3-bucket-url/' + encodeURIComponent(fileName);
    videoPlayer.load(); // Load the new video source
    videoPlayer.play(); // Play the video
}
