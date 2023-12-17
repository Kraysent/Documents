import React from "react";
import "components/BaseForm.scss";

interface BaseFormProps {
  title: string;
  children: React.ReactNode;
  onClose: React.MouseEventHandler<HTMLButtonElement>;
}

const BaseForm: React.FC<BaseFormProps> = (props: BaseFormProps) => {
  return (
    <div className="base-form">
      <div style={{ display: "flex" }}>
        <span className="form-title" style={{ flex: 1 }}>
          {props.title}
        </span>
        <button className="close-button" onClick={props.onClose}>
          <span style={{ flex: 1 }}>X</span>
        </button>
      </div>
      {props.children}
    </div>
  );
};

export default BaseForm;
