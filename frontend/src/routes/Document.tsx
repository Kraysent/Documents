import Heading from "components/heading";
import RowSection from "components/row";
import { createBackendClient } from "interactions/backend/root";
import React from "react";
import { useParams } from "react-router-dom";
import "routes/document.scss";

interface DocumentAttributeProps {
  field: string;
  value: any;
}

const DocumentAttribute: React.FC<DocumentAttributeProps> = (
  props: DocumentAttributeProps
) => {
  return (
    <div className="document-block-item">
      <RowSection flexSize={1} text={props.field} />
      <RowSection flexSize={5} text={props.value} alignment="center"/>
    </div>
  );
};

interface DocumentBlockProps {
  name: string;
  description: string;
  version: number;
}

const DocumentBlock: React.FC<DocumentBlockProps> = (
  props: DocumentBlockProps
) => {
  return (
    <div className="individual-document-block">
      <DocumentAttribute field="Name" value={props.name} />
      <DocumentAttribute field="Description" value={props.description} />
      <DocumentAttribute field="Version" value={props.version} />
    </div>
  );
};

interface DocumentsPageProps {
  apiHost: string;
}

const DocumentsPage: React.FC<DocumentsPageProps> = (
  props: DocumentsPageProps
) => {
  const { documentID } = useParams();

  if (documentID == undefined) {
    return <div>Lalala</div>;
  }

  let client = createBackendClient(props.apiHost);
  let [doc, mode, loading, error] = client.getDocument(documentID);

  return (
    <div>
      <Heading />
      {loading && <div>Loading....</div>}
      {error && <div>There was an error {error}</div>}
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
