//2025-10-21, for now i'll try to avoid using additional libraries such as jquery/ajax, let's stick to vanilla JS.
//TODO: move the DOM elements from go templates into the script. Let handle the form field generation as well as file drop funcionality.

export function init_file_selector()
{
    let fileBlobs = [];
    let fileList = new DataTransfer();

    const fileThumbs = document.getElementById("file-thumbs");

    const dropZone = document.getElementById("dropzone");
    dropZone.addEventListener("drop", dropHandler);

    const fileInput = document.getElementById("file-input");
    fileInput.addEventListener("change", (e) => {
        addUploadsToForm(e.target.files);
    });

    

    //Prevent all drop events except for if the user is dropping a file.
    window.addEventListener("drop", (e) => {
        if ([...e.dataTransfer.items].some((item) => item.kind === "file")) {
            e.preventDefault();
        }
    })

    dropZone.addEventListener("dragover", (e) => {
        const fileItems = [...e.dataTransfer.items].filter(
            (item) => item.kind === "file",
        );
        if (fileItems.length > 0) {
            e.preventDefault();
            if (fileItems.some((item) => item.type.startsWith("image/"))) {
                e.dataTransfer.dropEffect = "copy";
            } else {
                e.dataTransfer.dropEffect = "none";
            }
        }
    });

    window.addEventListener("dragover", (e) => {
        const fileItems = [...e.dataTransfer.items].filter(
            (item) => item.kind === "file",
        );
        if (fileItems.length > 0) {
            e.preventDefault();
            if (!dropZone.contains(e.target)) {
                e.dataTransfer.dropEffect = "none";
            }
        }
    });

    function addUploadsToForm(files){
        let file_idx = fileInput.files.length;

        for (const file of files) {
            if (file.type.startsWith("image/")) {
                const blobUrl = URL.createObjectURL(file);
                fileBlobs.push(blobUrl);

                const thumbContainer = document.createElement("div");
                const fileThumb = document.createElement("div");
                const fileName = document.createElement("div");
                const removeBtn = document.createElement("div");

                thumbContainer.classList.add("tmb-container");
                fileThumb.classList.add("file-tmb");
                fileName.classList.add("tmb-filename");
                removeBtn.classList.add("remove-btn");
                fileThumb.style = "background-image: url(" + blobUrl + ")";

                removeBtn.setAttribute("id", "remove-btn");
                fileThumb.setAttribute("id", file_idx);

                fileName.append(file.name);
                removeBtn.append("✖");

                thumbContainer.appendChild(removeBtn);
                thumbContainer.appendChild(fileThumb);
                thumbContainer.appendChild(fileName);
                fileThumbs.appendChild(thumbContainer);

                fileList.items.add(file);
                file_idx = file_idx + 1;

                removeBtn.addEventListener("click", (e) => {
                    let clicked_file = e.target.parentElement;
                    let file_idx = clicked_file.getElementsByClassName("file-tmb")[0].id;

                    URL.revokeObjectURL(fileBlobs[file_idx]);
                    fileList.items.remove(file_idx);
                    fileInput.files = fileList.files;

                    clicked_file.remove();
                })
            }
        }

        fileInput.files = fileList.files;
    }

    function dropHandler(ev) {
        ev.preventDefault();
        const files = [...ev.dataTransfer.items]
            .map((item) => item.getAsFile())
            .filter((file) => file);
        addUploadsToForm(files);
    }
}
