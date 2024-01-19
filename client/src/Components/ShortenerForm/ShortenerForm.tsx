import { ChangeEvent, useState, useRef } from "react";
import axios from "axios";
import ShortLinkContainer from "./ShortLinkContainer";

const ShortenerForm = () => {
  const [linkReady, setLinkReady] = useState(false);
  const [shortURL, setShortURL] = useState("https://l.timengle.dev/aBCDEF");
  const [urlInput, setUrlInput] = useState("");
  const [errHighlight, setErrHighlight] = useState(false);
  const [showErr, setShowErr] = useState(false);
  const urlInputRef = useRef<HTMLInputElement | null>(null);

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (showErr) {
      setShowErr(false);
    }
    setUrlInput(e.target.value);
  };

  const handleUrlPost = async () => {
    if (!urlInput) {
      return;
    }
    try {
      const { data } = await axios({
        method: "POST",
        url: "/create",
        data: {
          siteUrl: urlInput,
        },
      });
      setLinkReady(true);
      setShortURL(data.link);
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (error: any) {
      if (error.response.status === 422) {
        setShowErr(true);
        setErrHighlight(true);
        setTimeout(() => {
          setErrHighlight(false);
        }, 1500);
        if (urlInputRef.current) {
          urlInputRef.current.select();
        }
      }
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      handleUrlPost();
    }
  };

  const handleInputValueFocus = (e: React.FocusEvent<HTMLInputElement>) => {
    e.target.select();
  };

  return (
    <div className="flex  flex-col items-center rounded-md px-3 py-2">
      <div className="flex w-full max-w-96 flex-col gap-6">
        <div className="flex drop-shadow-lg">
          <input
            type="text"
            className={`flex-1 rounded-l-sm rounded-r-none border-[1px] border-r-0 border-solid ${errHighlight ? "border-red-500" : "border-blue-500"} py-1 indent-2 transition-colors duration-300 focus:outline-none`}
            placeholder="Enter your URL here"
            onChange={handleChange}
            onKeyDown={handleKeyDown}
            onFocus={handleInputValueFocus}
            value={urlInput}
            ref={urlInputRef}
          />
          <button
            className={`w-1/4 rounded-r-sm ${errHighlight ? "bg-red-500" : "bg-blue-500"} py-1 text-white transition-colors duration-300`}
            onClick={handleUrlPost}
          >
            Get URL
          </button>
        </div>
        <ShortLinkContainer active={linkReady} link={shortURL} />
        <div className="drop-shadow-none">
          <p
            className={`font-bold text-red-500 ${showErr ? "opacity-100" : "opacity-0"} drop-shadow-none duration-300`}
          >
            Invalid URL format
          </p>
          <p
            className={`${showErr ? "opacity-100" : "opacity-0"} drop-shadow-none duration-300`}
          >
            Include <span className="font-bold">"https://"</span> or{" "}
            <span className="font-bold">"http://"</span>
          </p>
        </div>
      </div>
    </div>
  );
};

export default ShortenerForm;
