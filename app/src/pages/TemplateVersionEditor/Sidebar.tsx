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

interface SidebarEntryProps {
    entry: WorkspaceTemplateVersionTreeItem
}

function TreeEntry({ entry }: SidebarEntryProps) {
    return (
        <SidebarTreeItem
            itemId={entry.full_path}
            label={entry.name}
            labelIcon={GetTypeForFile(entry.name).icon}
            labelInfo="90"
            type={entry.type}
        >
            {entry.children.map((e, index) => <TreeEntry entry={e} key={index} />)}
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
    const [selectedItem, setSelectedItem] = React.useState<string>();

    React.useEffect(() => {
        if (selectedItem) {
            onSelectionChange(selectedItem);
        }
    }, [onSelectionChange, selectedItem]);

    const [contextMenu, setContextMenu] = React.useState<{
        mouseX: number;
        mouseY: number;
    } | null>(null);

    const handleContextMenu = (event: React.MouseEvent) => {
        event.preventDefault();
        setContextMenu(
            contextMenu === null
                ? {
                    mouseX: event.clientX + 2,
                    mouseY: event.clientY - 6
                }
                : // repeated contextmenu when it is already open closes it with Chrome 84 on Ubuntu
                // Other native context menus might behave different.
                // With this behavior we prevent contextmenu from the backdrop to re-locale existing context menus.
                null
        );
    };

    const handleClose = () => {
        setContextMenu(null);
    };

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

    const handleCreateFile = React.useCallback(async () => {
        // get parent folder
        var parentFolderPath = "";
        var parentEntry: WorkspaceTemplateVersionTreeItem | null = null;
        if (selectedItem) {
            parentFolderPath = selectedItem;
            parentEntry = GetTreeEntryByPath(parentFolderPath, treeItems);
            if (parentEntry) {
                if (parentEntry.type === "file") {
                    var p = parentFolderPath.split("/");
                    parentFolderPath = p.filter((part, index) => index < p.length - 1).join("/");
                    if (parentFolderPath !== "") {
                        parentEntry = GetTreeEntryByPath(parentFolderPath, treeItems);
                    }
                }
            }
        }

        var r = await Swal.fire({
            title: 'Enter the name of the file',
            input: 'text',
            inputLabel: 'File name',
            inputPlaceholder: 'Enter file name here',
            showCancelButton: true,
            reverseButtons: true,
            inputValidator: (value) => {
                if (!value) {
                    return 'You need to write something!'
                }

                if (parentFolderPath !== "" && parentEntry === null) {
                    return 'Parent folder does not exist!'
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
            let [status] = await Http.Request(
                `${Http.GetServerURL()}/api/v1/templates/${template.id}/versions/${templateVersion.id}/entries`,
                "POST",
                JSON.stringify({
                    path: parentFolderPath + (parentFolderPath.endsWith("/") || parentFolderPath === "" ? "" : "/") + r.value,
                    type: "file",
                    content: "",
                })
            );
            if (status !== RequestStatus.OK) {
                toast.error("Cannot add file");
            }
        }

        fetchTreeItems();
    }, [fetchTreeItems, selectedItem, template.id, templateVersion.id, treeItems]);

    const handleCreateFolder = React.useCallback(async () => {
        // get parent folder
        var parentFolderPath = "";
        var parentEntry: WorkspaceTemplateVersionTreeItem | null = null;
        if (selectedItem) {
            parentFolderPath = selectedItem;
            parentEntry = GetTreeEntryByPath(parentFolderPath, treeItems);
            if (parentEntry) {
                if (parentEntry.type === "file") {
                    var p = parentFolderPath.split("/");
                    parentFolderPath = p.filter((part, index) => index < p.length - 1).join("/");
                    if (parentFolderPath !== "") {
                        parentEntry = GetTreeEntryByPath(parentFolderPath, treeItems);
                    }
                }
            }
        }

        var r = await Swal.fire({
            title: 'Enter the name of the folder',
            input: 'text',
            inputLabel: 'Folder name',
            inputPlaceholder: 'Enter folder name here',
            showCancelButton: true,
            reverseButtons: true,
            inputValidator: (value) => {
                if (!value) {
                    return 'You need to write something!'
                }

                if (parentFolderPath !== "" && parentEntry === null) {
                    return 'Parent folder does not exist!'
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
            let [status] = await Http.Request(
                `${Http.GetServerURL()}/api/v1/templates/${template.id}/versions/${templateVersion.id}/entries`,
                "POST",
                JSON.stringify({
                    path: parentFolderPath + (parentFolderPath.endsWith("/") || parentFolderPath === "" ? "" : "/") + r.value,
                    type: "dir",
                    content: "",
                })
            );
            if (status !== RequestStatus.OK) {
                toast.error("Cannot add file");
            }
        }

        fetchTreeItems();
    }, [fetchTreeItems, selectedItem, template.id, templateVersion.id, treeItems]);

    React.useEffect(() => {
        fetchTreeItems();
    }, [fetchTreeItems]);

    return (
        <React.Fragment>
            <div className='d-flex justify-content-between pt-1'>
                <span className='text-uppercase px-2' style={{ fontFamily: "Consolas", fontSize: "12", fontWeight: "bold" }}>
                    {template.name}
                </span>
                <div className='d-flex justify-content-end'>
                    <Button size='sm' className='text-center' style={{ background: "none", border: "none", fontSize: 14 }} onClick={handleCreateFolder}>
                        <FontAwesomeIcon icon={faFolderPlus} />
                    </Button>
                    <Button size='sm' className='text-center' style={{ background: "none", border: "none", fontSize: 14 }} onClick={handleCreateFile}>
                        <FontAwesomeIcon icon={faFileCirclePlus} />
                    </Button>
                </div>
            </div>
            <div onContextMenu={handleContextMenu} style={{ cursor: "context-menu" }}>
                <Box sx={{ minHeight: 352, minWidth: 250 }}>
                    <SimpleTreeView
                        onSelectedItemsChange={(e, item) => {
                            if (item) {
                                // find item
                                setSelectedItem(item);
                            }
                        }}
                    >
                        {treeItems.map((item, index) =>
                            <TreeEntry entry={item} key={index} />
                        )}
                    </SimpleTreeView>
                    <Menu
                        open={contextMenu !== null}
                        onClose={handleClose}
                        anchorReference="anchorPosition"
                        anchorPosition={
                            contextMenu !== null
                                ? { top: contextMenu.mouseY, left: contextMenu.mouseX }
                                : undefined
                        }
                    // PaperProps={{
                    //     style: {
                    //         backgroundColor: '#1e1e1e',
                    //         color: '#fff',
                    //         borderRadius: 8,
                    //         minWidth: 160,
                    //     },
                    // }}
                    // MenuListProps={{
                    //     sx: {
                    //         paddingY: 0.5,
                    //     },
                    // }}
                    >
                        <MenuItem
                            sx={{
                                '&:hover': {
                                    backgroundColor: '#333',
                                },
                            }}
                        >
                            Delete node
                        </MenuItem>
                    </Menu>
                </Box>
            </div>
        </React.Fragment >
    )
}