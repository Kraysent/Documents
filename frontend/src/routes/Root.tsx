import DocumentRow from "components/document-row";
import ErrorPopup from "components/error-popup";
import Heading from "components/heading";
import "components/heading.scss";
import { createBackendClient } from "interactions/backend/root";
import React from "react";
import GoogleButton from "react-google-button";
import "routes/Root.scss";
import "routes/documents-list-block.scss";

interface RootProps {
  apiHost: string;
  host: string;
}

const TextBlock: React.FC = () => {
  return (
    <div className="textblock">
      <div className="textblock-title">Where am I?</div>
      <div className="textblock-body">
        You are in the docarchive - document storage.
      </div>
    </div>
  );
};

function loginRedirect(apiHost: string) {
  let url = `${apiHost}/auth/google/login`;
  console.log(`redirecting to ${url}`);

  window.location.href = url;
}

const LoginSection: React.FC<RootProps> = (props: RootProps) => {
  return (
    <div className="login-section">
      <GoogleButton
        className="google-button"
        type="dark"
        label="Log in or register"
        onClick={() => {
          loginRedirect(props.apiHost);
        }}
      />
    </div>
  );
};

const App: React.FC<RootProps> = (props: RootProps) => {
  let client = createBackendClient(props.apiHost);
  let [docs, mode, loading, error] = client.getDocumentsList();

  return (
    <div className="App">
      <Heading />
      {error && <ErrorPopup error={error} />}
      {mode == "noauth" && (
        <div>
          <TextBlock />
          <LoginSection host={props.host} apiHost={props.apiHost} />
        </div>
      )}
      {mode == "auth" && (
        <div className="document-container">
          {docs.map((doc, i) => {
            return <DocumentRow host={props.host} key={i} document={doc} />;
          })}
        </div>
      )}
    </div>
  );
};

export default App;
