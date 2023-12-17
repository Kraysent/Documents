import React from "react";
import RowSection from "components/row-section";
import "components/clickable-row.scss";
import "components/links-row.scss";
import { ReactComponent as CopyToClipboardIcon } from "assets/copy.svg";

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
}

class LinkRow extends React.Component<LinkRowProps> {
  public static defaultProps = {};

  render() {
    let isExpired = new Date(this.props.link.expiry_date) < new Date();
    let status = isExpired ? "expired" : this.props.link.status;
    let sharedLink = `${this.props.host}/share/${this.props.link.id}`;
    return (
      <div
        className="clickable-row"
        onClick={() => {
          navigator.clipboard.writeText(sharedLink);
        }}
      >
        {!isExpired && <CopyToClipboardIcon className="copy-icon" />}
        <RowSection
          style={{ textAlign: "center", flex: 1, overflow: "hidden" }}
          active={!isExpired}
        >
          {this.props.link.id.substring(0, 8)}
        </RowSection>
        <RowSection
          style={{ textAlign: "center", flex: 1, overflow: "hidden" }}
          active={!isExpired}
        >
          {status}
        </RowSection>
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
