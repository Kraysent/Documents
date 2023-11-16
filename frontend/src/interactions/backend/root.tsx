import Document from "interactions/backend/entities";
import { useEffect, useState } from "react";

export interface BackendClient {
  getDocumentsList(): [Document[], string, boolean, any];
  getDocument(id: string): [Document, string, boolean, any];
}

class GetUserDocumentsResponse {
  documents: Document[];
}

class BackendClientImpl {
  host: string;

  constructor() {}

  getDocumentsList(): [Document[], string, boolean, any] {
    const [docs, setDocs] = useState<Document[]>([]);
    const [mode, setMode] = useState("noauth");
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
      fetch(`${this.host}/v1/user/documents`, { credentials: "include" })
        .then((response) => {
          if (!response.ok) {
            throw new Error(
              `This is an HTTP error: The status is ${response.status}`
            );
          }

          return response.json();
        })
        .then((data) => {
          let responseData: GetUserDocumentsResponse = data.data;

          setLoading(false);
          setDocs(responseData.documents);
          setMode("auth");
        })
        .catch((err) => {
          setError(err);
          setLoading(false);
          console.error(err);
        });
    }, []);

    return [docs, mode, loading, error];
  }

  getDocument(id: string): [Document, string, boolean, any] {
    const [doc, setDoc] = useState<Document>(new Document("", "", 1, ""));
    const [mode, setMode] = useState("noauth");
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
      fetch(`${this.host}/v1/document/id?id=${id}`, { credentials: "include" })
        .then((response) => {
          if (!response.ok) {
            throw new Error(
              `This is an HTTP error: The status is ${response.status}`
            );
          }

          return response.json();
        })
        .then((data) => {
          let docResponse: Document = data.data;

          setLoading(false);
          setDoc(docResponse);
          setMode("auth");
        })
        .catch((err) => {
          setError(err);
          setLoading(false);
          console.error(err);
        });
    }, []);

    return [doc, mode, loading, error];
  }
}

export function createBackendClient(host: string): BackendClient {
  let client = new BackendClientImpl();
  client.host = host;

  return client;
}
