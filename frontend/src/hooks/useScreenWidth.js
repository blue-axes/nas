import { useState, useEffect, useCallback, useRef } from "react";

export default function useScreenWidth() {
  const [width, setWidth] = useState(window.innerWidth);
  const timer = useRef(null);

  const handler = useCallback(() => {
    if (timer.current) clearTimeout(timer.current);
    timer.current = setTimeout(() => setWidth(window.innerWidth), 80);
  }, []);

  useEffect(() => {
    window.addEventListener("resize", handler);
    return () => {
      window.removeEventListener("resize", handler);
      if (timer.current) clearTimeout(timer.current);
    };
  }, [handler]);

  return width;
}