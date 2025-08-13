import * as React from 'react';
import Box from '@mui/material/Box';
import { SimpleTreeView } from '@mui/x-tree-view/SimpleTreeView';
import { Menu, MenuItem } from '@mui/material';
import { WorkspaceTemplate, WorkspaceTemplateVersion, WorkspaceTemplateVersionTreeItem } from '../../types/templates';
import { SidebarTreeItem } from './SidebarTreeItem';
import { toast } from 'react-toastify';
import { GetTypeForFile } from './FileType';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Button } from 'reactstrap';
import { faFileCirclePlus, faFolderPlus } from '@fortawesome/free-solid-svg-icons';
import Swal from 'sweetalert2';
import { Link } from 'react-router-dom';
import { APICreateTemplateVersionEntry, APIDeleteTemplateVersionEntry, APIListTemplateVersionEntry, APIRetrieveTemplateVersionEntry, APIUpdateTemplateVersionEntry } from '../../api/templates';

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
            {entry.children.map((e, index) => <TreeEntry entry={e} key={index} onContextMenu={onContextMenu}/>)}
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
    onDelete: (fullpath: string) => void;
}

export function TemplateVersionEditorSidebar({ template, templateVersion, onSelectionChange, onDelete }: TemplateVersionEditorSidebarProps) {
    const [treeItems, setTreeItems] = React.useState<WorkspaceTemplateVersionTreeItem[]>([]);
    const [selectedItem, setSelectedItem] = React.useState<WorkspaceTemplateVersionTreeItem | null>(null);
    const [contextMenuEntry, setContextMenuEntry] = React.useState<WorkspaceTemplateVersionTreeItem | null>(null);

    const GetDirName = React.useCallback((path: string) => {
        if (!path || path === "/") return "";
        const segments = path.split("/").filter(Boolean);
        if (segments.length <= 1) return "";
        return segments.slice(0, -1).join("/");
    }, [])

    React.useEffect(() => {
        if (selectedItem) {
            onSelectionChange(selectedItem.full_path);
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
        const entries = await APIListTemplateVersionEntry(template.id, templateVersion.id);
        if (entries) {
            setTreeItems(entries);
        } else {
            toast.error("Failed to fetch entries");
        }
    }, [template.id, templateVersion.id]);

    const handleCreateItem = React.useCallback(async (type: "dir" | "file") => {
        var defaultInputValue = "";
        if (selectedItem) {
            defaultInputValue = selectedItem.type === "dir" ? selectedItem.full_path : GetDirName(selectedItem.full_path);
            if (!defaultInputValue.endsWith("/") && defaultInputValue.length > 0) {
                defaultInputValue += "/";
            }
        }

        var r = await Swal.fire({
            title: `Enter the name of the ${type === "dir" ? "folder" : "file"}`,
            input: 'text',
            inputLabel: `${type === "dir" ? "Folder" : "File"} name`,
            inputPlaceholder: `Enter ${type === "dir" ? "folder" : "file"} name here`,
            inputValue: defaultInputValue,
            showCancelButton: true,
            reverseButtons: true,
            inputValidator: async (value) => {
                if (!value) {
                    return 'You need to write something!'
                }

                if (value.startsWith("/")) {
                    return "Path cannot start with '/'";
                }

                if (value.endsWith("/")) {
                    return "Path cannot end with a trailing slash";
                }

                // check if parent entry is a folder
                if (GetDirName(value) !== "") {
                    var segments = GetDirName(value).split("/");
                    for (let i = 0; i < segments.length; i++) {
                        var parentEntry: WorkspaceTemplateVersionTreeItem | null = GetTreeEntryByPath(segments.slice(0, i + 1).join("/"), treeItems);
                        if (parentEntry) {
                            if (parentEntry.type !== "dir") {
                                return `Cannot create ${type === "dir" ? "folder" : "file"}, parent path is not a folder.`;
                            }
                        }
                    }
                }

                if (await APIRetrieveTemplateVersionEntry(template.id, templateVersion.id, value)) {
                    return 'Path already exists';
                }
            },
            customClass: {
                confirmButton: "btn btn-primary",
                cancelButton: "btn btn-accent me-1",
                popup: "bg-dark text-light",
            },
            buttonsStyling: false,
        });

        if (r.isConfirmed && r.value) {
            if ((await APICreateTemplateVersionEntry(template.id, templateVersion.id, r.value, type, "")) === undefined) {
                toast.error(`Failed to create ${type === "dir" ? "folder" : "file"}`);
            }
        }

        fetchTreeItems();
    }, [GetDirName, fetchTreeItems, selectedItem, template.id, templateVersion.id, treeItems]);

    const handleDeleteEntry = React.useCallback(async () => {
        setContextMenu(null);
        if (contextMenuEntry) {
            if ((await Swal.fire({
                icon: "question",
                title: `Are you sure you want to permanently delete '${contextMenuEntry.name}'?`,
                text: "You cannot restore it",
                showCancelButton: true,
                reverseButtons: true,
                customClass: {
                    confirmButton: "btn btn-primary",
                    cancelButton: "btn btn-accent me-1",
                    popup: "bg-dark text-light",
                },
                buttonsStyling: false,
            })).isConfirmed) {
                await APIDeleteTemplateVersionEntry(template.id, templateVersion.id, contextMenuEntry.full_path);
                fetchTreeItems();
                onDelete(contextMenuEntry.full_path);
            }
        }
        setContextMenuEntry(null);
    }, [contextMenuEntry, fetchTreeItems, template.id, templateVersion.id]);

    const handleRenameEntry = React.useCallback(async () => {
        setContextMenu(null);
        if (contextMenuEntry) {
            var r = await Swal.fire({
                title: `Change ${contextMenuEntry.type === "dir" ? "folder" : "file"} name`,
                input: 'text',
                inputLabel: `${contextMenuEntry.type === "dir" ? "Folder" : "File"} name`,
                inputPlaceholder: `Enter ${contextMenuEntry.type === "dir" ? "folder" : "file"} name here`,
                inputValue: contextMenuEntry.full_path,
                showCancelButton: true,
                reverseButtons: true,
                inputValidator: async (value) => {
                    if (!value) {
                        return 'You need to write something!'
                    }

                    if (value.startsWith("/")) {
                        return "Path cannot start with '/'";
                    }

                    if (value.endsWith("/")) {
                        return "Path cannot end with a trailing slash";
                    }

                    // check if parent entry is a folder
                    if (GetDirName(value) !== "") {
                        var segments = GetDirName(value).split("/");
                        for (let i = 0; i < segments.length; i++) {
                            var parentEntry: WorkspaceTemplateVersionTreeItem | null = GetTreeEntryByPath(segments.slice(0, i + 1).join("/"), treeItems);
                            if (parentEntry) {
                                if (parentEntry.type !== "dir") {
                                    return `Cannot move ${contextMenuEntry.type === "dir" ? "folder" : "file"}, parent path is not a folder.`;
                                }
                            }
                        }
                    }

                    if (value !== contextMenuEntry.full_path) {
                        const entry = await APIRetrieveTemplateVersionEntry(template.id, templateVersion.id, value);
                        if (entry) {
                            return 'Path already exists';
                        } else if (entry === undefined) {
                            return 'Failed to check if path already exists';
                        }
                    }
                },
                customClass: {
                    confirmButton: "btn btn-primary",
                    cancelButton: "btn btn-accent me-1",
                    popup: "bg-dark text-light",
                },
                buttonsStyling: false,
            });

            if (r.isConfirmed && r.value) {
                var fileContent = btoa("");
                if (contextMenuEntry.type === "file") {
                    const entry = await APIRetrieveTemplateVersionEntry(template.id, templateVersion.id, contextMenuEntry.full_path);
                    if (entry) {
                        fileContent = entry.content;
                    }
                }

                if ((await APIUpdateTemplateVersionEntry(
                    template.id,
                    templateVersion.id,
                    contextMenuEntry.full_path,
                    r.value,
                    fileContent
                )) === undefined) {
                    toast.error(`Failed to rename ${contextMenuEntry.type === "dir" ? "folder" : "file"}`);
                } else {
                    onSelectionChange(r.value);
                }
            }

            fetchTreeItems();
        }
        setContextMenuEntry(null);
    }, [GetDirName, contextMenuEntry, fetchTreeItems, onSelectionChange, template.id, templateVersion.id, treeItems]);

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
                            setSelectedItem(item ? GetTreeEntryByPath(item, treeItems) : null);
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