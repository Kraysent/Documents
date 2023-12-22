import BaseForm from "components/BaseForm";
import React, { ChangeEvent, useState } from "react";
import "routes/DocumentPage/CreateLinkForm.scss";

interface CreateLinkForm {
  onClose: React.MouseEventHandler<HTMLButtonElement>;
  onCreate: (expiry: number) => void;
}

const CreateLinkForm: React.FC<CreateLinkForm> = (props: CreateLinkForm) => {
  const [expiry, setExpiry] = useState(0);

  const handleInputChange = (event: ChangeEvent<HTMLInputElement>) => {
    setExpiry(+event.target.value);
  };

  const handleCreate = () => {
    props.onCreate(expiry);
  };

  return (
    <BaseForm title="Create link" onClose={props.onClose}>
      <span className="expiry-field-name">How long should the link be valid (in days)?</span>
      <input className="expiry-input-field-box" type="number" defaultValue={14} onChange={handleInputChange} />
      <button className="create-button" onClick={handleCreate}>
        Create
      </button>
    </BaseForm>
  );
};

export default CreateLinkForm;
