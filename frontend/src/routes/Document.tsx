import Heading from "components/heading";
import React from "react";
import { useParams } from "react-router-dom";

const DocumentsPage: React.FC = () => {
  const { documentID } = useParams();

  return (
    <div>
      <Heading />
      {documentID}
    </div>
  );
};

export default DocumentsPage;
