import os
# Configuration file for the Sphinx documentation builder.
#
# For the full list of built-in configuration values, see the documentation:
# https://www.sphinx-doc.org/en/master/usage/configuration.html

# -- Project information -----------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#project-information

project = 'codebox'
copyright = '2025, Davide Bianchi'
author = 'Davide Bianchi'
version = os.getenv("CI_COMMIT_TAG", "dbg-v1.0.0")
server_base_url = os.getenv("CI_PAGES_URL", "http://127.0.0.1:8000")

# -- General configuration ---------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#general-configuration

extensions = [
    'myst_parser',
    # 'sphinxcontrib.openapi',
    # 'sphinxcontrib.httpdomain',
    # 'sphinx.ext.extlinks',
    'sphinxcontrib.redoc',
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
redoc = [
    {
        'name': 'Codebox Server API',
        'page': 'api/codebox-server/index',
        'spec': '_specs/swagger.yaml',
        'opts': {
            'lazy-rendering': True
        },
    },
]

myst_enable_extensions = [
    "colon_fence",
    "html_admonition",
    "html_image",
    "substitution",
]

myst_substitutions = {
    "server_base_url": server_base_url,
}
