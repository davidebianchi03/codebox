const darkThemeCdnUrl = "https://cdn.jsdelivr.net/npm/swagger-ui-themes@3.0.0/themes/3.x/theme-monokai.css";

function setSwaggerTheme() {
    const isDark = document.body.getAttribute("data-theme") === "dark";
    let darkCss = document.getElementById("swagger-dark-theme");

    if (!darkCss) {
        darkCss = document.createElement("link");
        darkCss.id = "swagger-dark-theme";
        darkCss.rel = "stylesheet";
        darkCss.href = darkThemeCdnUrl;
        darkCss.disabled = true; // start disabled
        document.head.appendChild(darkCss);
    }

    darkCss.disabled = !isDark;
}

window.onload = function () {
    // Applica tema prima di inizializzare Swagger
    setSwaggerTheme();

    SwaggerUIBundle({
        url: "../../_static/swagger.json",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
            SwaggerUIBundle.presets.apis,
            SwaggerUIBundle.SwaggerUIStandalonePreset
        ],
        layout: "BaseLayout"
    });

    // Osserva cambiamenti tema
    // const observer = new MutationObserver(setSwaggerTheme);
    // observer.observe(document.body, {
    //     attributes: true,
    //     attributeFilter: ["data-theme"],
    // });
};
