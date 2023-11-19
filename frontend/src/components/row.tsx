import React from "react";
import "components/row-section.scss"

interface RowSectionProps {
  flexSize: number;
  text: string;
  alignment: string;
}

class RowSection extends React.Component<RowSectionProps> {
  public static defaultProps = {
    flexSize: 1,
    alignment: "left",
  };

  render() {
    const style = {
      flex: this.props.flexSize,
      "text-align": this.props.alignment,
    } as React.CSSProperties;

    return (
      <div className="row-section" style={style}>
        {this.props.text}
      </div>
    );
  }
}

export default RowSection;
