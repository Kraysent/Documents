import React from "react";
import { Link } from "react-router-dom";

const Heading: React.FC = () => {
  return (
    <div className="heading">
      <Link to="/">
        <div className="heading-box">
          <span className="leftheading">doc</span>
          <span className="rightheading">archive</span>
        </div>
      </Link>
    </div>
  );
};

export default Heading;
