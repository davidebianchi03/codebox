import React, { useCallback, useEffect, useState, useMemo } from "react";
import { Workspace, WorkspaceType } from "../../types/workspace";
import { Link, useNavigate } from "react-router-dom";
import { APIListWorkspaces, APIListWorkspacesTypes } from "../../api/workspace";
import { WorkspaceItem } from "./WorkspaceItem";
import { Button, Card, Container, Form, Pagination } from "react-bootstrap";
import { useNotifications } from "../../hooks/useNotifications";
import { CodeboxNotification } from "../../types/notifications";

export default function HomePage() {
  const navigate = useNavigate();

  const [searchText, setSearchText] = useState<string>("");
  const [allWorkspaces, setAllWorkspaces] = useState<Workspace[]>([]);
  const [workspaceTypes, setWorkspaceTypes] = useState<WorkspaceType[]>([]);
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);
  const [isLoading, setIsLoading] = useState<boolean>(false);

  const FetchWorkspaces = useCallback(async () => {
    setIsLoading(true);
    const w = await APIListWorkspaces();
    if (w) {
      setAllWorkspaces(w);
    }
    setIsLoading(false);
  }, []);

  const FetchWorkspaceTypes = useCallback(async () => {
    const wt = await APIListWorkspacesTypes();
    if (wt) {
      setWorkspaceTypes(wt);
    }
  }, []);

  // Handle notifications for real-time workspace updates
  const handleNotification = useCallback((notification: CodeboxNotification) => {
    if (notification.type === "workspace" && notification.workspace) {
      // Update the workspace in the list
      setAllWorkspaces((prevWorkspaces) =>
        prevWorkspaces.map((w) =>
          w.id === notification.workspace!.id ? notification.workspace! : w
        )
      );
    }
  }, []);

  useNotifications({
    onNotification: handleNotification,
    enabled: true,
  });

  useEffect(() => {
    FetchWorkspaces();
    FetchWorkspaceTypes();
  }, [FetchWorkspaces, FetchWorkspaceTypes]);

  // Filter and sort workspaces
  const filteredAndSortedWorkspaces = useMemo(() => {
    return allWorkspaces
      .filter((w) => w.name.toLowerCase().includes(searchText.toLowerCase()))
      .sort((w1: any, w2: any) => {
        const d1 = new Date(w1.updated_at);
        const d2 = new Date(w2.updated_at);
        if (d1 < d2) return 1;
        else if (d1 > d2) return -1;
        else return 0;
      });
  }, [allWorkspaces, searchText]);

  // Calculate pagination
  const totalWorkspaces = filteredAndSortedWorkspaces.length;
  const totalPages = Math.ceil(totalWorkspaces / pageSize);
  const offset = (currentPage - 1) * pageSize;
  const paginatedWorkspaces = filteredAndSortedWorkspaces.slice(offset, offset + pageSize);

  // Reset to first page when search text changes
  useEffect(() => {
    setCurrentPage(1);
  }, [searchText]);

  const handlePreviousPage = () => {
    if (currentPage > 1) {
      setCurrentPage(currentPage - 1);
      window.scrollTo({ top: 0, behavior: "smooth" });
    }
  };

  const handleNextPage = () => {
    if (currentPage < totalPages) {
      setCurrentPage(currentPage + 1);
      window.scrollTo({ top: 0, behavior: "smooth" });
    }
  };

  const handlePageSizeChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setPageSize(parseInt(e.target.value));
    setCurrentPage(1);
  };

  return (
    <Container className="pb-4 mt-4">
      <div className="row g-2 align-items-center">
        <div className="col">
          <div className="page-pretitle">Overview</div>
          <h2 className="page-title">Workspaces</h2>
        </div>
        <div className="col-auto ms-auto d-print-none">
          <div className="btn-list">
            <Button
              variant="light"
              onClick={() => navigate("/create-workspace")}
            >
              Create workspace
            </Button>
          </div>
        </div>
      </div>

      <Card className="my-5">
        <Card.Body>
          <div className="row align-items-center">
            <div className="col">
              <Form.Label>Filter workspaces:</Form.Label>
              <Form.Control
                placeholder="my awesome workspace"
                value={searchText}
                onChange={(e) => setSearchText(e.target.value)}
              />
            </div>
            <div className="col-auto">
              <Form.Label>Items per page:</Form.Label>
              <Form.Select
                value={pageSize}
                onChange={handlePageSizeChange}
                style={{ width: "auto" }}
              >
                <option value="5">5</option>
                <option value="10">10</option>
                <option value="20">20</option>
                <option value="50">50</option>
              </Form.Select>
            </div>
          </div>
        </Card.Body>
      </Card>

      <Card className="my-1">
        <Card.Body>
          {isLoading ? (
            <p className="text-center">Loading workspaces...</p>
          ) : paginatedWorkspaces.length > 0 ? (
            <>
              {paginatedWorkspaces.map((workspace: Workspace) => (
                <React.Fragment key={workspace.id}>
                  <WorkspaceItem workspace={workspace} workspaceTypes={workspaceTypes} />
                </React.Fragment>
              ))}

              {/* Pagination Controls */}
              {totalPages > 1 && (
                <div className="d-flex justify-content-between align-items-center mt-4 pt-3 border-top">
                  <div className="text-muted small">
                    Showing {paginatedWorkspaces.length > 0 ? offset + 1 : 0} to{" "}
                    {Math.min(offset + pageSize, totalWorkspaces)} of {totalWorkspaces} workspaces
                  </div>
                  <Pagination className="mb-0">
                    <Pagination.First
                      onClick={() => {
                        setCurrentPage(1);
                        window.scrollTo({ top: 0, behavior: "smooth" });
                      }}
                      disabled={currentPage === 1}
                    />
                    <Pagination.Prev
                      onClick={handlePreviousPage}
                      disabled={currentPage === 1}
                    />

                    {/* Show page numbers */}
                    {Array.from({ length: Math.min(5, totalPages) }).map((_, idx) => {
                      // Calculate which page number to show
                      let startPage = Math.max(1, currentPage - 2);
                      // Ensure we show 5 pages if possible, but don't go over totalPages
                      if (totalPages - startPage < 4) {
                        startPage = Math.max(1, totalPages - 4);
                      }
                      const pageNum = startPage + idx;

                      // Skip if page number exceeds total pages
                      if (pageNum > totalPages) return null;

                      return (
                        <Pagination.Item
                          key={pageNum}
                          active={pageNum === currentPage}
                          onClick={() => {
                            setCurrentPage(pageNum);
                            window.scrollTo({ top: 0, behavior: "smooth" });
                          }}
                        >
                          {pageNum}
                        </Pagination.Item>
                      );
                    })}

                    <Pagination.Next
                      onClick={handleNextPage}
                      disabled={currentPage === totalPages}
                    />
                    <Pagination.Last
                      onClick={() => {
                        setCurrentPage(totalPages);
                        window.scrollTo({ top: 0, behavior: "smooth" });
                      }}
                      disabled={currentPage === totalPages}
                    />
                  </Pagination>
                </div>
              )}
            </>
          ) : (
            <p className="text-center">
              {allWorkspaces.length > 0 ? (
                <>
                  No workspaces found matching '{searchText}'{" "}
                  <span>
                    <Link to={`/create-workspace?name=${searchText}`}>create it</Link>
                  </span>
                </>
              ) : (
                <>
                  No workspaces found,{" "}
                  <span>
                    <u
                      onClick={() => navigate("/create-workspace")}
                      style={{ cursor: "pointer" }}
                    >
                      create your first workspace
                    </u>
                  </span>
                </>
              )}
            </p>
          )}
        </Card.Body>
      </Card>
    </Container>
  );
}
