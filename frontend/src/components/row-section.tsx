import React from "react";
import "components/row-section.scss";

interface RowSectionProps {
  flexSize: number;
  text: string;
  alignment: string;
  active: boolean;
}

class RowSection extends React.Component<RowSectionProps> {
  public static defaultProps = {
    flexSize: 1,
    alignment: "left",
    active: true,
  };

  render() {
    const style = {
      flex: this.props.flexSize,
      textAlign: this.props.alignment,
    } as React.CSSProperties;

    let className = "row-section";

    if (!this.props.active) {
      className = "row-section-inactive";
    }

    return (
      <div className={className} style={style}>
        {this.props.text}
      </div>
    );
  }
}

export default RowSection;
