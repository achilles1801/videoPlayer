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
        fetchVideos();
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
    <div className="App bg-gray-100 min-h-screen">
      <header className="bg-gray-800 text-white text-center py-4">
        <h1 className="text-4xl font-bold">VideoPlayer</h1>
      </header>

      <div className="container mx-auto p-8">
        <div className="flex flex-wrap justify-center items-start gap-8">
          <div className="w-full lg:w-1/2">
            {selectedVideo && (
              <video key={selectedVideo} className="w-full aspect-video rounded-lg shadow-lg" controls>
                <source src={selectedVideo} />
              </video>
            )}
          </div>

          <div className="w-full lg:w-1/2">
            <div className="mb-4 text-center">
              <label className="block w-full lg:w-3/4 mx-auto text-lg py-2 px-4 bg-gray-200 text-gray-700 rounded-md cursor-pointer hover:bg-gray-300 transition duration-300 ease-in-out">
                <input type="file" accept="video/mp4,video/quicktime,video/3gpp,video/webm,video/ogg" onChange={handleFileChange} className="hidden" />
                Upload Video
              </label>
              <p className="text-sm text-gray-600 mt-2">Only MP4, MOV, 3GP, WebM, or Ogg video files.</p>
            </div>

            <div className="overflow-auto max-h-[550px] p-2">
              {videos.map((videoUrl, index) => (
                <div key={index} className="bg-white rounded-lg overflow-hidden shadow-md mb-4">
                  <video className="w-full h-40 object-contain" onClick={() => setSelectedVideo(videoUrl)}  playsInline>
                    <source src={videoUrl} />
                  </video>
                  <button onClick={(event) => handleDelete(event, videoUrl)}
                    className='w-full bg-red-700 hover:bg-red-900 text-white font-bold py-2 rounded-b-lg'>
                    Delete
                  </button>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      <footer className="bg-gray-800 text-white text-center py-4 mt-8">
        © 2024 My Video Stream by <a href="https://www.linkedin.com/in/majdkhawaldeh/" className="text-blue-500 hover:underline">Majd Khawaldeh</a> and <a href="https://www.linkedin.com/in/ayub-hunter/" className="text-blue-500 hover:underline">Ayub Hunter</a>
      </footer>
    </div>
  );
}

export default App;
