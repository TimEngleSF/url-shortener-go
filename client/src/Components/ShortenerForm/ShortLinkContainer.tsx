import { useState } from "react";

interface ShortLinkProps {
  link: string;
  active: boolean;
}

const ShortLinkContainer = ({ link, active }: ShortLinkProps) => {
  const [displayCopied, setDisplayCopied] = useState(false);
  const borderColor = active ? "border-green-500" : "border-gray-400";
  const buttonClases = active
    ? "bg-green-500"
    : "bg-gray-400 cursor-not-allowed";

  const handleCopyClick = () => {
    // Commented out code will allow copy to work on iphone when hosted on a site without SSL cert

    // const input = document.createElement("textarea");
    // input.value = link;
    // document.body.appendChild(input);
    // input.select();
    // document.execCommand("copy");
    // document.body.removeChild(input);
    setDisplayCopied(true);

    navigator.clipboard.writeText(link);
    setTimeout(() => {
      setDisplayCopied(false);
    }, 2000);
  };
  return (
    <>
      <div className="flex drop-shadow-lg ">
        <input
          type="text"
          className={`flex-1 rounded-l-sm rounded-r-none border-[1px] border-solid ${borderColor} border-r-0 py-1 indent-2`}
          placeholder={link}
          disabled
          defaultValue={active ? link : ""}
          onClick={handleCopyClick}
        />
        <button
          className={`flex w-1/4 justify-center rounded-l-none rounded-r-sm py-1 text-white ${buttonClases}`}
          onClick={handleCopyClick}
          onTouchEnd={handleCopyClick}
          // disabled={!active}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={1.5}
            stroke="currentColor"
            className="h-6 w-6 text-white"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M15.666 3.888A2.25 2.25 0 0 0 13.5 2.25h-3c-1.03 0-1.9.693-2.166 1.638m7.332 0c.055.194.084.4.084.612v0a.75.75 0 0 1-.75.75H9a.75.75 0 0 1-.75-.75v0c0-.212.03-.418.084-.612m7.332 0c.646.049 1.288.11 1.927.184 1.1.128 1.907 1.077 1.907 2.185V19.5a2.25 2.25 0 0 1-2.25 2.25H6.75A2.25 2.25 0 0 1 4.5 19.5V6.257c0-1.108.806-2.057 1.907-2.185a48.208 48.208 0 0 1 1.927-.184"
            />
          </svg>
        </button>
      </div>
      <p
        className={`font-bold text-green-500 ${displayCopied ? "opacity-100" : "opacity-0"} duration-300`}
      >
        Link Copied!
      </p>
    </>
  );
};

export default ShortLinkContainer;
