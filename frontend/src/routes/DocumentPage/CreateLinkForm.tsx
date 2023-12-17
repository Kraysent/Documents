import BaseForm from "components/BaseForm";
import React, { ChangeEvent, useEffect, useState } from "react";
import "routes/DocumentPage/CreateLinkForm.scss";

interface CreateLinkForm {
  onClose: React.MouseEventHandler<HTMLButtonElement>;
  onCreate: (expiry: string) => void;
}

const CreateLinkForm: React.FC<CreateLinkForm> = (props: CreateLinkForm) => {
  const [expiry, setExpiry] = useState("");

  const handleInputChange = (event: ChangeEvent<HTMLInputElement>) => {
    setExpiry(event.target.value);
  };

  const handleCreate = () => {
    props.onCreate(expiry);
  };

  return (
    <BaseForm title="Create link" onClose={props.onClose}>
      <span className="expiry-field-name">Expiry date</span>
      <input className="expiry-input-field-box" onChange={handleInputChange} />
      <button className="create-button" onClick={handleCreate}>
        Create
      </button>
    </BaseForm>
  );
};

export default CreateLinkForm;
