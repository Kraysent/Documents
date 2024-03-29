import DocumentBlock from "components/document-view";
import ErrorPopup from "components/error-popup";
import Heading from "components/heading";
import { createBackendClient } from "interactions/backend/root";
import React from "react";
import { useParams } from "react-router-dom";

interface SharedDocumentPageProps {
  apiHost: string;
}

const SharedDocumentPage: React.FC<SharedDocumentPageProps> = (
  props: SharedDocumentPageProps
) => {
  const { linkID } = useParams();

  if (linkID == undefined) {
    return <div>Link is undefined</div>;
  }

  let client = createBackendClient(props.apiHost);
  let [doc, loading, error] = client.getDocumentViaLink(linkID);

  return (
    <div>
      <Heading />
      {loading && <div>Loading....</div>}
      {error && <ErrorPopup error={error} />}
      {!error && (
        <DocumentBlock
          name={doc.name}
          description={doc.description}
          version={doc.version}
        />
      )}
    </div>
  );
};

export default SharedDocumentPage;
