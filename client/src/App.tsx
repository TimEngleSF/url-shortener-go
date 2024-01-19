import "./App.css";
import ShortenerForm from "./Components/ShortenerForm/ShortenerForm";
import { FaGithub, FaLinkedin } from "react-icons/fa6";

function App() {
  const queryParams = new URLSearchParams(window.location.search);
  const error = queryParams.get("error");
  return (
    <div className="flex h-full flex-col justify-center">
      <div className="  p-8">
        <h1 className="mb-12 text-4xl font-bold tracking-wider text-blue-500">
          URL Shortener
        </h1>
        {error && (
          <p className="mx-auto mb-4 max-w-80 rounded-md bg-red-100 py-2 text-red-500">
            {error}
          </p>
        )}
        <ShortenerForm />
      </div>
      <footer className="mt-auto flex items-center justify-center gap-5 bg-blue-500 py-4 text-white">
        <div className="flex gap-4">
          <a
            href="https://github.com/TimEngleSF/url-shortener-go"
            target="_blank"
          >
            <FaGithub className="h-6 w-6" />
          </a>
          <a href="https://linkedin.com/in/tim-engle" target="_blank">
            <FaLinkedin className="h-6 w-6" />
          </a>
        </div>
      </footer>
    </div>
  );
}

export default App;
