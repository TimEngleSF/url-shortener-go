import React, { useState } from "react";
import ShortLinkContainer from "./ShortLinkContainer";

const ShortenerForm = () => {
  const [linkReady, setLinkReady] = useState(false);
  const [shortURL, setShortURL] = useState("https://examp.le/aBCDEF");
  return (
    <div className="flex w-full flex-col items-center rounded-md px-3 py-2 drop-shadow-lg ">
      <div className="flex w-full flex-col gap-6 md:w-3/4">
        <div className="flex ">
          <input
            type="text"
            className="flex-1 rounded-l-sm border-[1px] border-solid  border-blue-500 py-1 indent-2"
            placeholder="Enter your URL here"
          />
          <button className="w-1/4 rounded-r-sm bg-blue-500 py-1 text-white">
            Get URL
          </button>
        </div>
        <ShortLinkContainer active={linkReady} link={shortURL} />
      </div>
    </div>
  );
};

export default ShortenerForm;
