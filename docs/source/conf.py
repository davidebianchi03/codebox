import os
import subprocess
# Configuration file for the Sphinx documentation builder.
#
# For the full list of built-in configuration values, see the documentation:
# https://www.sphinx-doc.org/en/master/usage/documentation.html#project-information

# -- Project information -----------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/documentation.html#project-information

project = 'codebox'
copyright = '2025, Davide Bianchi'
author = 'Davide Bianchi'

# Get version from CI_COMMIT_TAG, git tag, or fallback to dbg version
def get_version():
    # Try CI_COMMIT_TAG first (CI/CD environment)
    if ci_tag := os.getenv("CI_COMMIT_TAG"):
        return ci_tag
    
    # Try to get the latest git tag (local development)
    try:
        tag = subprocess.check_output(
            ["git", "describe", "--tags", "--abbrev=0"],
            stderr=subprocess.DEVNULL,
            text=True
        ).strip()
        if tag:
            return tag
    except (subprocess.CalledProcessError, FileNotFoundError):
        pass
    
    # Fallback to debug version
    return "dbg-v1.0.0"

version = get_version()
server_base_url = os.getenv("CI_PAGES_URL", "http://127.0.0.1:8000")

# -- General configuration ---------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#general-configuration

extensions = [
    'myst_parser',
    # 'sphinxcontrib.openapi',
    # 'sphinxcontrib.httpdomain',
    # 'sphinx.ext.extlinks',
    # 'sphinxcontrib.redoc',
]

templates_path = ['_templates']
exclude_patterns = []


# -- Options for HTML output -------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#options-for-html-output
html_theme = "furo"
html_static_path = ["_static"]
html_theme_options = {
    "light_logo": "codebox-logo.png",
    "dark_logo": "codebox-logo-white.png",
}
html_title = f"Codebox {version} documentation"

# The suffix of source filenames.
source_suffix = {
    '.rst': 'restructuredtext',
    '.txt': 'markdown',
    '.md': 'markdown',
}

# The master toctree document.
master_doc = 'index'

# redoc (Open API)
# redoc = [
#     {
#         'name': 'Codebox Server API',
#         'page': 'api/codebox-server/index',
#         'spec': '_specs/swagger.yaml',
#         'opts': {
#             'lazy-rendering': True
#         },
#     },
# ]

myst_enable_extensions = [
    "colon_fence",
    "html_admonition",
    "html_image",
    "substitution",
]

myst_substitutions = {
    "server_base_url": server_base_url,
    "version": version,
}

myst_html_meta = {
    "description": "Documentation for Codebox, self-hosted remote development workspaces",
    "keywords": "codebox, IDE, CLI, API, self-hosted, remote development, workspaces",
    "author": author,
}

html_meta = {
    "viewport": "width=device-width, initial-scale=1.0",
}

# RST substitutions for code blocks
rst_prolog = f"""
.. |version| replace:: {version}
"""

# Hook to replace {version} placeholders in generated HTML files
def post_process_html(app, exception):
    """Replace {version} placeholders with actual version in HTML files"""
    if exception:
        return
    
    from pathlib import Path
    
    build_dir = Path(app.outdir)
    for html_file in build_dir.rglob("*.html"):
        try:
            content = html_file.read_text(encoding='utf-8')
            if "{version}" in content:
                # Replace {version} placeholders with actual version
                new_content = content.replace("{version}", version)
                html_file.write_text(new_content, encoding='utf-8')
        except Exception:
            pass

def setup(app):
    app.connect("build-finished", post_process_html)
