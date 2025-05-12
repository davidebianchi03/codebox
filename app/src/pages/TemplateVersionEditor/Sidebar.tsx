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
            itemId={entry.name}
            label={entry.name}
            labelIcon={GetIconForFile(entry.name)}
            labelInfo="90"
        >
            {entry.children.map(e => <TreeEntry entry={e} />)}
        </SidebarTreeItem>
    )
}

interface TemplateVersionEditorSidebarProps {
    template: WorkspaceTemplate
    templateVersion: WorkspaceTemplateVersion
}

export function TemplateVersionEditorSidebar({ template, templateVersion }: TemplateVersionEditorSidebarProps) {
    const [treeItems, setTreeItems] = React.useState<WorkspaceTemplateVersionTreeItem[]>([]);

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
                    <SimpleTreeView>
                        {treeItems.map((item) =>
                            <TreeEntry entry={item} />
                        )}

                        {/* <SidebarTreeItem itemId="1" label="All Mail" labelIcon={MailIcon} />
                        <SidebarTreeItem itemId="2" label="Trash" labelIcon={DeleteIcon} />
                        <SidebarTreeItem itemId="3" label="Categories" labelIcon={Label}>
                            <SidebarTreeItem
                                itemId="5"
                                label="Social"
                                labelIcon={SupervisorAccountIcon}
                                labelInfo="90"
                            />
                            <SidebarTreeItem
                                itemId="6"
                                label="Updates"
                                labelIcon={InfoIcon}
                                labelInfo="2,294"
                            />
                            <SidebarTreeItem
                                itemId="7"
                                label="Forums"
                                labelIcon={ForumIcon}
                                labelInfo="3,566"
                            />
                            <SidebarTreeItem
                                itemId="8"
                                label="Promotions"
                                labelIcon={LocalOfferIcon}
                                labelInfo="733"
                            /> */}
                        {/* </SidebarTreeItem> */}
                        {/* <SidebarTreeItem itemId="4" label="History" labelIcon={Label} /> */}
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