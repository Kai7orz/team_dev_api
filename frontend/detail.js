const params = new URLSearchParams(window.location.search);
const artId = params.get("id");
const currentPage = params.get("page");

parentDiv = document.createElement("div");

fetch(`http://localhost:8080/artworks/${artId}`)
    .then(response => response.json())
    .then(data => {

        for(key in data){
            if(data[key] == "") data[key] = "Unknown";
        }

        parentDiv.innerHTML = `
            <div class="row card-screen">
                <div class="col-7">
                    <img src=${data.primaryImage} class="image w-100" alt="No Image">
                </div>
                <div class="col-5 pl-5 mt-4">
                    
                    <h2>${data.title}</h2>
                    <h5 class="pt-3">Date:  ${data.objectDate}</h5>
                    
                    <h5>Artist: ${data.artistDisplayName}</h5>
                    <h5>Culture: ${data.culture}</h5>
                    <h5>ID: ${data.objectID}</h5>
                </div>
            </div>
        `;
        document.getElementById("art-detail").append(parentDiv);
    });

document.getElementById("back-button").addEventListener("click", () => {
    location.href=`/frontend/index.html?page=${currentPage}`;
});