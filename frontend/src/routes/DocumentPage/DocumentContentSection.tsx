import DocumentBlock from "components/document-view";
import ErrorPopup from "components/error-popup";
import { createBackendClient } from "interactions/backend/root";
import { useParams } from "react-router-dom";
import React from "react";

interface DocumentContentSectionProps {
  apiHost: string;
}

const DocumentContentSection: React.FC<DocumentContentSectionProps> = (
  props: DocumentContentSectionProps
) => {
  const { documentID } = useParams();

  if (documentID == undefined) {
    return <div>Document is undefined</div>;
  }

  let client = createBackendClient(props.apiHost);
  let [doc, mode, loading, error] = client.getDocument(documentID);

  return (
    <div>
      {loading && <div>Loading....</div>}
      {error && <ErrorPopup error={error} />}
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

export default DocumentContentSection;
