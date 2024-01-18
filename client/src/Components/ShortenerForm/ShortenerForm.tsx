import { ChangeEvent, useState } from "react";
import axios from "axios";
import ShortLinkContainer from "./ShortLinkContainer";

const ShortenerForm = () => {
  const [linkReady, setLinkReady] = useState(false);
  const [shortURL, setShortURL] = useState("https://examp.le/aBCDEF");
  const [urlInput, setUrlInput] = useState("");

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
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
    } catch (error) {
      console.log(error);
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
    <div className="flex w-full flex-col items-center rounded-md px-3 py-2 drop-shadow-lg ">
      <div className="flex w-full flex-col gap-6 md:w-3/4">
        <div className="flex ">
          <input
            type="text"
            className="flex-1 rounded-l-sm rounded-r-none border-[1px] border-r-0 border-solid border-blue-500 py-1 indent-2 focus:outline-none"
            placeholder="Enter your URL here"
            onChange={handleChange}
            onKeyDown={handleKeyDown}
            onFocus={handleInputValueFocus}
            value={urlInput}
          />
          <button
            className="w-1/4 rounded-r-sm bg-blue-500 py-1 text-white"
            onClick={handleUrlPost}
          >
            Get URL
          </button>
        </div>
        <ShortLinkContainer active={linkReady} link={shortURL} />
      </div>
    </div>
  );
};

export default ShortenerForm;
