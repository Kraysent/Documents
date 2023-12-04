import React from "react";
import "components/row-section.scss";

interface RowSectionProps {
  flexSize: number;
  text: string;
  alignment: string;
  shadowedText: boolean;
}

class RowSection extends React.Component<RowSectionProps> {
  public static defaultProps = {
    flexSize: 1,
    alignment: "left",
    shadowedText: false,
  };

  render() {
    const style = {
      flex: this.props.flexSize,
      textAlign: this.props.alignment,
      color: this.props.shadowedText ? "grey" : "black",
    } as React.CSSProperties;

    return (
      <div className="row-section" style={style}>
        {this.props.text}
      </div>
    );
  }
}

export default RowSection;
