#!/bin/bash

version_tag="{{ version }}"

print_logo() {
    BLUE='\033[38;5;25m'
    WHITE='\033[97m'
    RESET='\033[0m'

    clear

    echo -e "${WHITE}"
    cat << "EOF"

   ██████╗ ██████╗ ██████╗ ███████╗██████╗  ██████╗ ██╗  ██╗
  ██╔════╝██╔═══██╗██╔══██╗██╔════╝██╔══██╗██╔═══██╗╚██╗██╔╝
  ██║     ██║   ██║██║  ██║█████╗  ██████╔╝██║   ██║ ╚███╔╝
  ██║     ██║   ██║██║  ██║██╔══╝  ██╔══██╗██║   ██║ ██╔██╗
  ╚██████╗╚██████╔╝██████╔╝███████╗██████╔╝╚██████╔╝██╔╝ ██╗
   ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝╚═════╝  ╚═════╝ ╚═╝  ╚═╝
EOF
    year=$(date +"%Y")
    echo -e "${RESET}"
    echo -e "Copyright (c) $year Davide Bianchi. All rights reserved.\n\n"

    echo -e "Welcome to the Codebox installer!"
    echo -e "This installer will set up Codebox using Docker and Docker Compose."
    echo -e "Please make sure you have Docker and Docker Compose installed and running before proceeding."
    echo -e "For more information, visit the Codebox documentation: https://codebox4073715.gitlab.io/codebox/"
}

get_command() {
    # get command from cli args
    local default_command="start"

    if [ $# -eq 0 ]; then
        echo $default_command
    else
        echo $1
    fi
}

# get the stack name from cli args (stack name is the second argument, default is "codebox")
get_codebox_stack_name() {
    local codebox_stack_name="codebox"
    if [ $# -ge 2 ]; then
        codebox_stack_name=$2
    fi
    echo $codebox_stack_name
}

# checks if docker and docker-compose are installed and if the user has 
# permission to run docker commands
check_docker() {
    # check if Docker is installed
    if ! command -v docker &> /dev/null; then
        echo "Docker is not installed. Please install Docker and try again."
        exit 1
    fi

    # check if user can run Docker without sudo
    if ! docker info &> /dev/null; then
        echo "You do not have permission to run Docker commands. Please add your user to the docker group or run this script with sudo."
        exit 1
    fi
}

# get the path to the .env file for the current stack
get_env_file_path() {
    echo "/etc/codebox/$1.env"
}

# create env file
create_env_file() {
    local env_file_path=$1
    
    # ask the use if they want to create the env file
    read -p "The env file for this stack does not exist. Do you want to create it now? (y/n) " answer
    if [[ "$answer" != "y" ]]; then
        echo "Aborting setup. Please create the env file and run the setup script again."
        exit 1
    fi

    # ask the user for external url
    read -p "Enter the external URL for your Codebox instance (e.g. https://codebox.example.com): " external_url
    
    # check that the external url is valid
    if ! [[ "$external_url" =~ ^https?:// ]]; then
        echo "Invalid URL. Please enter a valid URL starting with http:// or https://"
        exit 1
    fi
    
    read -p "Do you want to use subdomains for workspaces? (y/n) " use_subdomains
    if [[ "$use_subdomains" == "y" ]]; then
        use_subdomains=true
        read -p "Enter the wildcard domain for your workspaces (e.g. codebox.example.com): " wildcard_domain

        # check that the wildcard domain is valid
        if ! [[ "$wildcard_domain" =~ ^[a-zA-Z0-9.-]+$ ]]; then
            echo "Invalid domain. Please enter a valid domain (e.g. codebox.example.com)"
            exit 1
        fi
    else
        use_subdomains=false
    fi

    read -p "Do you want to use the default port (12800) for your Codebox instance? (y/n) " use_default_port
    if [[ "$use_default_port" == "y" ]]; then
        port=12800
    else
        read -p "Enter the port for your Codebox instance: " port

        # check that the port is a valid number between 1 and 65535
        if ! [[ "$port" =~ ^[0-9]+$ ]] || [ "$port" -lt 1 ] || [ "$port" -gt 65535 ]; then
            echo "Invalid port. Please enter a valid port number between 1 and 65535."
            exit 1
        fi
    fi

    mkdir -p "$(dirname "$env_file_path")"
    cat << EOF > "$env_file_path"
# Codebox environment variables
# You can customize these variables to configure your Codebox instance
CODEBOX_EXTERNAL_URL=$external_url
CODEBOX_USE_SUBDOMAINS=$use_subdomains
CODEBOX_WILDCARD_DOMAIN=$wildcard_domain
CODEBOX_PORT=$port
EOF
}

# stop the stack and remove containers
stop_stack() {
    local stack_name=$1
    echo "Stopping stack $stack_name..."
    docker compose -p "$stack_name" down
}

# checks if stack is running
is_stack_running() {
    local stack_name=$1
    if docker compose -p "$stack_name" ps -q | grep -q .; then
        return 0
    else
        return 1
    fi
}

# print help message
print_help() {
    echo "Usage: $0 [command] [stack_name]"
    echo "Commands:"
    echo "  setup   - create the .env file for the stack (default command)"
    echo "  start   - start the stack"
    echo "  stop    - stop the stack"
    echo "  help    - print this help message"
    echo "Stack name is optional and defaults to 'codebox'. You can use different stack names to run multiple instances of Codebox on the same machine."
}

# sequentially execute the setup steps
print_logo
check_docker

if [ "$version_tag" = "{{ version }}" ]; then
    read -p "Enter the version tag of Codebox to install (e.g. v0.0.58): " version_tag
fi

command=$(get_command "$@")
stack_name=$(get_codebox_stack_name "$@")
env_file_path=$(get_env_file_path "$stack_name")

echo "Using stack name: $stack_name"
echo "Using config file: $env_file_path"

# setup command
if [ $command == "setup" ]; then
    if [ -f "$env_file_path" ]; then
        echo "Config file already exists at $env_file_path. Do you want to overwrite it? (y/n) "
        read answer
        if [[ "$answer" != "y" ]]; then
            echo "Aborting setup..."
            exit 1
        fi
    fi

    create_env_file "$env_file_path"

    echo "Config file created successfully at $env_file_path. Edit this file to customize your Codebox instance and run the start command to start the stack."

    exit 0
fi

# start command
if [ $command == "start" ]; then
    if [ ! -f "$env_file_path" ]; then
        create_env_file "$env_file_path"

        echo "Config file created successfully at $env_file_path. Edit this file to customize your Codebox instance and run the start command to start the stack."
    fi

    if is_stack_running "$stack_name"; then
        echo "Stack $stack_name is already running. Do you want to stop it and start a new instance? (y/n) "
        read answer
        if [[ "$answer" != "y" ]]; then
            echo "Aborting start..."
            exit 1
        fi

        stop_stack "$stack_name"
    fi

    # download the docker-compose file in a temp dir
    echo "Downloading docker-compose.yml from the Codebox repository..."
    url="https://gitlab.com/codebox4073715/codebox/-/raw/${version_tag}/docker-compose.yml?ref_type=tags"
    temp_dir=$(mktemp -d)
    curl -L "$url" -o "$temp_dir/docker-compose.yml"

    # export env vars
    export $(grep -v '^#' "$env_file_path" | xargs)
    export CODEBOX_VERSION="$version_tag"

    # pull the images
    echo "Pulling Docker images..."
    docker compose -p "$stack_name" -f "$temp_dir/docker-compose.yml" pull

    # TODO: check if config is correct before starting the stack

    # start the stack    
    echo "Starting stack $stack_name..."
    docker compose -p "$stack_name" -f "$temp_dir/docker-compose.yml" up -d

    echo "Stack $stack_name started successfully!"
fi

# stop command
if [ $command == "stop" ]; then
    if ! is_stack_running "$stack_name"; then
        echo "Stack $stack_name is not running."
        exit 1
    fi

    stop_stack "$stack_name"

    echo "Stack $stack_name stopped successfully!"
fi

# help command
if [ $command == "help" ]; then
    echo -e "\nCodebox Installer version $version_tag"
    print_help
    exit 0
fi

echo -e "\nCodebox Installer version $version_tag"
echo "Invalid command: $command"
print_help
exit 1
