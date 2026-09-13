import json
import os
import subprocess
from html import escape
from pathlib import Path
from urllib.parse import urljoin

# Configuration file for the Sphinx documentation builder.
#
# For the full list of built-in configuration values, see the documentation:
# https://www.sphinx-doc.org/en/master/usage/documentation.html#project-information

# -- Project information -----------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/documentation.html#project-information

project = 'codebox'
copyright = '2025, Davide Bianchi'
author = 'Davide Bianchi'

SITE_NAME = 'Codebox Documentation'
SITE_DESCRIPTION = (
    'Codebox documentation for self-hosted remote development workspaces, ' 
    'Docker and Dev Container setup, runner configuration, CLI usage and API reference.'
)
SITE_KEYWORDS = (
    'codebox, remote development, dev environments, docker, self-hosted, workspaces, ' 
    'runner, cli, api, devcontainers'
)
SITE_URL = os.getenv('CODEBOX_DOCS_URL', 'https://codebox.example.com')

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
# https://www.sphinx-doc.org/en/master/usage/documentation.html#general-configuration

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
html_baseurl = SITE_URL
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
    "description": SITE_DESCRIPTION,
    "keywords": SITE_KEYWORDS,
    "author": author,
    "robots": "index, follow",
}

html_meta = {
    "viewport": "width=device-width, initial-scale=1.0",
    "description": SITE_DESCRIPTION,
    "keywords": SITE_KEYWORDS,
    "robots": "index, follow",
}

# RST substitutions for code blocks
rst_prolog = f"""
.. |version| replace:: {version}
"""


def build_sitemap(app, exception):
    """Build a sitemap.xml file for the generated HTML documentation."""
    if exception:
        return

    build_dir = Path(app.outdir)
    if not build_dir.exists():
        return

    urls = []
    for html_file in build_dir.rglob("*.html"):
        relative_path = html_file.relative_to(build_dir).as_posix()
        if relative_path.startswith("_"):
            continue
        if relative_path in {"search.html", "genindex.html", "sitemap.xml", "robots.txt"}:
            continue
        urls.append(urljoin(SITE_URL.rstrip('/') + '/', relative_path))

    urls = sorted(set(urls))
    url_entries = "\n".join(
        f"  <url><loc>{escape(url, quote=True)}</loc></url>" for url in urls
    )

    sitemap_content = (
        '<?xml version="1.0" encoding="UTF-8"?>\n'
        '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n'
        f'{url_entries}\n'
        '</urlset>\n'
    )
    (build_dir / 'sitemap.xml').write_text(sitemap_content, encoding='utf-8')

    robots_content = (
        'User-agent: *\n'
        'Allow: /\n\n'
        f'Sitemap: {SITE_URL.rstrip("/")}/sitemap.xml\n'
    )
    (build_dir / 'robots.txt').write_text(robots_content, encoding='utf-8')


def inject_seo_metadata(app, exception):
    """Add description, social tags and JSON-LD metadata to generated HTML pages."""
    if exception:
        return

    build_dir = Path(app.outdir)
    if not build_dir.exists():
        return

    for html_file in build_dir.rglob("*.html"):
        try:
            content = html_file.read_text(encoding='utf-8')
        except Exception:
            continue

        if "</head>" not in content:
            continue

        relative_path = html_file.relative_to(build_dir).as_posix()
        page_url = urljoin(SITE_URL.rstrip('/') + '/', relative_path)
        page_title = "Codebox Documentation"

        title_match = __import__('re').search(r'<title>(.*?)</title>', content, flags=__import__('re').IGNORECASE | __import__('re').MULTILINE | __import__('re').DOTALL)
        if title_match:
            page_title = title_match.group(1).strip()

        meta_block = f"""
    <meta name="description" content="{escape(SITE_DESCRIPTION, quote=True)}">
    <meta name="keywords" content="{escape(SITE_KEYWORDS, quote=True)}">
    <meta name="robots" content="index, follow">
    <meta property="og:title" content="{escape(page_title, quote=True)}">
    <meta property="og:description" content="{escape(SITE_DESCRIPTION, quote=True)}">
    <meta property="og:type" content="website">
    <meta property="og:url" content="{escape(page_url, quote=True)}">
    <meta property="og:site_name" content="{escape(SITE_NAME, quote=True)}">
    <meta name="twitter:card" content="summary_large_image">
    <meta name="twitter:title" content="{escape(page_title, quote=True)}">
    <meta name="twitter:description" content="{escape(SITE_DESCRIPTION, quote=True)}">
    <link rel="canonical" href="{escape(page_url, quote=True)}">
    <link rel="sitemap" type="application/xml" title="Sitemap" href="{escape(urljoin(SITE_URL.rstrip('/') + '/', 'sitemap.xml'), quote=True)}">
"""

        jsonld_data = {
            "@context": "https://schema.org",
            "@type": "WebSite",
            "name": SITE_NAME,
            "url": SITE_URL.rstrip('/'),
            "description": SITE_DESCRIPTION,
            "inLanguage": "en",
            "publisher": {
                "@type": "Organization",
                "name": author,
                "url": SITE_URL.rstrip('/'),
            },
        }
        jsonld_block = (
            "\n    <script type=\"application/ld+json\">\n"
            f"{json.dumps(jsonld_data, ensure_ascii=False)}\n"
            "    </script>\n"
        )

        if "meta name=\"description\"" in content:
            content = __import__('re').sub(
                r'<meta name="description" content="[^"]*">',
                '<meta name="description" content="' + escape(SITE_DESCRIPTION, quote=True) + '">',
                content,
                count=1,
            )

        content = content.replace("</head>", meta_block + jsonld_block + "\n</head>", 1)
        html_file.write_text(content, encoding='utf-8')


# Hook to replace {version} placeholders in generated HTML files
def post_process_html(app, exception):
    """Replace {version} placeholders with actual version in HTML files"""
    if exception:
        return

    build_dir = Path(app.outdir)
    for html_file in build_dir.rglob("*.html"):
        try:
            content = html_file.read_text(encoding='utf-8')
            if "{version}" in content:
                new_content = content.replace("{version}", version)
                html_file.write_text(new_content, encoding='utf-8')
        except Exception:
            pass


def setup(app):
    app.connect("build-finished", post_process_html)
    app.connect("build-finished", build_sitemap)
    app.connect("build-finished", inject_seo_metadata)
