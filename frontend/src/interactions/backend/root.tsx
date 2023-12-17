import {
  Document,
  Link,
  LinkRequest,
  LinkResponse,
} from "interactions/backend/entities";
import { useEffect, useState } from "react";
import BackendError from "interactions/backend/error";

export interface BackendClient {
  getDocumentsList(): [Document[], string, boolean, any];
  getDocument(id: string): [Document, string, boolean, any];
  getDocumentViaLink(id: string): [Document, boolean, any];
  getLinksList(documentID: string): [Link[], boolean, any];
  createLink(request: LinkRequest): Promise<LinkResponse>;
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
    const [error, setError] = useState<BackendError | null>(null);

    useEffect(() => {
      fetch(`${this.host}/v1/document/id?id=${id}`, { credentials: "include" })
        .then((response) => {
          return response.text();
        })
        .then((rawResp) => {
          console.log(rawResp);

          return JSON.parse(rawResp);
        })
        .then((data) => {
          if (data.data != undefined) {
            let docResponse: Document = data.data;

            setLoading(false);
            setDoc(docResponse);
            setMode("auth");
          } else if (data.code != undefined) {
            let errResponse: BackendError = data;

            setError(errResponse);
            setLoading(false);
          }
        })
        .catch((err) => {
          setError(err);
          setLoading(false);
          console.error(err);
        });
    }, []);

    return [doc, mode, loading, error];
  }

  getDocumentViaLink(id: string): [Document, boolean, any] {
    const [doc, setDoc] = useState<Document>(new Document("", "", 1, ""));
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<BackendError | null>(null);

    useEffect(() => {
      fetch(`${this.host}/v1/link?id=${id}`, { credentials: "include" })
        .then((response) => {
          return response.text();
        })
        .then((rawResp) => {
          console.log(rawResp);

          return JSON.parse(rawResp);
        })
        .then((data) => {
          if (data.data != undefined) {
            let docResponse: Document = data.data;

            setLoading(false);
            setDoc(docResponse);
          } else if (data.code != undefined) {
            let errResponse: BackendError = data;

            setError(errResponse);
            setLoading(false);
          }
        })
        .catch((err) => {
          setError(err);
          setLoading(false);
          console.error(err);
        });
    }, []);

    return [doc, loading, error];
  }

  getLinksList(documentID: string): [Link[], boolean, any] {
    const [links, setLinks] = useState<Link[] | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<BackendError | null>(null);

    useEffect(() => {
      fetch(`${this.host}/v1/links?document_id=${documentID}`, {
        credentials: "include",
      })
        .then((response) => {
          return response.text();
        })
        .then((rawResp) => {
          console.log(rawResp);

          return JSON.parse(rawResp);
        })
        .then((data) => {
          if (data.data != undefined) {
            let response: Link[] = data.data.links;

            setLoading(false);
            setLinks(response);
          } else if (data.code != undefined) {
            let errResponse: BackendError = data;

            setError(errResponse);
            setLoading(false);
          }
        })
        .catch((err) => {
          setError(err);
          setLoading(false);
          console.error(err);
        });
    }, []);

    return [links!!, loading, error];
  }

  async createLink(request: LinkRequest): Promise<LinkResponse> {
    return fetch(`${this.host}/v1/link`, {
      credentials: "include",
      method: "POST",
      body: JSON.stringify(request),
    })
      .then((response) => {
        return response.text();
      })
      .then((rawResp) => {
        console.log(rawResp);

        return JSON.parse(rawResp);
      })
      .then((data) => {
        if (data.code == undefined) {
          let response: LinkResponse = data.data;

          return response;
        } else {
          let errResponse: BackendError = data;

          throw errResponse;
        }
      })
      .catch((err) => {
        console.error(err);

        throw new Error(err);
      });
  }
}

export function createBackendClient(host: string): BackendClient {
  let client = new BackendClientImpl();
  client.host = host;

  return client;
}
