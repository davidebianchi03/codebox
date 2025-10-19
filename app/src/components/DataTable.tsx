import React, { useMemo, useState } from "react";

// Make sure to install bootstrap in your project and import its CSS once (e.g. in index.tsx):
// import 'bootstrap/dist/css/bootstrap.min.css';

// Usage: <DataTable columns={columns} data={data} initialPageSize={10} pageSizeOptions={[5,10,20]} />

type Column = {
  key: string;
  label: string;
  sortable?: boolean;
  // optional render function for custom cell content
  render?: (value: any, row: any) => React.ReactNode;
};

type DataTableProps = {
  columns: Column[];
  data: any[];
  initialPageSize?: number;
  pageSizeOptions?: number[];
  className?: string;
};

const clamp = (v: number, a: number, b: number) => Math.max(a, Math.min(b, v));

export default function DataTable({
  columns,
  data,
  initialPageSize = 10,
  pageSizeOptions = [5, 10, 20, 50],
  className = "",
}: DataTableProps) {
  const [query, setQuery] = useState("");
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(initialPageSize);
  const [sortKey, setSortKey] = useState<string | null>(null);
  const [sortDir, setSortDir] = useState<"asc" | "desc">("asc");

  // filter
  const filtered = useMemo(() => {
    if (!query) return data;
    const q = query.trim().toLowerCase();
    return data.filter((row) =>
      columns.some((col) => {
        const raw = row[col.key];
        if (raw === null || raw === undefined) return false;
        return String(raw).toLowerCase().includes(q);
      })
    );
  }, [data, query, columns]);

  // sort
  const sorted = useMemo(() => {
    if (!sortKey) return filtered;
    const copy = [...filtered];
    copy.sort((a, b) => {
      const A = a[sortKey];
      const B = b[sortKey];
      if (A == null && B == null) return 0;
      if (A == null) return -1;
      if (B == null) return 1;
      // numeric comparison when both are numbers
      if (typeof A === "number" && typeof B === "number") {
        return sortDir === "asc" ? A - B : B - A;
      }
      // date comparison when both are dates
      const dateA = new Date(A);
      const dateB = new Date(B);
      if (!isNaN(dateA.getTime()) && !isNaN(dateB.getTime())) {
        return sortDir === "asc" ? dateA.getTime() - dateB.getTime() : dateB.getTime() - dateA.getTime();
      }
      // fallback to string compare
      const sa = String(A).toLowerCase();
      const sb = String(B).toLowerCase();
      if (sa < sb) return sortDir === "asc" ? -1 : 1;
      if (sa > sb) return sortDir === "asc" ? 1 : -1;
      return 0;
    });
    return copy;
  }, [filtered, sortKey, sortDir]);

  const total = sorted.length;
  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  const currentPage = clamp(page, 1, totalPages);

  const pageData = useMemo(() => {
    const start = (currentPage - 1) * pageSize;
    return sorted.slice(start, start + pageSize);
  }, [sorted, currentPage, pageSize]);

  // page number range (show up to 7 pages window)
  const pageRange = useMemo(() => {
    const maxButtons = 7;
    const half = Math.floor(maxButtons / 2);
    let start = clamp(currentPage - half, 1, Math.max(1, totalPages - maxButtons + 1));
    let end = Math.min(totalPages, start + maxButtons - 1);
    // adjust start if we're at the tail
    start = Math.max(1, end - maxButtons + 1);
    const arr: number[] = [];
    for (let i = start; i <= end; i++) arr.push(i);
    return arr;
  }, [currentPage, totalPages]);

  // handlers
  function onSearchChange(e: React.ChangeEvent<HTMLInputElement>) {
    setQuery(e.target.value);
    setPage(1);
  }

  function onPageSizeChange(e: React.ChangeEvent<HTMLSelectElement>) {
    const ps = Number(e.target.value) || initialPageSize;
    setPageSize(ps);
    setPage(1);
  }

  function onSort(col: Column) {
    if (!col.sortable) return;
    if (sortKey === col.key) {
      setSortDir((d) => (d === "asc" ? "desc" : "asc"));
    } else {
      setSortKey(col.key);
      setSortDir("asc");
    }
  }

  return (
    <div className={`datatable ${className}`}>
      <div className="d-flex flex-column flex-md-row justify-content-between align-items-center gap-2 mb-3">
        <div className="d-flex gap-2 w-100">
          <input
            className="form-control"
            placeholder="Search..."
            value={query}
            onChange={onSearchChange}
            aria-label="Search table"
          />
          <select className="form-select" value={pageSize} onChange={onPageSizeChange} style={{ width: "auto" }}>
            {pageSizeOptions.map((opt) => (
              <option key={opt} value={opt}>
                {opt} / page
              </option>
            ))}
          </select>
        </div>
        <div className="text-muted small mt-1 text-end" style={{ minWidth: "100px" }}>
          Showing <strong>{pageData.length}</strong> of <strong>{total}</strong>
        </div>
      </div>

      <div className="table-responsive">
        <table className="table table-striped table-hover">
          <thead className="table-light">
            <tr>
              {columns.map((col) => (
                <th
                  key={col.key}
                  scope="col"
                  style={{ cursor: col.sortable ? "pointer" : "default" }}
                  onClick={() => onSort(col)}
                  aria-sort={sortKey === col.key ? (sortDir === "asc" ? "ascending" : "descending") : "none"}
                >
                  <div className="d-flex align-items-center">
                    <span>{col.label}</span>
                    {col.sortable && (
                      <small className="ms-2 text-muted">{sortKey === col.key ? (sortDir === "asc" ? "▲" : "▼") : "↕"}</small>
                    )}
                  </div>
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {pageData.length === 0 ? (
              <tr>
                <td colSpan={columns.length} className="text-center text-muted py-4">
                  No results
                </td>
              </tr>
            ) : (
              pageData.map((row, idx) => (
                <tr key={idx}>
                  {columns.map((col) => (
                    <td key={col.key}>
                      {col.render ? col.render(row[col.key], row) : String(row[col.key] ?? "")}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      <nav aria-label="Table pagination" className="mt-3">
        <ul className="pagination justify-content-center">
          <li className={`page-item ${currentPage === 1 ? "disabled" : ""}`}>
            <button className="page-link" onClick={() => setPage(1)} aria-label="First" disabled={currentPage === 1}>
              «
            </button>
          </li>
          <li className={`page-item ${currentPage === 1 ? "disabled" : ""}`}>
            <button className="page-link" onClick={() => setPage((p) => clamp(p - 1, 1, totalPages))} aria-label="Previous" disabled={currentPage === 1}>
              &#8249;
            </button>
          </li>

          {pageRange[0] > 1 && (
            <li className="page-item disabled">
              <span className="page-link">…</span>
            </li>
          )}

          {pageRange.map((p) => (
            <li key={p} className={`page-item ${p === currentPage ? "active" : ""}`}>
              <button className="page-link" onClick={() => setPage(p)}>
                {p}
              </button>
            </li>
          ))}

          {pageRange[pageRange.length - 1] < totalPages && (
            <li className="page-item disabled">
              <span className="page-link">…</span>
            </li>
          )}

          <li className={`page-item ${currentPage === totalPages ? "disabled" : ""}`}>
            <button className="page-link" onClick={() => setPage((p) => clamp(p + 1, 1, totalPages))} aria-label="Next" disabled={currentPage === totalPages}>
              &#8250;
            </button>
          </li>
          <li className={`page-item ${currentPage === totalPages ? "disabled" : ""}`}>
            <button className="page-link" onClick={() => setPage(totalPages)} aria-label="Last" disabled={currentPage === totalPages}>
              »
            </button>
          </li>
        </ul>
      </nav>
      <style>{`
        .datatable .table thead th {
          user-select: none;
        }
      `}</style>
    </div>
  );
}

// ---------------------------
// Example usage (put in App.tsx)
// ---------------------------
/*
import React from 'react';
import DataTable from './DataTable';
import 'bootstrap/dist/css/bootstrap.min.css';

const columns = [
  { key: 'id', label: 'ID', sortable: true },
  { key: 'name', label: 'Name', sortable: true },
  { key: 'email', label: 'Email' },
  { key: 'registered', label: 'Registered', sortable: true, render: (v) => new Date(v).toLocaleDateString() }
];

const data = Array.from({ length: 137 }).map((_, i) => ({
  id: i + 1,
  name: `User ${i + 1}`,
  email: `user${i + 1}@example.com`,
  registered: new Date(Date.now() - Math.random() * 1000 * 60 * 60 * 24 * 365 * 3).toISOString(),
}));

export default function App() {
  return (
    <div className="container mt-4">
      <h2>Example DataTable</h2>
      <DataTable columns={columns} data={data} initialPageSize={10} pageSizeOptions={[5,10,25,50]} />
    </div>
  );
}
*/
