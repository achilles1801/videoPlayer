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
        // Replace with your server URL and add any needed query parameters
        const presignUrlEndpoint = 'http://localhost:8080/presign?filename=' + encodeURIComponent(file.name);

        fetch(presignUrlEndpoint)
            .then(response => response.json())
            .then(data => {
                if (data.url) {
                    return uploadFileToS3(data.url, file);
                } else {
                    throw new Error('Failed to get a presigned URL.');
                }
            })
            .then(() => {
                // Here you would update the video player source to the public URL or the presigned GET URL
                // If using presigned URLs for viewing, you'd need to implement generating those in your backend
                const videoUrl = 'https://your-s3-bucket.s3.amazonaws.com/' + encodeURIComponent(file.name);
                updateVideoPlayer(videoUrl);
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
