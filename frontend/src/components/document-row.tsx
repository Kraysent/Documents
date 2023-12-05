import React from "react";
import RowSection from "components/row-section";
import "components/document-row.scss";

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
        className="document-block"
        onClick={() =>
          (window.location.href =
            this.props.host + "/document/" + this.props.document.id)
        }
      >
        <RowSection text={this.props.document.name} alignment="center" />
        <RowSection
          flexSize={0.1}
          text={this.props.document.version.toString()}
          alignment="center"
        />
        {this.props.showDescription &&
          this.props.document.description != "" && (
            <RowSection
              flexSize={2}
              text={this.props.document.description}
              alignment="right"
            />
          )}
        {this.props.showDescription &&
          this.props.document.description == "" && (
            <RowSection
              flexSize={2}
              text="No description"
              active={false}
              alignment="right"
            />
          )}
      </div>
    );
  }
}

export default DocumentRow;
