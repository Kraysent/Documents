import Heading from "components/heading";
import React from "react";
import { useParams } from "react-router-dom";
import "routes/DocumentPage/index.scss";
import DocumentContentSection from "routes/DocumentPage/DocumentContentSection";
import DocumentsListSection from "routes/DocumentPage/DocumentListSection";
import LinksListSection from "routes/DocumentPage//DocumentLinksSection";

interface DocumentViewPageProps {
  host: string;
  apiHost: string;
}

const DocumentViewPage: React.FC<DocumentViewPageProps> = (
  props: DocumentViewPageProps
) => {
  const { documentID } = useParams();

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
        <div style={{ flex: 1, overflow: "scroll" }}>
          <LinksListSection
            host={props.host}
            apiHost={props.apiHost}
            documentID={documentID!!}
          />
        </div>
      </div>
    </div>
  );
};

export default DocumentViewPage;
