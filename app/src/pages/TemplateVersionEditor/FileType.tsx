import ConsoleIcon from "../../assets/icons/console.svg";
import DockerIcon from "../../assets/icons/docker.svg";
import JSONIcon from "../../assets/icons/json.svg";
import MarkdownIcon from "../../assets/icons/markdown.svg";
import PowerShellIcon from "../../assets/icons/powershell.svg";
import TerraformIcon from "../../assets/icons/terraform.svg";
import YamlIcon from "../../assets/icons/yaml.svg";
import FileIcon from "../../assets/icons/file.svg";


export interface FileMap {
    extensions: string[]
    icon: string
    language: string
}

const FilesMap: FileMap[] = [
    {
        extensions: ["sh"],
        icon: ConsoleIcon,
        language: "shell",
    },
    {
        extensions: ["Dockerfile", "dockerfile"],
        icon: DockerIcon,
        language: "dockerfile",
    },
    {
        extensions: ["json"],
        icon: JSONIcon,
        language: "json",
    },
    {
        extensions: ["md"],
        icon: MarkdownIcon,
        language: "markdown",
    },
    {
        extensions: ["ps1"],
        icon: PowerShellIcon,
        language: "powershell",
    },
    {
        extensions: ["tf"],
        icon: TerraformIcon,
        language: "terraform",
    },
    {
        extensions: ["docker-compose.yaml", "docker-compose.yml"],
        icon: DockerIcon,
        language: "yaml",
    },
    {
        extensions: ["yaml", "yml"],
        icon: YamlIcon,
        language: "yaml",
    },
]

export function GetTypeForFile(filename: string): FileMap {
    var file: FileMap = {
        extensions: ["yaml", "yml"],
        icon: FileIcon,
        language: "",
    };
    FilesMap.forEach(ft => {
        ft.extensions.forEach(ex => {
            if (filename.endsWith(ex) && file.language === "") {
                file = ft;
            }
        })
    });
    return file;
}
