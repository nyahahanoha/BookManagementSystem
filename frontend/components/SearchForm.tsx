import { useState, useEffect } from "preact/hooks";
//import { API_BASE_URL } from "../utils/api.ts";

interface SearchFormProps {
  onSearch: (query: string, type: 'isbn' | 'title') => void;
  loading: boolean;
}

export default function SearchForm({ onSearch, loading }: SearchFormProps) {
  const [scanning, setScanning] = useState(false);
  //const [scanLoading, setScanLoading] = useState(false);

  // 初期化時にlocalStorageから状態復元
  useEffect(() => {
    const saved = localStorage.getItem("scan-active");
    if (saved === "true") setScanning(true);
  }, []);

  // 状態変更時にlocalStorageへ保存
  useEffect(() => {
    localStorage.setItem("scan-active", scanning ? "true" : "false");
  }, [scanning]);

  const handleSubmit = (e: JSX.TargetedEvent<HTMLFormElement, Event>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const query = formData.get("query") as string;
    const searchType = formData.get("searchType") as 'isbn' | 'title';
    onSearch(query.trim(), searchType);
  };

  //const handleScanToggle = async () => {
  //  setScanLoading(true);
  //  try {
  //    const endpoint = scanning ? "scan:stop" : "scan:start";
  //    await fetch(`${API_BASE_URL}/${endpoint}`, { method: "POST" });
  //    setScanning(!scanning);
  //  } catch (err) {
  //    // エラー処理は必要に応じて
  //  } finally {
  //    setScanLoading(false);
  //  }
  //};

  return (
    <div class="searchform-blue-bg">
      <div class="searchform-header-row">
        <h2 class="searchform-blue-title">Book Search</h2>
        {/*
        <button
          type="button"
          class={`scan-btn${scanning ? " scan-btn-active" : ""}`}
          onClick={handleScanToggle}
          disabled={scanLoading}
        >
          {scanLoading
            ? "Loading..."
            : scanning
              ? "Stop Scan"
              : "Start Scan"}
        </button>
        */}
      </div>
      <form onSubmit={handleSubmit} class="searchform-blue-form">
        <input
          type="text"
          name="query"
          placeholder="Enter ISBN or book title..."
          class="searchform-blue-input"
        />
        <select
          name="searchType"
          class="searchform-blue-select"
        >
          <option value="title">Title</option>
          <option value="isbn">ISBN</option>
        </select>
        <button
          type="submit"
          disabled={loading}
          class={`searchform-blue-btn${loading ? " searchform-blue-btn-loading" : ""}`}
        >
          {loading ? (
            <span class="searchform-blue-spinner"></span>
          ) : (
            'Search'
          )}
        </button>
      </form>
    </div>
  );
}