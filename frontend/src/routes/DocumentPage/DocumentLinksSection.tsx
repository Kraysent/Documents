import ErrorPopup from "components/error-popup";
import { AddLinkRow, LinkRow } from "components/links-row";
import { createBackendClient } from "interactions/backend/root";
import React, { useState } from "react";
import CreateLinkForm from "routes/DocumentPage/CreateLinkForm";
import "routes/DocumentPage/DocumentLinksSection.scss";

interface LinksListSectionProps {
  host: string;
  apiHost: string;
  documentID: string;
}

const LinksListSection: React.FC<LinksListSectionProps> = (
  props: LinksListSectionProps
) => {
  let client = createBackendClient(props.apiHost);
  let [links, loading, error] = client.getLinksList(props.documentID);

  const [createLinkWindowOpen, setCreateLinkWindowOpen] = useState(false);

  const handleCreate = async (expiry: string) => {
    try {
      let response = await client.createLink({
        document_id: props.documentID,
        expiry_date: expiry,
      });

      window.location.reload();
    } catch (e) {
      console.log(JSON.stringify(e));
    }

    setCreateLinkWindowOpen(false);
  };

  return (
    <div>
      {loading && <div>Loading....</div>}
      {error && <ErrorPopup error={error} />}
      {createLinkWindowOpen && (
        <CreateLinkForm
          onClose={() => setCreateLinkWindowOpen(false)}
          onCreate={handleCreate}
        />
      )}
      {links != null && (
        <div className="links-list-container">
          <AddLinkRow onClick={() => setCreateLinkWindowOpen(true)} />
          {links.map((link, i) => {
            return <LinkRow host={props.host} key={i} link={link} />;
          })}
        </div>
      )}
    </div>
  );
};

export default LinksListSection;
