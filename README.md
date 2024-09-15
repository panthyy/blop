# Blop - Project Template Generator

Blop is a command-line tool that helps you quickly scaffold new projects based on customizable templates.

## Installation

To install Blop, you need to have Go installed on your system. Then, you can use the following command:

```
go install github.com/panthyy/blop@latest
```

This will install the `blop` binary in your `$GOPATH/bin` directory. Make sure this directory is in your system's PATH.

## Manifest File

The manifest file is a YAML file that defines the generation process. It includes the following sections:

- `id`: A unique identifier for the generation process.
- `name`: The name of the generation process.
- `description`: A brief description of the generation process.
- `variables`: A list of variables that the user can customize.
    - `type`: The type of the variable. 
        - `input`: A text input field.
        - `select`: A dropdown menu.
    - `message`: The message to display to the user.
    - `options`: The options to display in the dropdown menu.
- `files`: A list of files that will be created or modified during the generation process.
    - `path`: The path to the file.
    - `content`: The content of the file.

## Example

```yaml
id: "gp"
name: "Generate Project"
description: "Generate a new project with the selected framework and language"
variables:
  name:
    type: input
    message: "Enter the project name"
files:
  - path: "{{ .name }}/package.json"
    content: |
      {
        "name": "{{ .name }}",
        "version": "1.0.0",
        "description": "{{ .description }}",
        "scripts": {},
      }
```

After you have created a manifest file, you can import it using the following command:

```
blop import -f manifest.yaml
```

This will make the cli remember the manifest and you can use it by running:

```
blop gen gp  or blop gp
```

This will prompt you to enter the values for the variables and generate the project.

## Usage

After you have created a manifest file, you can import it using the following command:

```
blop import -f manifest.yaml
```

This will make the cli remember the manifest and you can use it by running:

```
blop gen gp  or blop gp
```

This will prompt you to enter the values for the variables and generate the project.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.


