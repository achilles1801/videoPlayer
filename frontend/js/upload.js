// upload.js
document.addEventListener('DOMContentLoaded', function() {
    const uploadForm = document.getElementById('upload-form');
    const fileInput = document.getElementById('file-input');
    const videoPlayer = document.getElementById('video-player');

    uploadForm.addEventListener('submit', function(e) {
        e.preventDefault();
        const file = fileInput.files[0];
        if (file) {
            getPresignedUrlAndUpload(file);
        } else {
            alert('Please select a file to upload.');
        }
    });

    function getPresignedUrlAndUpload(file) {
        const presignUrlEndpoint = 'http://localhost:8080/presign?filename=upload/' + encodeURIComponent(file.name);    
        fetch(presignUrlEndpoint)
            .then(response => response.json())
            .then(data => {
                if (data.url) {
                    // Use the presigned URL obtained from the server
                    return uploadFileToS3(data.url, file).then(() => data.url);
                } else {
                    throw new Error('Failed to get a presigned URL.');
                }
            })
            .then(url => {
                // Update the video player source with the presigned URL
                updateVideoPlayer(url);
            })
            .catch(error => {
                console.error('An error occurred!', error);
            });
    }
    

    function uploadFileToS3(presignedUrl, file) {
        return fetch(presignedUrl, {
            method: 'PUT',
            body: file,
            headers: new Headers({
                'Content-Type': file.type
            })
        });
    }

    function updateVideoPlayer(videoUrl) {
        videoPlayer.src = videoUrl;
        videoPlayer.load();
        videoPlayer.play();
    }
});
