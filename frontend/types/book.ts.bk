// types/book.ts
export interface BookInfo {
  ISBN: string;
  Title: string;
  Authors: string[];
  Description: string;
  Publishdate: string;
  Language: number;
  Image: {
    Source: {
      Scheme: string;
      Opaque: string;
      User: null;
      Host: string;
      Path: string;
      RawPath: string;
      OmitHost: boolean;
      ForceQuery: boolean;
      RawQuery: string;
      Fragment: string;
      RawFragment: string;
    };
    Path: string;
  };
}

export interface BooksResponse {
  books: BookInfo[] | null;
  count: number;
}

export const LanguageMap = {
  0: "Unknown",
  1: "Japanese",
  2: "English",
} as const;
