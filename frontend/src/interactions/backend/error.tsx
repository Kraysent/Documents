class BackendError extends Error {
  code: string;
  message: string;
}

export default BackendError;
