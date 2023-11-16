import Heading from "components/heading";
import React from "react";
import { useRouteError } from "react-router-dom";

const ErrorPage: React.FC = () => {
  const error: any = useRouteError();
  console.error(error);

  return (
    <div id="error-page">
      <Heading />
      <h1>404!</h1>
      <p>Seems that this is not the page you are looking for!</p>
      <p>
        <i>{error.statusText || error.message}</i>
      </p>
    </div>
  );
};

export default ErrorPage;
