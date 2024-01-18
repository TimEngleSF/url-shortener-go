import "./App.css";
import ShortenerForm from "./Components/ShortenerForm/ShortenerForm";

function App() {
  return (
    <>
      <h1 className="mb-12 text-4xl font-bold tracking-wider text-blue-500">
        URL Shortener
      </h1>
      <ShortenerForm />
    </>
  );
}

export default App;
