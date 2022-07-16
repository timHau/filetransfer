const fileForm = document.querySelector('.file-form');
const fileInput = fileForm.querySelector('input[type="file"]');
const fileSubmit = fileForm.querySelector('input[type="submit"]');
const metaInfo = document.querySelector('.meta-info');

function formatFileSize(bytes) {
    const ending = ['B', 'KB', 'MB', 'GB'];
    let res = bytes;
    let i = 0;
    while (res >= 1024 && i < ending.length) {
        res /= 1024;
        i++;
    }
    return `${res.toFixed(2)} ${ending[i]}`;
}

function addMetaInfo(file) {
    const fileSize = file.size; // in bytes
    const fileSizeFormatted = formatFileSize(fileSize);
    const fileType = file.type;
    const fileName = file.name;
    const fileDate = file.lastModified;
    const fileDateFormatted = new Date(fileDate).toLocaleDateString();

    metaInfo.innerHTML = `
        <p>File name: ${fileName}</p>
        <p>File size: ${fileSizeFormatted}</p>
        <p>File type: ${fileType}</p>
        <p>File date: ${fileDateFormatted}</p>
    `;
}

function checkFileSize(file) {
    if (file.size > 10737418240) {
        alert('File size is too big');
        return false;
    }
    return true;
}

fileInput.addEventListener('change', (e) => {
    if (fileInput.files.length > 0) {
        const file = fileInput.files[0];
        addMetaInfo(file);
    }
});

fileSubmit.addEventListener('click', async (e) => {
    e.preventDefault();

    if (fileInput.files.length > 0) {
        const file = fileInput.files[0];
        if (checkFileSize(file)) {
            const data = new FormData();
            data.append('file', file);
            data.append('fileName', "test");

            await fetch('/upload', {
                method: 'POST',
                body: data,
            });
        }
    }
});