const params = new URLSearchParams(window.location.search);
let currentPage = params.get("page") === null ? 1 : parseInt(params.get("page"));

const config = {
    url: `http://localhost:8080/artworks/`,
    parentId: "item-list",
    seeMoreButtonId: "loadMore",
    searchButton: "search-button",
    searchField: "search-field",
    pageID: "page-id"
}

fetchArtworks();
pagination();


document.getElementById(config.searchButton).addEventListener("click", () => {
    let fieldValue = document.getElementById(config.searchField).value;
    let currentQuery = fieldValue;
});

function fetchArtworks() {
    let url = config.url + currentPage;

    fetch(url).then(response=>response.json()).then(function(data){
        //data.forEach(element => {
            generateArtCard(data);
        //});
    });
}

function generateArtCard(art) {
    let cardDiv = document.createElement("div");

    if(art.artistDisplayName === "") {
        art.artistDisplayName = "\"the artist is unknown\"";
    }

    cardDiv.innerHTML = 
    `
        <button class="art-button mb-2" onclick="location.href='details.html?id=${art.objectID}&page=${currentPage}'"><strong>${art.title}</strong> by ${art.artistDisplayName}</button>
    `
    document.getElementById(config.parentId).append(cardDiv);
}

function pagination(){
    let parentNav = document.getElementById(config.pageID);
    parentNav.innerHTML = 
    `
        <ul class="pagination pagination-lg justify-content-center mt-5">
            <li id="page-back-button" class="page-item">
                <a class="page-link" aria-label="Previous">
                    <span aria-hidden="true">&laquo;</span>
                </a>
            </li>
            <li class="page-item"><a class="page-link" href="index.html?page=${currentPage}">${currentPage}</a></li>
            <li class="page-item"><a class="page-link" href="index.html?page=${currentPage+1}">${currentPage+1}</a></li>
            <li class="page-item"><a class="page-link" href="index.html?page=${currentPage+2}">${currentPage+2}</a></li>
            <li id="page-next-button" class="page-item">
                <a class="page-link" aria-label="Next">
                    <span aria-hidden="true">&raquo;</span>
                </a>
            </li>
        </ul>
    `;

    document.getElementById("page-back-button").addEventListener("click", () => {
        if(currentPage != 1)    currentPage--;
        pagination();
    });

    document.getElementById("page-next-button").addEventListener("click", () => {
        currentPage++;
        pagination();
    })   
}

