import React, { ChangeEventHandler } from "react";
import { Form } from "react-bootstrap";

type FormControlElement = HTMLInputElement | HTMLTextAreaElement;

export interface GitRepoSelectorField {
    value: string;
    placeholder?: string;
    onChange?: ChangeEventHandler<FormControlElement> | undefined;
    error?: string;
}

export interface GitRepoSelectorProps {
    repoUrl: GitRepoSelectorField,
    refName: GitRepoSelectorField,
    configFilePath: GitRepoSelectorField,
}

export function GitRepoSelector({
    repoUrl,
    refName,
    configFilePath,
}: GitRepoSelectorProps) {
    return (
        <React.Fragment>
            <Form.Group className="mt-2">
                <Form.Label>Repository URL</Form.Label>
                <Form.Control
                    name="gitRepositoryURL"
                    placeholder={repoUrl.placeholder || "git@example.com/my-awesome-project"}
                    value={repoUrl.value}
                    onChange={repoUrl.onChange}
                    isInvalid={
                        repoUrl.error !== undefined
                    }
                />
                <Form.Control.Feedback>
                    {repoUrl.error}
                </Form.Control.Feedback>
            </Form.Group>
            <Form.Group className="mt-2">
                <Form.Label>Ref Name</Form.Label>
                <Form.Control
                    name="gitRefName"
                    placeholder={refName.placeholder || "refs/heads/main"}
                    value={refName.value}
                    onChange={refName.onChange}
                    isInvalid={refName.error !== undefined}
                />
                <Form.Control.Feedback>
                    {refName.error}
                </Form.Control.Feedback>
            </Form.Group>
            <Form.Group className="mt-2">
                <Form.Label>Config files path</Form.Label>
                <Form.Control
                    name="configFilesPath"
                    placeholder={configFilePath.placeholder || "path/to/config"}
                    value={configFilePath.value}
                    onChange={configFilePath.onChange}
                    isInvalid={
                        configFilePath.error !== undefined
                    }
                />
                <Form.Control.Feedback>
                    {configFilePath.error}
                </Form.Control.Feedback>
            </Form.Group>
        </React.Fragment>
    )
}
