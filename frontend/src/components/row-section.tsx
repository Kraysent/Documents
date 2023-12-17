import React from "react";
import "components/row-section.scss";

interface RowSectionProps {
  active: boolean;
  inverted: boolean;
  style: React.CSSProperties;
  className: string;
  children: any;
}

class RowSection extends React.Component<RowSectionProps> {
  public static defaultProps = {
    active: true,
    inverted: false,
    style: {},
    className: "",
  };

  render() {
    let className = "row-section";

    if (this.props.inverted) {
      className += "-inverted";
    }

    if (!this.props.active) {
      className += "-inactive";
    }

    className += this.props.className;

    return (
      <div className={className} style={this.props.style}>
        {this.props.children}
      </div>
    );
  }
}

export default RowSection;
