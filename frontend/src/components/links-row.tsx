import React from "react";
import RowSection from "components/row-section";
import "components/clickable-row.scss";
import "components/links-row.scss";
import { ReactComponent as CopyToClipboardIcon } from "assets/copy.svg";
import { ReactComponent as DeleteIcon } from "assets/delete.svg";

class Link {
  id: string;
  document_id: string;
  creation_date: string;
  expiry_date: string;
  status: string;
}

interface LinkRowProps {
  key: number;
  host: string;
  link: Link;
  onDelete: (link: Link) => void;
}

class LinkRow extends React.Component<LinkRowProps> {
  public static defaultProps = {};

  render() {
    let isExpired = new Date(this.props.link.expiry_date) < new Date();
    let isDisabled = this.props.link.status == "disabled";
    let isActive = !isExpired && !isDisabled;
    let status = isExpired ? "expired" : this.props.link.status;
    let sharedLink = `${this.props.host}/share/${this.props.link.id}`;

    const handleDeleteClick = () => {
      return this.props.onDelete(this.props.link);
    };

    return (
      <div className="clickable-row">
        {isActive && (
          <CopyToClipboardIcon
            className="copy-icon"
            onClick={() => {
              navigator.clipboard.writeText(sharedLink);
            }}
          />
        )}
        <RowSection
          style={{ textAlign: "center", flex: 1, overflow: "hidden" }}
          active={isActive}
        >
          {this.props.link.id.substring(0, 8)}
        </RowSection>
        <RowSection
          style={{ textAlign: "center", flex: 1, overflow: "hidden" }}
          active={isActive}
        >
          {status}
        </RowSection>
        {isActive && (
          <DeleteIcon className="delete-icon" onClick={handleDeleteClick} />
        )}
      </div>
    );
  }
}

interface AddLinkRowProps {
  onClick: React.MouseEventHandler<HTMLDivElement>;
}

class AddLinkRow extends React.Component<AddLinkRowProps> {
  render() {
    return (
      <div className="clickable-row" onClick={this.props.onClick}>
        <RowSection style={{ textAlign: "center", flex: 1 }}>+</RowSection>
      </div>
    );
  }
}

export { LinkRow, AddLinkRow };
