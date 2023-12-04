import DocumentBlock from "components/document-block";
import Heading from "components/heading";
import { createBackendClient } from "interactions/backend/root";
import React from "react";
import { useParams } from "react-router-dom";
import "routes/document.scss";

interface DocumentsPageProps {
  apiHost: string;
}

const DocumentsPage: React.FC<DocumentsPageProps> = (
  props: DocumentsPageProps
) => {
  const { documentID } = useParams();

  if (documentID == undefined) {
    return <div>Document is undefined</div>;
  }

  let client = createBackendClient(props.apiHost);
  let [doc, mode, loading, error] = client.getDocument(documentID);

  return (
    <div>
      <Heading />
      {loading && <div>Loading....</div>}
      {error && <div>There was an error {JSON.stringify(error)}</div>}
      {mode == "auth" && (
        <DocumentBlock
          name={doc.name}
          description={doc.description}
          version={doc.version}
        />
      )}
    </div>
  );
};

export default DocumentsPage;
