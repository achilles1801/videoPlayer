import React, { useState, useEffect } from 'react';

function App() {
  const [videos, setVideos] = useState([]);


  const handleFileChange = (event) => {
    const file = event.target.files[0];
    if (file) {
      uploadVideo(file);
    }
  };

  const uploadVideo = (file) => {
    fetch('http://localhost:8080/presign?filename=' + file.name)
      .then(response => response.json())
      .then(data => {
        const { url } = data;
        return fetch(url, {
          method: 'PUT',
          body: file,
        });
      })
      .then(response => {
        if (response.ok) {
          console.log('Video uploaded successfully');
          fetchVideos(); // Fetch the updated list of videos
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
    <div className="App">
      <h2>Video Streaming Site</h2>
      <div>
      <input type="file" accept="video/*" onChange={handleFileChange} />
    </div>
      {videos.map((videoUrl, index) => (
  <video key={index} width="700px" height="400px" controls>
    <source src={videoUrl} type="video/mp4" />
  </video>
))}
    </div>
  );
}

export default App;
