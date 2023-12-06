import React from "react";
import "components/row-section.scss";

interface RowSectionProps {
  flexSize: number;
  text: string;
  alignment: string;
  active: boolean;
  inverted: boolean;
}

class RowSection extends React.Component<RowSectionProps> {
  public static defaultProps = {
    flexSize: 1,
    alignment: "left",
    active: true,
    inverted: false,
  };

  render() {
    const style = {
      flex: this.props.flexSize,
      textAlign: this.props.alignment,
    } as React.CSSProperties;

    let className = "row-section";

    if (this.props.inverted) {
      className += "-inverted";
    }

    if (!this.props.active) {
      className += "-inactive";
    }

    return (
      <div className={className} style={style}>
        {this.props.text}
      </div>
    );
  }
}

export default RowSection;
