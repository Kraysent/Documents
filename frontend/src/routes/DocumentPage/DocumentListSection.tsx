import DocumentRow from "components/document-row";
import ErrorPopup from "components/error-popup";
import { createBackendClient } from "interactions/backend/root";
import React from "react";

interface DocumentsListSectionProps {
  host: string;
  apiHost: string;
}

const DocumentsListSection: React.FC<DocumentsListSectionProps> = (
  props: DocumentsListSectionProps
) => {
  let client = createBackendClient(props.apiHost);
  let [docs, mode, loading, error] = client.getDocumentsList();

  return (
    <div>
      {loading && <div>Loading....</div>}
      {error && <ErrorPopup error={error} />}
      {mode == "auth" && (
        <div>
          {docs.map((doc, i) => {
            return (
              <DocumentRow
                host={props.host}
                key={i}
                document={doc}
                showDescription={false}
              />
            );
          })}
        </div>
      )}
    </div>
  );
};

export default DocumentsListSection;
