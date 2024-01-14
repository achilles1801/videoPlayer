import React, { useState, useEffect } from 'react';

function App() {
  const [videos, setVideos] = useState([]);
  const [selectedVideo, setSelectedVideo] = useState(null);


  const handleFileChange = (event) => {
    const file = event.target.files[0];
    const allowedTypes = ["video/mp4", "video/quicktime", "video/3gpp", "video/webm","video/ogg"];
    if (file && allowedTypes.includes(file.type)) {
      uploadVideo(file);
    } else {
      alert("Only MP4, MOV, WebM, Ogg and 3GP videos are allowed.");
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
            'Content-Type': file.type,
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
  const handleDelete = (event, videoUrl) => {
    event.stopPropagation(); // Prevent the click event from triggering the video selection
  
    fetch('http://localhost:8080/delete?videoUrl=' + encodeURIComponent(videoUrl), {
      method: 'DELETE',
    })
      .then(response => {
        if (response.ok) {
          console.log('Video deleted successfully');
          fetchVideos(); // Fetch the updated list of videos
          if (selectedVideo === videoUrl) {
            setSelectedVideo(null); // Unselect the video if it's currently selected
          }
        } else {
          console.error('Delete failed', response.statusText);
        }
      })
      .catch(error => {
        console.error('Error deleting video:', error);
        console.error('Error details:', error.message);
      });
  };

  

  return (
    <div className="App bg-gray-100 min-h-screen p-8">
      <div className="flex flex-col md:flex-row justify-between items-center mb-8">
      <h2 className="text-3xl font-bold mb-4 md:mb-0 text-center flex items-center justify-center text-gray-700 tracking-wider uppercase">VideoPlayer</h2>      </div>
      <div className="flex flex-col md:flex-row justify-between mb-8">
        <div className="flex justify-center mb-8 md:mb-0 md:mr-8 flex-grow">
          {selectedVideo && (
            <video key={selectedVideo} className="w-1/2 aspect-video" controls>
              <source src={selectedVideo} />
            </video>
          )}
        </div>
  
        <div className="flex flex-col md:max-w-2xl">
          <div>
            <label className="block w-full text-lg py-2 px-4 border border-gray-300 rounded-md mb-2 cursor-pointer hover:bg-gray-200">
              <input type="file" accept="video/mp4,video/quicktime,video/3gpp,video/webm,video/ogg" onChange={handleFileChange} className="hidden" />
              Upload Video
            </label>
            <p className="text-sm text-gray-600">Please upload only MP4, MOV, 3GP, WebM, or Ogg video files.</p>
          </div>
  
          <div className="overflow-auto p-2 max-h-96">
            {videos.map((videoUrl, index) => (
              <div key={index} className="bg-white rounded-lg overflow-hidden shadow-md mb-4 cursor-pointer">
                <video className="w-full h-40 object-scale-down" onClick={() => setSelectedVideo(videoUrl)}>
                  <source src={videoUrl} />
                </video>
                <button onClick={(event) => handleDelete(event, videoUrl)}
                className='bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded mt-2'
                >Delete</button>
              </div>
            ))}
          </div>
        </div>
      </div>
  
      <footer className="bg-gray-800 text-white text-center p-4 mt-8">
        © 2024 My Video Stream
      </footer>
    </div>
  );
}

export default App;
