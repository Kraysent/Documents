export class Document {
  id: string;
  name: string;
  version: number;
  description: string;

  constructor(id: string, name: string, version: number, description: string) {
    this.id = id;
    this.name = name;
    this.version = version;
    this.description = description;
  }
}

export class Link {
  id: string;
  document_id: string;
  creation_date: string;
  expiry_date: string;
  status: string;
}

