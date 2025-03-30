package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

const githubLink = "https://github.com/MyArchiveProjects/pip-doctor"

func main() {
	for {
		printMenu()
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			repairPip()
		case "2":
			checkPyPI()
		case "3":
			fullRepairPython()
		case "4":
			fmt.Println("Exiting pipDoctor.")
			fmt.Println("Project: " + githubLink)
			return
		default:
			fmt.Println("Invalid option. Please enter 1, 2, 3 or 4.")
		}
	}
}

func printMenu() {
	fmt.Println("==============================")
	fmt.Println("        pipDoctor v1.0")
	fmt.Println("  Python environment fixer")
	fmt.Println("==============================")
	fmt.Println("GitHub: " + githubLink)
	fmt.Println("")
	fmt.Println("1. Repair pip installation")
	fmt.Println("2. Check network access (PyPI)")
	fmt.Println("3. Repair full Python (beta)")
	fmt.Println("4. Exit")
	fmt.Println("")
	fmt.Print("Select an option: ")
}

func repairPip() {
	fmt.Println("[*] Starting pip repair...")

	pythonPath := findPython()
	if pythonPath == "" {
		pythonPath = promptForPythonPath()
		if pythonPath == "" {
			fmt.Println("[!] Python path not provided. Aborting.")
			return
		}
	}
	fmt.Println("[✓] Python found at:", pythonPath)

	if checkCommand("pip", "--version") {
		fmt.Println("[✓] pip is already installed.")
		return
	}

	fmt.Println("[!] pip not found. Downloading installer...")

	if err := downloadGetPip(); err != nil {
		fmt.Println("[✘] Failed to download get-pip.py:", err)
		return
	}
	fmt.Println("[⇩] get-pip.py downloaded")

	if err := runGetPip(pythonPath); err != nil {
		fmt.Println("[✘] Failed to install pip:", err)
		return
	}
	fmt.Println("[✓] pip installed successfully")

	addToPath(filepath.Dir(pythonPath))

	_ = os.Remove("get-pip.py")
	fmt.Println("[✓] Temporary files cleaned up")
	fmt.Println("[✓] pip repair completed")
}

func checkPyPI() {
	fmt.Println("[*] Checking access to PyPI...")

	resp, err := http.Head("https://pypi.org/simple/")
	if err != nil {
		fmt.Println("[✘] Could not connect to PyPI:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("[✓] PyPI is reachable. No issues detected.")
	} else {
		fmt.Printf("[✘] PyPI returned status code %d. Please check your network or firewall.\n", resp.StatusCode)
	}
}

func fullRepairPython() {
	fmt.Println("[*] Starting full Python repair...")

	pythonPath := findPython()
	if pythonPath == "" {
		fmt.Println("[✘] Python is completely broken or not found.")
		fmt.Println("[→] Please install Python 3.11 manually:")
		fmt.Println("    https://www.python.org/ftp/python/3.11.7/python-3.11.7-amd64.exe")
		fmt.Println("[→] Or visit: " + githubLink)
		return
	}

	fmt.Println("[!] Attempting to reset Python environment...")

	cmd := exec.Command(pythonPath, "-m", "ensurepip", "--upgrade")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("[✘] Python environment is too corrupted to recover.")
		fmt.Println("[→] Please reinstall Python 3.11 from:")
		fmt.Println("    https://www.python.org/ftp/python/3.11.7/python-3.11.7-amd64.exe")
		fmt.Println("[→] Or visit: " + githubLink)
		return
	}

	addToPath(filepath.Dir(pythonPath))
	fmt.Println("[✓] Python core modules repaired.")
	fmt.Println("[✓] pip and base packages reset.")
}

func checkCommand(name string, arg string) bool {
	cmd := exec.Command(name, arg)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run() == nil
}

func downloadGetPip() error {
	url := "https://bootstrap.pypa.io/get-pip.py"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create("get-pip.py")
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func runGetPip(pythonPath string) error {
	cmd := exec.Command(pythonPath, "get-pip.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func addToPath(dir string) {
	path := os.Getenv("PATH")
	if strings.Contains(strings.ToLower(path), strings.ToLower(dir)) {
		fmt.Println("[✓] PATH already contains Python")
		return
	}
	newPath := path + ";" + dir
	_ = syscall.Setenv("PATH", newPath)
	fmt.Println("[✓] Python path added to current session PATH")
}

func promptForPythonPath() string {
	fmt.Print("[-] Python was not found. Please enter the full path to your Python executable:\n> ")
	reader := bufio.NewReader(os.Stdin)
	path, _ := reader.ReadString('\n')
	path = strings.TrimSpace(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("[!] Provided path does not exist.")
		return ""
	}
	return path
}

func findPython() string {
	if checkCommand("python", "--version") {
		pythonPath, _ := exec.LookPath("python")
		return pythonPath
	}

	paths := []string{
		`C:\Python39\python.exe`,
		`C:\Python310\python.exe`,
		`C:\Python311\python.exe`,
		os.Getenv("LocalAppData") + `\Programs\Python\Python39\python.exe`,
		os.Getenv("LocalAppData") + `\Programs\Python\Python310\python.exe`,
		os.Getenv("LocalAppData") + `\Programs\Python\Python311\python.exe`,
		`C:\Program Files\Python39\python.exe`,
		`C:\Program Files\Python310\python.exe`,
		`C:\Program Files\Python311\python.exe`,
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	userDir := os.Getenv("USERPROFILE")
	searchPaths := []string{
		userDir + `\AppData\Local\Programs`,
		userDir + `\AppData\Local`,
		userDir + `\AppData\Roaming`,
		`C:\Program Files`,
		`C:\`,
	}

	var found string
	for _, base := range searchPaths {
		filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
			if err == nil && strings.ToLower(info.Name()) == "python.exe" {
				found = path
				return io.EOF
			}
			return nil
		})
		if found != "" {
			break
		}
	}

	return found
}
