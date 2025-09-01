import React, { useEffect, useState } from "react";
import { Button, Input, InputGroup, InputGroupText, Row } from "reactstrap";
import { faPlus, faTrashCan } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

interface EnvEditorProps {
    value?: string;
    onChange?: (value: string) => void;
    invalid?: boolean;
}

export function EnvEditor({ value = "", onChange = () => { }, invalid }: EnvEditorProps) {
    const lines = value.split("\n");
    const [advancedModeEnabled, setAdvancedModeEnabled] = useState(false);

    return (
        <React.Fragment>
            {advancedModeEnabled ? (
                <React.Fragment>
                    <Button color="link" className="p-0 my-2" onClick={(e) => {
                        e.preventDefault();
                        setAdvancedModeEnabled(false);
                        onChange(value.trimEnd())
                    }}>
                        Simple mode
                    </Button>
                    <textarea
                        className={`form-control ${!!invalid ? "is-invalid" : ""}`}
                        rows={10}
                        placeholder="VAR1=VALUE1"
                        name="environment"
                        onChange={(e) => onChange(e.target.value)}
                        value={value}
                    ></textarea>
                </React.Fragment>
            ) : (
                <React.Fragment>
                    <Button color="link" className="p-0 my-2" onClick={(e) => {
                        e.preventDefault();
                        setAdvancedModeEnabled(true);
                    }}>
                        Advanced mode
                    </Button>
                    {
                        lines.map((line, index) => (
                            <EnvEditorRow
                                key={index}
                                value={line}
                                onChange={(newValue) => {
                                    const newLines = [...lines];
                                    newLines[index] = newValue;
                                    onChange(newLines.join("\n"));
                                }} onDelete={() => {
                                    const newLines = [...lines];
                                    newLines.splice(index, 1);
                                    onChange(newLines.join("\n"));
                                }}
                            />
                        ))
                    }
                    <Button
                        color="accent"
                        className="mt-3"
                        onClick={(e) => {
                            e.preventDefault();
                            onChange(value + "\n");
                        }}
                    >
                        <FontAwesomeIcon
                            icon={faPlus}
                            className="me-2"
                        />
                        Add an environment variable
                    </Button>
                </React.Fragment>
            )
            }
        </React.Fragment>
    )
}

interface EnvEditorRowProps {
    value?: string;
    onChange?: (value: string) => void;
    onDelete?: () => void;
}


export function EnvEditorRow({ value = "", onChange, onDelete = () => { } }: EnvEditorRowProps) {
    const name = value.split("=")[0] || "";
    const val = value.split("=").slice(1).join("=") || "";

    return (
        <React.Fragment>
            <div className="d-flex gap-3 my-2">
                <InputGroup>
                    <InputGroupText>
                        name*
                    </InputGroupText>
                    <Input
                        placeholder="VAR1"
                        value={name}
                        onChange={(e) => onChange?.(`${e.target.value}=${val}`)}
                    />
                </InputGroup>
                <InputGroup>
                    <InputGroupText>
                        value
                    </InputGroupText>
                    <Input
                        placeholder="VALUE1"
                        value={val}
                        onChange={(e) => onChange?.(`${name}=${e.target.value}`)}
                    />
                </InputGroup>
                <Button color="danger" onClick={(e) => {
                    e.preventDefault();
                    onDelete();
                }}>
                    <FontAwesomeIcon icon={faTrashCan} />
                </Button>
            </div>
        </React.Fragment>
    )
}
