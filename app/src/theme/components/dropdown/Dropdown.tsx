import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import "./Dropdown.css"

import React, { useState } from 'react';
import { IconProp } from "@fortawesome/fontawesome-svg-core";

export interface DropdownOption {
    title: string
    icon?: IconProp
}

interface DropdownProps {
    title?: string
    titleIcon?: IconProp
    options: DropdownOption[];
    onSelect: (option: DropdownOption) => void;
    minWidth?: string
    left?: string
}

export const Dropdown: React.FC<DropdownProps> = ({ title, titleIcon, options, onSelect, minWidth, left }) => {
    const [isOpen, setIsOpen] = useState<boolean>(false);

    // Toggle the visibility of the dropdown
    const toggleDropdown = () => {
        setIsOpen(!isOpen);
    };

    // Handle item selection and close dropdown
    const handleSelect = (option: DropdownOption) => {
        onSelect(option);
        setIsOpen(false);
    };

    return (
        <div className="dropdown">
            <button onClick={toggleDropdown} className="dropdown-toggle">
                {titleIcon ?
                    <span style={{ marginRight: title ? "10px" : "auto" }}>
                        <FontAwesomeIcon icon={titleIcon} />
                    </span>
                    :
                    null
                }
                {title}
            </button>
            {isOpen && (
                <ul className="dropdown-menu" style={{ minWidth: minWidth, left: left }}>
                    {options.map((option, index) => (
                        <li key={index} onClick={() => handleSelect(option)} className="dropdown-item">
                            {
                                option.icon ?
                                    <FontAwesomeIcon icon={option.icon} /> :
                                    null
                            }
                            <span style={{ marginLeft: option.icon ? "10px" : undefined }}>
                                {option.title}
                            </span>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
};