import React from "react";
import RowSection from "components/row-section";
import "components/document-view.scss";

interface DocumentAttributeProps {
  field: string;
  value: any;
}

const DocumentAttribute: React.FC<DocumentAttributeProps> = (
  props: DocumentAttributeProps
) => {
  return (
    <div className="document-block-item">
      <RowSection style={{ flex: 1 }} inverted={true}>
        {props.field}
      </RowSection>
      <RowSection style={{ flex: 5, textAlign: "center" }} inverted={true}>
        {props.value}
      </RowSection>
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
      {props.description != "" && (
        <DocumentAttribute field="Description" value={props.description} />
      )}
      <DocumentAttribute field="Version" value={props.version} />
    </div>
  );
};

export default DocumentBlock;
