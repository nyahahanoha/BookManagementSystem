export enum Language {
  UNKNOWN = 0,
  ENGLISH = 1,
  JAPANESE = 2,
}

export interface Book {
  isbn: string;
  title: string;
  authors: string[];
  description: string;
  publishdate: string;
  language: Language;
  imageurl: string;
}
