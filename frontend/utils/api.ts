// utils/api.ts
import { BookInfo, BooksResponse } from "../types/book.ts";

export const API_BASE_URL = "http://localhost:8080";
export const TOKEN = "eqVVujtPwO0Xs6G0bX1Qjs3j3^G6$M";

export class BookAPI {
  static async getAllBooks(): Promise<BooksResponse> {
    try {
      const response = await fetch(`${API_BASE_URL}/books`, {
        headers: {
          "Authorization": `${TOKEN}`,
        },
      });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      return await response.json();
    } catch (error) {
      console.error("Error fetching all books:", error);
      return { books: null, count: 0 };
    }
  }

  static async getBookByISBN(isbn: string): Promise<BooksResponse> {
    try {
      const response = await fetch(`${API_BASE_URL}/book:${isbn}`, {
        headers: {
          "Authorization": `${TOKEN}`,
        },
      });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      return await response.json();
    } catch (error) {
      console.error("Error fetching book by ISBN:", error);
      return { books: null, count: 0 };
    }
  }

  static async searchBooksByTitle(title: string): Promise<BooksResponse> {
    try {
      const encodedTitle = encodeURIComponent(title);
      const response = await fetch(`${API_BASE_URL}/books/search:${encodedTitle}`, {
        headers: {
          "Authorization": `${TOKEN}`,
        },
      });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      return await response.json();
    } catch (error) {
      console.error("Error searching books:", error);
      return { books: null, count: 0 };
    }
  }
}