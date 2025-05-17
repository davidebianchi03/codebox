import * as React from 'react';
import Box from '@mui/material/Box';
import { SimpleTreeView } from '@mui/x-tree-view/SimpleTreeView';
import { Menu, MenuItem } from '@mui/material';
import { WorkspaceTemplate, WorkspaceTemplateVersion, WorkspaceTemplateVersionTreeItem } from '../../types/templates';
import { SidebarTreeItem } from './SidebarTreeItem';
import { Http } from '../../api/http';
import { RequestStatus } from '../../api/types';
import { toast } from 'react-toastify';
import { GetTypeForFile } from './FileType';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Button } from 'reactstrap';
import { faFileCirclePlus, faFolderPlus } from '@fortawesome/free-solid-svg-icons';
import Swal from 'sweetalert2';
import { Link } from 'react-router-dom';

interface SidebarEntryProps {
    entry: WorkspaceTemplateVersionTreeItem;
    onContextMenu: (event: React.MouseEvent, fullpath: WorkspaceTemplateVersionTreeItem) => void;
}

function TreeEntry({ entry, onContextMenu }: SidebarEntryProps) {
    return (
        <SidebarTreeItem
            itemId={entry.full_path}
            label={entry.name}
            labelIcon={GetTypeForFile(entry.name).icon}
            labelInfo="90"
            type={entry.type}
            onContextMenu={(e) => {
                e.stopPropagation();
                onContextMenu(e, entry);
            }}
        >
            {entry.children.map((e, index) => <TreeEntry entry={e} key={index} onContextMenu={onContextMenu} />)}
        </SidebarTreeItem>
    )
}

function GetTreeEntryByPath(path: string, tree: WorkspaceTemplateVersionTreeItem[]): WorkspaceTemplateVersionTreeItem | null {
    for (let i = 0; i < tree.length; i++) {
        if (tree[i].full_path === path) {
            return tree[i];
        }

        var fromChild = GetTreeEntryByPath(path, tree[i].children);
        if (fromChild) {
            return fromChild;
        }
    }
    return null;
}

interface TemplateVersionEditorSidebarProps {
    template: WorkspaceTemplate
    templateVersion: WorkspaceTemplateVersion
    onSelectionChange: (selectedItem: string) => void;
}

export function TemplateVersionEditorSidebar({ template, templateVersion, onSelectionChange }: TemplateVersionEditorSidebarProps) {
    const [treeItems, setTreeItems] = React.useState<WorkspaceTemplateVersionTreeItem[]>([]);
    const [selectedItem, setSelectedItem] = React.useState<string | null>(null);
    const [contextMenuEntry, setContextMenuEntry] = React.useState<WorkspaceTemplateVersionTreeItem | null>(null);

    const GetDirName = React.useCallback((path: string) => {
        if (!path || path === "/") return "";
        const segments = path.split("/").filter(Boolean);
        if (segments.length <= 1) return "/";
        return "/" + segments.slice(0, -1).join("/");
    }, [])

    React.useEffect(() => {
        if (selectedItem) {
            onSelectionChange(selectedItem);
        }
    }, [onSelectionChange, selectedItem]);

    const [contextMenu, setContextMenu] = React.useState<{
        mouseX: number;
        mouseY: number;
    } | null>(null);

    const handleContextMenu = React.useCallback((event: React.MouseEvent, entry: WorkspaceTemplateVersionTreeItem) => {
        event.preventDefault();
        setContextMenu(
            contextMenu === null
                ? {
                    mouseX: event.clientX + 2,
                    mouseY: event.clientY - 6
                }
                :
                null
        );

        setContextMenuEntry(entry);
    }, [contextMenu]);

    const handleCloseContextMenu = React.useCallback(() => {
        setContextMenu(null);
        setContextMenuEntry(null);
    }, []);

    const fetchTreeItems = React.useCallback(async () => {
        let [status, statusCode, responseBody] = await Http.Request(
            `${Http.GetServerURL()}/api/v1/templates/${template.id}/versions/${templateVersion.id}/entries`,
            "GET",
            null
        );
        if (status === RequestStatus.OK && statusCode === 200) {
            setTreeItems(responseBody as WorkspaceTemplateVersionTreeItem[]);
        } else {
            toast.error("Failed to fetch entries");
        }
    }, [template.id, templateVersion.id]);

    const handleCreateItem = React.useCallback(async (type: "dir" | "file") => {
        // get parent folder
        var parentEntry: WorkspaceTemplateVersionTreeItem | null = null;
        if (selectedItem) {
            parentEntry = GetTreeEntryByPath(selectedItem, treeItems);
            if (parentEntry) {
                if (parentEntry.type === "file") {
                    if (GetDirName(parentEntry.full_path) !== "") {
                        parentEntry = GetTreeEntryByPath(GetDirName(parentEntry.full_path), treeItems);
                    }
                }
            }
        }

        var r = await Swal.fire({
            title: `Enter the name of the ${type === "dir" ? "folder" : "file"}`,
            input: 'text',
            inputLabel: `${type === "dir" ? "Folder" : "File"} name`,
            inputPlaceholder: `Enter ${type === "dir" ? "folder" : "file"} name here`,
            showCancelButton: true,
            reverseButtons: true,
            inputValidator: async (value) => {
                if (!value) {
                    return 'You need to write something!'
                }

                var itemPath = value;
                // if parent entry is null the new item will be created in the root folder
                if (parentEntry !== null) {
                    itemPath = `${parentEntry.full_path}${!parentEntry.full_path.endsWith("/") && "/"}${value}`;
                }

                // check if item already exists
                let [status, statusCode] = await Http.Request(
                    `${Http.GetServerURL()}/api/v1/templates/${template.id}/versions/${templateVersion.id}/entries/${encodeURIComponent(itemPath)}`,
                    "GET",
                    null
                );

                if (status === RequestStatus.OK) {
                    return 'Item already exists';
                } else if (statusCode !== 404) {
                    return 'Failed to check if item already exists';
                }
            },
            customClass: {
                confirmButton: "btn btn-primary",
                cancelButton: "btn btn-light me-1",
                popup: "bg-dark text-light",
            },
            buttonsStyling: false,
        });

        if (r.isConfirmed && r.value) {
            var itemPath = r.value;
            // if parent entry is null the new item will be created in the root folder
            if (parentEntry !== null) {
                itemPath = `${parentEntry.full_path}${!parentEntry.full_path.endsWith("/") && "/"}${r.value}`;
            }

            // create file
            let [status] = await Http.Request(
                `${Http.GetServerURL()}/api/v1/templates/${template.id}/versions/${templateVersion.id}/entries`,
                "POST",
                JSON.stringify({
                    path: itemPath,
                    type: type,
                    content: "",
                })
            );
            if (status !== RequestStatus.OK) {
                toast.error("Cannot add file");
            }
        }

        fetchTreeItems();
    }, [GetDirName, fetchTreeItems, selectedItem, template.id, templateVersion.id, treeItems]);

    const handleDeleteEntry = React.useCallback(async () => {
        if (contextMenuEntry) {
            await Http.Request(
                `
                ${Http.GetServerURL()}/api/v1/templates/${template.id}/versions/
                ${templateVersion.id}/entries/${encodeURIComponent(contextMenuEntry.full_path)}
                `,
                "DELETE",
                null,
            );
            setContextMenu(null);
            setContextMenuEntry(null);
            fetchTreeItems();
        }
    }, [contextMenuEntry, fetchTreeItems, template.id, templateVersion.id]);

    const handleRenameEntry = React.useCallback(async () => {
        if (contextMenuEntry) {
            var r = await Swal.fire({
                title: `Change ${contextMenuEntry.type === "dir" ? "folder" : "file"} name`,
                input: 'text',
                inputLabel: `${contextMenuEntry.type === "dir" ? "Folder" : "File"} name`,
                inputPlaceholder: `Enter ${contextMenuEntry.type === "dir" ? "folder" : "file"} name here`,
                showCancelButton: true,
                reverseButtons: true,
                inputValidator: async (value) => {
                    // if (!value) {
                    //     return 'You need to write something!'
                    // }

                    // if (parentFolderPath !== "" && parentEntry === null) {
                    //     return 'Parent folder does not exist!'
                    // }

                    // // check if item already exists
                    // var itemPath = parentFolderPath + (parentFolderPath.endsWith("/") || parentFolderPath === "" ? "" : "/") + value;
                    // let [status, statusCode] = await Http.Request(
                    //     `${Http.GetServerURL()}/api/v1/templates/${template.id}/versions/${templateVersion.id}/entries/${encodeURIComponent(itemPath)}`,
                    //     "GET",
                    //     null
                    // );

                    // if (status === RequestStatus.OK) {
                    //     return 'Item already exists';
                    // } else if (statusCode !== 404) {
                    //     return 'Failed to check if item already exists';
                    // }
                },
                customClass: {
                    confirmButton: "btn btn-primary",
                    cancelButton: "btn btn-light me-1",
                    popup: "bg-dark text-light",
                },
                buttonsStyling: false,
            });


            setContextMenu(null);
            setContextMenuEntry(null);
            fetchTreeItems();
        }
    }, [contextMenuEntry, fetchTreeItems]);

    React.useEffect(() => {
        fetchTreeItems();
    }, [fetchTreeItems]);

    const sidebarContainer = React.useRef<any>(null);

    return (
        <React.Fragment>
            <div className='d-flex justify-content-between pt-1' style={{ height: 25 }}>
                <span className='text-uppercase px-2' style={{ fontFamily: "Consolas", fontSize: "12", fontWeight: "bold" }}>
                    <Link to={`/templates/${template.id}`}>
                        {template.name}
                    </Link>
                </span>
                <div className='d-flex justify-content-end'>
                    <Button
                        size='sm'
                        className='text-center'
                        style={{ background: "none", border: "none", fontSize: 14 }}
                        onClick={() => handleCreateItem("dir")}
                    >
                        <FontAwesomeIcon icon={faFolderPlus} />
                    </Button>
                    <Button
                        size='sm'
                        className='text-center'
                        style={{ background: "none", border: "none", fontSize: 14 }}
                        onClick={() => handleCreateItem("file")}
                    >
                        <FontAwesomeIcon icon={faFileCirclePlus} />
                    </Button>
                </div>
            </div>
            <div
                // onContextMenu={handleContextMenu}
                style={{ cursor: "context-menu" }}
                onClick={(e) => {
                    if (sidebarContainer.current === e.target) {
                        setSelectedItem(null);
                    }
                }}
            >
                <Box sx={{ minHeight: 352, minWidth: 250, height: "calc(100% - 45px)" }} ref={sidebarContainer}>
                    <SimpleTreeView
                        onSelectedItemsChange={(e, item) => {
                            setSelectedItem(item);
                        }}
                    >
                        {treeItems.map((item, index) =>
                            <TreeEntry entry={item} key={index} onContextMenu={handleContextMenu} />
                        )}
                    </SimpleTreeView>
                    <Menu
                        open={contextMenu !== null}
                        onClose={handleCloseContextMenu}
                        anchorReference="anchorPosition"
                        anchorPosition={
                            contextMenu !== null
                                ? { top: contextMenu.mouseY, left: contextMenu.mouseX }
                                : undefined
                        }
                        slotProps={{
                            paper: {
                                style: {
                                    backgroundColor: '#1e1e1e',
                                    color: '#fff',
                                    borderRadius: 8,
                                    minWidth: 160,
                                },
                            }
                        }}
                    >
                        <MenuItem
                            sx={{
                                '&:hover': {
                                    backgroundColor: '#333',
                                },
                                fontSize: 14
                            }}
                            onClick={handleRenameEntry}
                        >
                            Rename
                        </MenuItem>
                        <MenuItem
                            sx={{
                                '&:hover': {
                                    backgroundColor: '#333',
                                },
                                fontSize: 14
                            }}
                            onClick={handleDeleteEntry}
                        >
                            Delete permanently
                        </MenuItem>
                    </Menu>
                </Box>
            </div>
        </React.Fragment >
    )
}