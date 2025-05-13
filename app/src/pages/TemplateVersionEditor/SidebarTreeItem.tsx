import * as React from 'react';
import { styled } from '@mui/material/styles';
import Box from '@mui/material/Box';
import { SvgIconProps } from '@mui/material/SvgIcon';
import {
    TreeItemContent,
    TreeItemIconContainer,
    TreeItemRoot,
    TreeItemGroupTransition,
} from '@mui/x-tree-view/TreeItem';
import { useTreeItem, UseTreeItemParameters } from '@mui/x-tree-view/useTreeItem';
import { TreeItemProvider } from '@mui/x-tree-view/TreeItemProvider';
import { TreeItemIcon } from '@mui/x-tree-view/TreeItemIcon';

interface StyledTreeItemProps
    extends Omit<UseTreeItemParameters, 'rootRef'>,
    React.HTMLAttributes<HTMLLIElement> {
    labelIcon: string;
    labelInfo?: string;
}

const SidebarTreeItemContent = styled(TreeItemContent)(({ theme }) => ({
    marginBottom: theme.spacing(0.3),
    borderRadius: 4,
    paddingRight: theme.spacing(1),
    paddingLeft: `calc(${theme.spacing(1)} + var(--TreeView-itemChildrenIndentation) * var(--TreeView-itemDepth) * 0.8)`,
    fontWeight: theme.typography.fontWeightMedium,
    fontSize: 12,
}));

const SidebarTreeItemIconContainer = styled(TreeItemIconContainer)(({ theme }) => ({
    marginRight: theme.spacing(1),
}));

export const SidebarTreeItem = React.forwardRef(function CustomTreeItem(
    props: StyledTreeItemProps,
    ref: React.Ref<HTMLLIElement>,
) {
    const {
        id,
        itemId,
        label,
        disabled,
        children,
        color,
        labelIcon,
        labelInfo,
        ...other
    } = props;

    const {
        getContextProviderProps,
        getRootProps,
        getContentProps,
        getIconContainerProps,
        getLabelProps,
        getGroupTransitionProps,
        status,
    } = useTreeItem({ id, itemId, children, label, disabled, rootRef: ref });

    return (
        <TreeItemProvider {...getContextProviderProps()}>
            <TreeItemRoot
                {...getRootProps(other)}
            >
                <SidebarTreeItemContent {...getContentProps()}>
                    <SidebarTreeItemIconContainer {...getIconContainerProps()}>
                        <TreeItemIcon status={status} />
                    </SidebarTreeItemIconContainer>
                    <Box
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            p: 0.1,
                            pr: 0,
                        }}
                    >
                        {(children as any[]).length == 0 && (
                            <img src={labelIcon} alt='file icon' width={16} height={16} className='me-2'/>
                        )}
                        <span>
                            {label}
                        </span>
                    </Box>
                </SidebarTreeItemContent>
                {children && <TreeItemGroupTransition {...getGroupTransitionProps()} />}
            </TreeItemRoot>
        </TreeItemProvider>
    );
});