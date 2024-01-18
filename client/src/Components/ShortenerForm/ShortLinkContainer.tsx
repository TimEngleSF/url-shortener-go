interface ShortLinkProps {
  link: string;
  active: boolean;
}

const ShortLinkContainer = ({ link, active }: ShortLinkProps) => {
  const borderColor = active ? "border-green-500" : "border-gray-400";
  const buttonClases = active
    ? "bg-green-500"
    : "bg-gray-400 cursor-not-allowed";

  return (
    <div className="flex">
      <input
        type="text"
        className={`flex-1 rounded-l-sm border-[1px] border-solid ${borderColor} py-1 indent-2`}
        placeholder={link}
        disabled
      />
      <button
        className={`flex w-1/4 justify-center rounded-r-sm py-1 text-white ${buttonClases}`}
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
            d="M13.19 8.688a4.5 4.5 0 0 1 1.242 7.244l-4.5 4.5a4.5 4.5 0 0 1-6.364-6.364l1.757-1.757m13.35-.622 1.757-1.757a4.5 4.5 0 0 0-6.364-6.364l-4.5 4.5a4.5 4.5 0 0 0 1.242 7.244"
          />
        </svg>
      </button>
    </div>
  );
};

export default ShortLinkContainer;
