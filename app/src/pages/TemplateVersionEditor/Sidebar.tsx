import * as React from 'react';
import Box from '@mui/material/Box';
import SupervisorAccountIcon from '@mui/icons-material/SupervisorAccount';
import { SimpleTreeView } from '@mui/x-tree-view/SimpleTreeView';
import { Menu, MenuItem } from '@mui/material';
import { WorkspaceTemplate, WorkspaceTemplateVersion, WorkspaceTemplateVersionTreeItem } from '../../types/templates';
import { SidebarTreeItem } from './SidebarTreeItem';
import { Http } from '../../api/http';
import { RequestStatus } from '../../api/types';
import { toast } from 'react-toastify';
import { GetIconForFile } from './FileIcon';

interface SidebarEntryProps {
    entry: WorkspaceTemplateVersionTreeItem
}

function TreeEntry({ entry }: SidebarEntryProps) {
    return (
        <SidebarTreeItem
            itemId={entry.full_path}
            label={entry.name}
            labelIcon={GetIconForFile(entry.name)}
            labelInfo="90"
        >
            {entry.children.map((e, index) => <TreeEntry entry={e} key={index} />)}
        </SidebarTreeItem>
    )
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
    }, [selectedItem]);

    const [contextMenu, setContextMenu] = React.useState<{
        mouseX: number;
        mouseY: number;
    } | null>(null);

    const [selectedNodeId, setSelectedNodeId] = React.useState<string | null>(
        null
    );

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
    }, []);

    React.useEffect(() => {
        fetchTreeItems();
    }, [fetchTreeItems]);

    return (
        <React.Fragment>
            <div onContextMenu={handleContextMenu} style={{ cursor: "context-menu" }}>
                <Box sx={{ minHeight: 352, minWidth: 250 }}>
                    <SimpleTreeView
                        onSelectedItemsChange={(e, item) => {
                            if (item) {
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
                    >
                        <MenuItem onClick={handleClose}>Delete node {selectedNodeId}</MenuItem>
                    </Menu>
                </Box>
            </div>
        </React.Fragment >
    )
}