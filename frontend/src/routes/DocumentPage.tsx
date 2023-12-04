import DocumentBlock from "components/document-view";
import DocumentRow from "components/document-row";
import Heading from "components/heading";
import { createBackendClient } from "interactions/backend/root";
import React from "react";
import { useParams } from "react-router-dom";
import "routes/document.scss";

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
      {error && <div>There was an error {JSON.stringify(error)}</div>}
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

interface DocumentViewPageProps {
  host: string;
  apiHost: string;
}

const DocumentViewPage: React.FC<DocumentViewPageProps> = (
  props: DocumentViewPageProps
) => {
  return (
    <div>
      <Heading />
      <div className="document-page-sections">
        <div style={{ flex: 1 }}>
          <DocumentsListSection host={props.host} apiHost={props.apiHost} />
        </div>
        <div style={{ flex: 3 }}>
          <DocumentContentSection apiHost={props.apiHost} />
        </div>
      </div>
    </div>
  );
};

export default DocumentViewPage;
