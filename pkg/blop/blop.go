package blop

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/panthyy/blop/internal/cache"
	"github.com/panthyy/blop/internal/manifest"
	"github.com/panthyy/blop/internal/template"
	"github.com/panthyy/blop/internal/utils"
	"github.com/spf13/cobra"
)

var (
	manifestURL string
	outputDir   string
	cacheDir    string
	localFile   string
	templateID  string
)

func init() {
	homeDir, _ := os.UserHomeDir()
	cacheDir = filepath.Join(homeDir, ".blop", "cache")
}

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "blop",
		Short: "Blop is a project template generator",
		Long:  `Blop helps you quickly scaffold new projects based on customizable templates.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return runUnknownCommand(args[0])
			}
			return cmd.Help()
		},
	}

	rootCmd.AddCommand(newGenCommand())
	rootCmd.AddCommand(newImportCommand())
	rootCmd.AddCommand(newListCommand())
	rootCmd.AddCommand(newRemoveCommand())

	return rootCmd
}

func newGenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen [template-id]",
		Short: "Generate a new project from a template",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runGen,
	}

	cmd.Flags().StringVar(&manifestURL, "manifest", "", "URL of the manifest file")
	cmd.Flags().StringVar(&outputDir, "output", ".", "Output directory for the new project")

	return cmd
}

func newImportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import [URL]",
		Short: "Import a template from a URL or local file",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runImport,
	}

	cmd.Flags().StringVarP(&localFile, "file", "f", "", "Path to a local manifest file")

	return cmd
}

func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all available templates",
		RunE:  runList,
	}
}

func newRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <template-id>",
		Short: "Remove a template from the cache",
		Args:  cobra.ExactArgs(1),
		RunE:  runRemove,
	}
}

func runGen(cmd *cobra.Command, args []string) error {
	if len(args) == 0 && manifestURL == "" {
		return listAndSelectTemplate()
	}

	if len(args) > 0 {
		templateID = args[0]
	}

	var data []byte
	var err error

	if manifestURL != "" {
		data, err = utils.DownloadFile(manifestURL)
		if err != nil {
			return fmt.Errorf("failed to download manifest: %w", err)
		}
	} else {
		c := cache.New(cacheDir)
		data, err = c.Get(templateID)
		if err != nil {
			return fmt.Errorf("failed to get template from cache: %w", err)
		}
	}

	m, err := manifest.Parse(data)
	if err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	variables, err := promptVariables(m.Variables)
	if err != nil {
		return fmt.Errorf("failed to prompt for variables: %w", err)
	}

	if err := template.Render(m, outputDir, variables); err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return nil
}

func runImport(cmd *cobra.Command, args []string) error {
	var data []byte
	var err error

	if localFile != "" {
		data, err = os.ReadFile(localFile)
		if err != nil {
			return fmt.Errorf("failed to read local file: %w", err)
		}
	} else if len(args) > 0 {
		url := args[0]
		data, err = utils.DownloadFile(url)
		if err != nil {
			return fmt.Errorf("failed to download template: %w", err)
		}
	} else {
		return fmt.Errorf("please provide either a URL or a local file using the -f flag")
	}

	m, err := manifest.Parse(data)
	if err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	c := cache.New(cacheDir)
	if err := c.Set(m.ID, data); err != nil {
		return fmt.Errorf("failed to cache template: %w", err)
	}

	fmt.Printf("Template '%s' (ID: %s) imported successfully!\n", m.Name, m.ID)
	return nil
}

func runList(cmd *cobra.Command, args []string) error {
	c := cache.New(cacheDir)

	templates, err := c.List()
	if err != nil {
		return fmt.Errorf("failed to list templates: %w", err)
	}

	if len(templates) == 0 {
		fmt.Println("No templates found.")
		return nil
	}

	fmt.Println("Available templates:")
	for _, t := range templates {
		fmt.Printf("- %s\n", t)
	}

	return nil
}

func runRemove(cmd *cobra.Command, args []string) error {
	templateID := args[0]
	c := cache.New(cacheDir)

	if err := c.Remove(templateID); err != nil {
		return fmt.Errorf("failed to remove template '%s': %w", templateID, err)
	}

	fmt.Printf("Template '%s' removed successfully\n", templateID)
	return nil
}

func promptVariables(variables map[string]manifest.Variable) (map[string]string, error) {
	result := make(map[string]string)

	for name, v := range variables {
		var value string
		var err error

		switch v.Type {
		case "select":
			value, err = promptSelect(name, v)
		default:
			value, err = promptInput(name, v)
		}

		if err != nil {
			return nil, fmt.Errorf("error prompting for %s: %w", name, err)
		}

		result[name] = value
	}

	return result, nil
}

func promptSelect(name string, v manifest.Variable) (string, error) {
	options := make([]string, len(v.Options))
	optionMap := make(map[string]string)
	for i, opt := range v.Options {
		options[i] = opt.Name
		optionMap[opt.Name] = opt.Value
	}

	prompt := promptui.Select{
		Label: fmt.Sprintf("Choose %s", name),
		Items: options,
	}

	_, selected, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return optionMap[selected], nil
}

func promptInput(name string, v manifest.Variable) (string, error) {
	validate := func(input string) error {
		if input == "" && v.Default == "" {
			return fmt.Errorf("%s cannot be empty", name)
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("Enter %s", name),
		Default:  v.Default,
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func listAndSelectTemplate() error {
	c := cache.New(cacheDir)
	templates, err := c.List()
	if err != nil {
		return fmt.Errorf("failed to list templates: %w", err)
	}

	if len(templates) == 0 {
		return fmt.Errorf("no templates found. Please import a template first")
	}

	prompt := promptui.Select{
		Label: "Select a template",
		Items: templates,
	}

	_, selected, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("failed to select template: %w", err)
	}

	templateID = selected
	return runGen(nil, []string{templateID})
}

func runUnknownCommand(commandName string) error {
	c := cache.New(cacheDir)
	data, err := c.Get(commandName)
	if err != nil {
		return fmt.Errorf("unknown command '%s'. Run 'blop --help' for usage", commandName)
	}

	m, err := manifest.Parse(data)
	if err != nil {
		return fmt.Errorf("failed to parse manifest for '%s': %w", commandName, err)
	}

	variables, err := promptVariables(m.Variables)
	if err != nil {
		return fmt.Errorf("failed to prompt for variables: %w", err)
	}

	if err := template.Render(m, outputDir, variables); err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	fmt.Printf("Project generated using template '%s'\n", commandName)
	return nil
}
