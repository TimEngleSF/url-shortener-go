import "./App.css";
import ShortenerForm from "./Components/ShortenerForm/ShortenerForm";

function App() {
  const queryParams = new URLSearchParams(window.location.search);
  const error = queryParams.get("error");
  return (
    <>
      <h1 className="mb-12 text-4xl font-bold tracking-wider text-blue-500">
        URL Shortener
      </h1>
      {error && <p className="text-red-500">{error}</p>}
      <ShortenerForm />
    </>
  );
}

export default App;
