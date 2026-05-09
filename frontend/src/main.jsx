import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import router from "./router.jsx";
import { BrowserRouter, RouterProvider } from "react-router";
import "./tech-theme.css";

createRoot(document.getElementById("root")).render(
  <StrictMode>
    <RouterProvider router={router}></RouterProvider>
  </StrictMode>
);
