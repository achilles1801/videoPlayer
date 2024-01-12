import React, { useState, useEffect } from 'react';

function App() {
  const [videos, setVideos] = useState([]);
  const [selectedVideo, setSelectedVideo] = useState(null);


  const handleFileChange = (event) => {
    const file = event.target.files[0];
    if (file && file.type === "video/mp4") {
      uploadVideo(file);
    } else {
      alert("Only MP4 videos are allowed.");
    }
  };

  const uploadVideo = (file) => {
    fetch('http://localhost:8080/presign?filename=' + file.name)
      .then(response => response.json())
      .then(data => {
        const { url } = data;
        console.log('Uploading to', url);
        return fetch(url, {
          method: 'PUT',
          headers: {
            'Content-Type': 'video/mp4',
          },
          body: file,
        });
      })
      .then(response => {
        if (response.ok) {
          console.log('Video uploaded successfully');
          fetchVideos(); // Fetch the updated list of videos
          const uploadedVideoUrl = URL.createObjectURL(file);
          setSelectedVideo(uploadedVideoUrl); // Set the newly uploaded video as the selected video
        } else {
          console.error('Upload failed', response.statusText);
        }
      })
      .catch(error => {
        console.error('Error uploading video:', error);
        console.error('Error details:', error.message);
      });
  };
  
  
  const fetchVideos = () => {
    fetch('http://localhost:8080/list')
      .then(response => response.json())
      .then(data => setVideos(data))
      .catch(error => console.error('Error fetching videos:', error));
  };
  
  useEffect(() => {
    fetchVideos();
  }, []);
  

  return (
    <div className="App bg-gray-100 min-h-screen p-8">
      <h2 className="text-3xl font-bold mb-4 text-center">Video Streaming Site</h2>
      <div className="mb-4">
  <label className="block w-full text-lg py-2 px-4 border border-gray-300 rounded-md mb-2 cursor-pointer hover:bg-gray-200">
    <input type="file" accept="video/mp4" onChange={handleFileChange} className="hidden" />
    Upload MP4 Video
  </label>
  <p className="text-sm text-gray-600">Please upload only MP4 video files.</p>
</div>

      <div className="flex">
      <div className="w-2/3 pr-4 mb-4">
  {selectedVideo && (
    <video className="w-full max-w-2xl aspect-video" controls>
      <source src={selectedVideo} type="video/mp4" />
    </video>
  )}
</div>

<div className="w-1/3 overflow-auto p-2">
  {videos.map((videoUrl, index) => (
    <div key={index} className="bg-white rounded-lg overflow-hidden shadow-md mb-4 cursor-pointer" onClick={() => setSelectedVideo(videoUrl)}>
      <video className="w-full h-40">
        <source src={videoUrl} type="video/mp4" />
      </video>
    </div>
  ))}
</div>


      </div>
      <footer className="bg-gray-800 text-white text-center p-4 mt-8">
  © 2024 My Video Stream
</footer>

    </div>
  );
}

export default App;
