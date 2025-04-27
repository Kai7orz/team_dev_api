const params = new URLSearchParams(window.location.search);
let currentPage = params.get("page") === null ? 1 : parseInt(params.get("page"));
let filter = params.get("culture") === null ? "" : params.get("culture");
let sortBy = params.get("sortBy") === null ? "" : params.get("sortBy");

const config = {
    url: `http://localhost:8080/artworks`,
    parentId: "item-list",
    seeMoreButtonId: "loadMore",
    searchButton: "search-button",
    searchField: "search-field",
    pageID: "page-id",
    topPageButton: "top-page-button",
    sortButton: "sort-button"
}

fetchArtworks();
pagination();

document.getElementById(config.sortButton).addEventListener("click", () => {
    const selectedSort = document.getElementById("sort-select").value;
    let url = `/demo/index.html?page=${currentPage}`;


    if (params.has("culture")) {
        url += `&culture=${params.get("culture")}`;
    }
    url += `&sortBy=${selectedSort}`;

    location.href = url;
});

document.getElementById(config.searchButton).addEventListener("click", () => {
    let value = document.getElementById(config.searchField).value;
    value = value.replace(/\s+/g, '');

    if(value !== "") {
        currentPage = 1;
        location.href=`index.html?page=${currentPage}&culture=${value}`;
    }

    if(sortBy !== "") {
        location.href = `index.html?page=${currentPage}&culture=${value}&sortBy=${sortBy}`;
    }

});

function fetchArtworks() {
    let url = config.url + "?page=" + currentPage;
    if(filter !== "")  url += "&culture=" + filter;
    if(sortBy !== "") url += "&sortBy=" + sortBy;


    console.log(url);
    fetch(url).then(response=>response.json()).then(function(data){
        data.forEach(element => {
            generateArtCard(element);
            console.log(element);
        });
    });

}

function generateArtCard(art) {
    let cardDiv = document.createElement("div");
    let detailParams = `/demo/details.html?id=${art.objectID}&page=${currentPage}`
    const urlParams = new URLSearchParams(window.location.search);

    if (urlParams.has("culture")) {
        detailParams += `&culture=${urlParams.get("culture")}`;
    }
    if (urlParams.has("sortBy")) {
        detailParams += `&sortBy=${urlParams.get("sortBy")}`;
    }

    if(art.artistDisplayName === "") {
        art.artistDisplayName = "\"the artist is unknown\"";
    }

    cardDiv.innerHTML =
    `
        <button class="art-button mb-2" onclick="location.href='${detailParams}'"><strong>${art.title}</strong> by ${art.artistDisplayName}</button>
    `
    document.getElementById(config.parentId).append(cardDiv);
}

function pagination(){
    let defaultParams1 = `/demo/index.html?page=${currentPage}`;
    let defaultParams2 = `/demo/index.html?page=${currentPage+1}`;
    let defaultParams3 = `/demo/index.html?page=${currentPage+2}`;
    if(filter !== "") {
        defaultParams1 += `&culture=${filter}`;
        defaultParams2 += `&culture=${filter}`;
        defaultParams3 += `&culture=${filter}`;
    }else {
        const urlParams = new URLSearchParams(window.location.search);
        if (urlParams.has("culture")) {
            defaultParams1 += `&culture=${urlParams.get("culture")}`;
            defaultParams2 += `&culture=${urlParams.get("culture")}`;
            defaultParams3 += `&culture=${urlParams.get("culture")}`;
        }
        if (urlParams.has("sortBy")) {
            defaultParams1 += `&sortBy=${urlParams.get("sortBy")}`;
            defaultParams2 += `&sortBy=${urlParams.get("sortBy")}`;
            defaultParams3 += `&sortBy=${urlParams.get("sortBy")}`;
        }
    }

    let parentNav = document.getElementById(config.pageID);
    parentNav.innerHTML =
    `
        <ul class="pagination pagination-lg justify-content-center mt-5">
            <li id="page-back-button" class="page-item">
                <a class="page-link" aria-label="Previous">
                    <span aria-hidden="true">&laquo;</span>
                </a>
            </li>
            <li class="page-item"><a class="page-link" href=${defaultParams1}>${currentPage}</a></li>
            <li class="page-item"><a class="page-link" href=${defaultParams2}>${currentPage+1}</a></li>
            <li class="page-item"><a class="page-link" href=${defaultParams3}>${currentPage+2}</a></li>
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

document.getElementById(config.topPageButton).addEventListener("click", () => {
    currentPage = 1;
    location.href='index.html';
})