
pipDoctor — Python Environment Repair Tool
==========================================

**Version:** 1.0  
**Author:** MyArchiveProjects  
**GitHub:** https://github.com/MyArchiveProjects/pip-doctor

Overview
--------

pipDoctor is a serious, professional CLI tool written in Go designed to automatically detect and fix broken or missing pip and Python environments. It is intended for beginners and developers who encounter problems with Python installations on Windows.

This utility:
- Repairs missing or broken pip installations
- Checks PyPI (Python Package Index) accessibility
- Offers full Python repair via ensurepip
- Helps users recover from corrupted or PATH-less Python setups
- Provides download link to official Python 3.11 if recovery is not possible

Features
--------

- 🛠 **Repair pip** — detects if pip is missing, downloads get-pip.py and reinstalls it
- 🌐 **Check PyPI** — tests connectivity to https://pypi.org/simple/ to verify if packages can be installed
- 🔧 **Full Python Repair** — uses `ensurepip` to rebuild Python's core environment
- 🔍 **Python Auto-Detect** — scans PATH, known folders and optionally full disk to find python.exe
- 📎 **GitHub Integrated** — links to the official project and Python download
- 🪪 **No dependencies** — fully portable `.exe`

How to Use
----------

1. **Run** `pipDoctor.exe`
2. Choose an option:

```
1. Repair pip installation
2. Check network access (PyPI)
3. Repair full Python (beta)
4. Exit
```

3. Follow prompts if Python path is required

Build from Source
-----------------

```bash
go install github.com/akavel/rsrc@latest
rsrc -ico pipdoctor.ico -o rsrc.syso
go build -ldflags="-H windowsgui" -o pipDoctor.exe main.go
```

You can also build without icon:

```bash
go build -o pipDoctor.exe main.go
```

Icon
----

Custom icon provided. Convert PNG to ICO via https://icoconvert.com and use with `rsrc`.

License
-------

MIT — Use at your own risk.
