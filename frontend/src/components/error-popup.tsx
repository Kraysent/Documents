import React, { useState } from "react";
import "components/error-popup.scss";
import BackendError from "interactions/backend/error";

interface ErrorPopupProps {
  error: BackendError;
}

const ErrorPopup: React.FC<ErrorPopupProps> = (props: ErrorPopupProps) => {
  const [trigger, setTrigger] = useState(true);

  return trigger ? (
    <div className="error-popup">
      <div className="header">
        <h4 className="error-title" style={{ margin: 0 }}>
          {props.error.code}
        </h4>
        <button
          className="close-button"
          onClick={() => {
            setTrigger(false);
          }}
        >
          <span className="close-text">X</span>
        </button>
      </div>
      <span>{props.error.message}</span>
    </div>
  ) : (
    <div></div>
  );
};

export default ErrorPopup;
