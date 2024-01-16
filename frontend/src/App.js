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
    const url = new URL('https://kcynhi0fr1.execute-api.us-east-2.amazonaws.com/new/presign');
    url.searchParams.append('filename', file.name);
    url.searchParams.append('contenttype', file.type);
    console.log('File type:', file.type);
    
    fetch(url, {mode: 'cors'})
      .then(response => {
        if (!response.ok) {
          console.error('Error fetching presigned URL:', response.status, response.statusText);
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
      })
      .then(data => {
        const { url } = data;
        console.log('Presigned URL:', url);
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
        if (!response.ok) {
          console.error('Error uploading file:', response.status, response.statusText);
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        console.log('Upload successful');
      })
      .catch(error => {
        console.error('Error uploading video:', error);
        console.error('Error details:', error.message);
      });
  };
  
  
  const fetchVideos = () => {
    const request = new Request('https://kcynhi0fr1.execute-api.us-east-2.amazonaws.com/new/list');
    fetch(request)
      .then(response => {
        if (!response.ok) {
          console.error('Error status:', response.status);
          return response.text().then(text => {
            console.error('Error body:', text);
            throw new Error('Network response was not ok');
          });
        }
        return response.json();
      })
      .then(data => {
        // Ensure that videos is always an array
        setVideos(Array.isArray(data) ? data : []);
      })
      .catch(error => {
        console.error('Error fetching videos:', error);
        // Set videos to an empty array in case of error
        setVideos([]);
      });
  };
  
  
  useEffect(() => {
    fetchVideos();
  }, []);


  const handleDelete = (event, videoUrl) => {
    event.stopPropagation(); // Prevent the click event from triggering the video selection
    fetch(`https://kcynhi0fr1.execute-api.us-east-2.amazonaws.com/new/delete?videoUrl=${encodeURIComponent(videoUrl)}`, {
      method: 'DELETE',
    })
      .then(response => {
        if (response.ok) {
          console.log('Video deleted successfully');
          // Optimistically update the state by filtering out the deleted video
          setVideos(videos => videos.filter(url => url !== videoUrl));
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
