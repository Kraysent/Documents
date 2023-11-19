import Heading from "components/heading";
import "components/heading.scss";
import RowSection from "components/row";
import { createBackendClient } from "interactions/backend/root";
import React from "react";
import GoogleButton from "react-google-button";
import { Link } from "react-router-dom";
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

class Document {
  id: string;
  name: string;
  version: number;
  description: string;

  constructor(id: string, name: string, version: number, description: string) {
    this.id = id;
    this.name = name;
    this.version = version;
    this.description = description;
  }
}

interface DocumentBlockProps {
  key: number;
  document: Document;
  host: string;
}

const DocumentBlock: React.FC<DocumentBlockProps> = (
  props: DocumentBlockProps
) => {
  return (
    <Link
      to={props.host + "/document/" + props.document.id}
      style={{ textDecoration: "none" }}
    >
      <div className="document-block" key={props.key}>
        <RowSection text={props.document.name} alignment="center" />
        <RowSection
          flexSize={0.1}
          text={props.document.version.toString()}
          alignment="center"
        />
        <RowSection
          flexSize={2}
          text={props.document.description}
          alignment="right"
        />
      </div>
    </Link>
  );
};

const App: React.FC<RootProps> = (props: RootProps) => {
  let client = createBackendClient(props.apiHost);
  let [docs, mode, loading, error] = client.getDocumentsList();

  return (
    <div className="App">
      <Heading />
      {mode == "noauth" && (
        <div>
          <TextBlock />
          <LoginSection host={props.host} apiHost={props.apiHost} />
        </div>
      )}
      {mode == "auth" && (
        <div className="document-container">
          {docs.map((doc, i) => {
            return <DocumentBlock host={props.host} key={i} document={doc} />;
          })}
        </div>
      )}
    </div>
  );
};

export default App;
