# SSH Server Selector

This project allows you to select an SSH server and an SSH key using [`fzf`] and then connect to the selected server using the selected key.

## Prerequisites

- Go (Golang) installed
- [`fzf`] installed
- SSH keys stored in the [`authorized_keys/`] directory
- A JSON file named [`servers.json`] with the server details

## JSON File Format

The [`servers.json`] file should be structured as follows:

```json
{
    "servers": {
        "group1": {
            "ssh-servers": [
                "ssh-server1 localhost 2222",
                "ssh-server2 localhost 2223",
                "ssh-server3 localhost 2224"
            ]
        }
    }
}
```

## How to Run

[![asciicast](https://asciinema.org/a/GkPTqrdm7xZw9kwv14JGq9jrm.svg)](https://asciinema.org/a/GkPTqrdm7xZw9kwv14JGq9jrm)

1. Clone the repository and navigate to the project directory.

2. Ensure you have a [`servers.json`] file in the project directory with the correct format.

3. Ensure you have your SSH keys in the [`authorized_keys/`] directory.

4. Build the Go program:

    ```sh
    go build -o ssh-selector main.go
    ```

5. Run the program:

    ```sh
    ./ssh-selector
    ```

## How It Works

1. The program reads the [`servers.json`] file to get the list of servers.
2. It writes the list of servers to a temporary file.
3. It uses [`fzf`] to allow you to select a server from the list.
4. It reads the [`authorized_keys/`] directory to get the list of SSH keys.
5. It writes the list of SSH keys to another temporary file.
6. It uses [`fzf`] to allow you to select an SSH key from the list.
7. It removes the old host key for the selected server.
8. It connects to the selected server using the selected SSH key.

## Example

```sh
./ssh-selector
```

You will be prompted to select a server and then an SSH key. After making your selections, the program will connect to the selected server using the selected SSH key.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Contributing

Feel free to submit issues and pull requests for improvements and bug fixes.