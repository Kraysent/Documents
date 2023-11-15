import React, { useEffect, useState } from 'react';
import './App.scss';
import './heading.scss';
import './document-block.scss';
import GoogleButton from 'react-google-button';

let apiHost: string;

function setDevEnv() {
  apiHost = "http://localhost:8080/api"
}

function setProdEnv() {
  apiHost = "https://docarchive.space/api"
}

const Heading: React.FC = () => {
  return <div className="heading">
    <div className="heading-box">
      <span className="leftheading">doc</span>
      <span className="rightheading">archive</span>
    </div>
  </div>
}

const TextBlock: React.FC = () => {
  return <div className="textblock">
    <div className="textblock-title">
      Where am I?
    </div>
    <div className="textblock-body">
      You are in the docarchive - document storage.
    </div>
  </div>
}

function loginRedirect() {
  let url = `${apiHost}/auth/google/login`
  console.log(`redirecting to ${url}`)

  window.location.href = url
}

const LoginSection: React.FC = () => {
  return <div className="login-section">
    <GoogleButton
      className="google-button"
      type="dark"
      label="Log in or register"
      onClick={() => { loginRedirect() }} />
  </div>
}

class Document {
  id: string
  document_type: string
  attributes: object

  constructor(id: string, document_type: string, attributes: object) {
    this.id = id
    this.document_type = document_type
    this.attributes = attributes
  }
}

class GetUserDocumentsResponse {
  documents: Document[]

  constructor(documents: Document[]) {
    this.documents = documents
  }
}

interface DocumentBlockProps {
  key: number
  document: Document
}

function DocumentBlock(props: DocumentBlockProps) {
  return <div className="document-block" key={props.key}>
    <div className="document-id-block">{props.document.id}</div>
    <div className="document-name-block">{props.document.document_type}</div>
    <div className="document-attributes-block">{JSON.stringify(props.document.attributes)}</div>
  </div>
}


const App: React.FC = () => {
  const [docs, setDocs] = useState<Document[]>([]);
  const [mode, setMode] = useState("noauth");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  if (process.env.NODE_ENV == "development") {
    setDevEnv()
  } else if (process.env.NODE_ENV == "production") {
    setProdEnv()
  }

  useEffect(() => {
    fetch(`${apiHost}/v1/user/documents`, { credentials: 'include' })
      .then((response) => {
        if (!response.ok) {
          throw new Error(
            `This is an HTTP error: The status is ${response.status}`
          );
        }

        return response.json()
      })
      .then((data) => {
        let responseData: GetUserDocumentsResponse = data.data

        console.log(JSON.stringify(docs))
        setLoading(false)
        setDocs(responseData.documents)
        setMode("auth")
      }).catch((err) => {
        setError(err)
        console.error(err)
      })
  }, [])

  return (
    <div className="App">
      <Heading />
      {
        mode == "noauth" &&
        <div>
          <TextBlock />
          <LoginSection />
        </div>
      }
      {mode == "auth" &&
        <div className='document-container'>
          {
            docs.map((doc, i) => {
              return <DocumentBlock key={i} document={doc} />
            })
          }
        </div>
      }
    </div>
  );
}

export default App;
