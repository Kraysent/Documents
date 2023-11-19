import React from "react";
import ReactDOM from "react-dom/client";
import "index.css";
import Root from "routes/Root";
import ErrorPage from "routes/error";
import reportWebVitals from "reportWebVitals";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import DocumentsPage from "routes/Document";

let apiHost!: string;
let host!: string;

function setDevEnv() {
  apiHost = "http://localhost:8080/api";
  host = "http://localhost:3000";
}

function setProdEnv() {
  apiHost = "https://docarchive.space/api";
  host = "https://docarchive.space";
}

if (process.env.NODE_ENV == "development") {
  setDevEnv();
} else if (process.env.NODE_ENV == "production") {
  setProdEnv();
}

const router = createBrowserRouter([
  {
    path: "/",
    element: <Root host={host} apiHost={apiHost} />,
    errorElement: <ErrorPage />,
  },
  {
    path: "document/:documentID",
    element: <DocumentsPage apiHost={apiHost} />,
  },
]);

const root = ReactDOM.createRoot(
  document.getElementById("root") as HTMLElement
);
root.render(<RouterProvider router={router} />);

reportWebVitals();
