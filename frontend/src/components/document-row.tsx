import React from "react";
import RowSection from "components/row-section";
import "components/clickable-row.scss";

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

interface DocumentRowProps {
  key: number;
  host: string;
  document: Document;
  showDescription: boolean;
}

class DocumentRow extends React.Component<DocumentRowProps> {
  public static defaultProps = {
    showDescription: true,
  };

  render() {
    return (
      <div
        className="clickable-row"
        onClick={() =>
          (window.location.href =
            this.props.host + "/document/" + this.props.document.id)
        }
      >
        <RowSection style={{ textAlign: "left", flex: 1 }}>
          {this.props.document.name}
        </RowSection>
        <RowSection style={{ textAlign: "center", flex: 0.1 }}>
          {this.props.document.version.toString()}
        </RowSection>
        {this.props.showDescription &&
          this.props.document.description != "" && (
            <RowSection style={{ textAlign: "right", flex: 2 }}>
              {this.props.document.description}
            </RowSection>
          )}
        {this.props.showDescription &&
          this.props.document.description == "" && (
            <RowSection active={false} style={{ textAlign: "right", flex: 2 }}>
              No description
            </RowSection>
          )}
      </div>
    );
  }
}

export default DocumentRow;
