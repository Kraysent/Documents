import React from 'react';
import './App.scss';
import './heading.scss';
import GoogleButton from 'react-google-button';

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
  let host = "localhost"

  console.log(`running this process in ${process.env.NODE_ENV} mode`)

  if (process.env.NODE_ENV == "production") {
    host = "docarchive.space"
  }

  let url = `http://${host}:8080/auth/google/login`
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

const App: React.FC = () => {
  return (
    <div className="App">
      <Heading />
      <TextBlock />
      <LoginSection />
    </div>
  );
}

export default App;
