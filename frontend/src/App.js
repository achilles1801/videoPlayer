import './App.css';

function App() {
  return (
    <div className="App">
      <h2>Video Streaming Site</h2>
      <video width="700px" height="400px" controls > 
      <source src="https://d2ufs6yhgycudn.cloudfront.net/uploads/20231231225832-1e94e06e-a501-482a-a4d9-178f00c9c17a-834331ab4211531d0e81f82841f82555_seamless_loop.mp4" type = "video/mp4"/>
      </video>

    </div>
  );
}

export default App;
