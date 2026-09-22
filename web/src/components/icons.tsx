type IconProps = { className?: string };

const base = {
  viewBox: "0 0 24 24",
  fill: "none",
  stroke: "currentColor",
  strokeWidth: 1.75,
  "aria-hidden": true,
} as const;

export function SearchIcon(props: IconProps) {
  return (
    <svg {...base} strokeLinecap="round" strokeLinejoin="round" {...props}>
      <circle cx="10.5" cy="10.5" r="6.5" />
      <line x1="20" y1="20" x2="15.5" y2="15.5" />
    </svg>
  );
}

export function GridIcon(props: IconProps) {
  return (
    <svg {...base} strokeLinecap="round" strokeLinejoin="round" {...props}>
      <rect x="3.5" y="3.5" width="7" height="7" rx="1.2" />
      <rect x="13.5" y="3.5" width="7" height="7" rx="1.2" />
      <rect x="3.5" y="13.5" width="7" height="7" rx="1.2" />
      <rect x="13.5" y="13.5" width="7" height="7" rx="1.2" />
    </svg>
  );
}

export function BookmarkIcon(props: IconProps) {
  return (
    <svg {...base} strokeLinecap="round" strokeLinejoin="round" {...props}>
      <path d="M6 4h12v16l-6-4-6 4V4Z" />
    </svg>
  );
}

export function BriefcaseIcon(props: IconProps) {
  return (
    <svg {...base} strokeLinecap="round" strokeLinejoin="round" {...props}>
      <rect x="3.5" y="7.5" width="17" height="12" rx="1.5" />
      <path d="M8.5 7.5V6a2 2 0 0 1 2-2h3a2 2 0 0 1 2 2v1.5" />
      <line x1="3.5" y1="13" x2="20.5" y2="13" />
    </svg>
  );
}

export function ListIcon(props: IconProps) {
  return (
    <svg {...base} strokeLinecap="round" strokeLinejoin="round" {...props}>
      <line x1="8" y1="6" x2="20.5" y2="6" />
      <line x1="8" y1="12" x2="20.5" y2="12" />
      <line x1="8" y1="18" x2="20.5" y2="18" />
      <circle cx="4" cy="6" r="1" />
      <circle cx="4" cy="12" r="1" />
      <circle cx="4" cy="18" r="1" />
    </svg>
  );
}

export function GearIcon(props: IconProps) {
  return (
    <svg {...base} strokeLinecap="round" strokeLinejoin="round" {...props}>
      <circle cx="12" cy="12" r="3.2" />
      <path d="M12 3.5v2.4M12 18.1v2.4M20.5 12h-2.4M5.9 12H3.5M17.8 6.2l-1.7 1.7M7.9 16.1l-1.7 1.7M17.8 17.8l-1.7-1.7M7.9 7.9 6.2 6.2" />
    </svg>
  );
}

export function StarIcon(props: IconProps) {
  return (
    <svg {...base} strokeLinejoin="round" {...props}>
      <path d="m12 4 2.4 5.1 5.6.7-4.1 3.9 1 5.6-4.9-2.7-4.9 2.7 1-5.6-4.1-3.9 5.6-.7Z" />
    </svg>
  );
}

export function XIcon(props: IconProps) {
  return (
    <svg {...base} strokeLinecap="round" {...props}>
      <line x1="6" y1="6" x2="18" y2="18" />
      <line x1="18" y1="6" x2="6" y2="18" />
    </svg>
  );
}
