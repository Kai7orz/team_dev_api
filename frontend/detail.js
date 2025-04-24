const params = new URLSearchParams(window.location.search);
const artId = params.get("id");
const currentPage = params.get("page");

parentDiv = document.createElement("div");

fetch(`http://localhost:8080/artworks/${artId}`)
    .then(response => response.json())
    .then(data => {
        parentDiv.innerHTML = `
            <div class="row card-screen">
                <div class="col-7">
                    <img src=${data.primaryImage} class="image img-fluid w-100">
                </div>
                <div class="col-5">
                    <div class="d-flex flex-row">
                        <h2>${data.title}</h2>
                        <p>${data.objectDate}</p>
                    </div>
                    <p>Artist: ${data.artistDisplayName}</p>
                    <p>${data.culture}</p>
                    <p>ID: ${data.objectID}</p>
                </div>
            </div>
        `;
        document.getElementById("art-detail").append(parentDiv);
    });

document.getElementById("back-button").addEventListener("click", () => {
    location.href=`index.html?page=${currentPage}`;
});