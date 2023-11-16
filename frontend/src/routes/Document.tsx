import React from "react";
import { useParams } from "react-router-dom";
import Heading from "./heading";

const DocumentsPage: React.FC = () => {
    const { documentID } = useParams()

    return (
        <div>
            <Heading />
            {documentID}
        </div>
    )
}

export default DocumentsPage
