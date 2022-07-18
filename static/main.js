const fileForm = document.querySelector('.file-form');
const fileInput = fileForm.querySelector('input[type="file"]');
const emailInput = document.querySelector('input[type="email"]');
const fileSubmit = fileForm.querySelector('input[type="submit"]');
const metaInfo = document.querySelector('.meta-info');
const framsesInAnimation = 150;
const maxFileSize = 10737418240;
const maxSingleTransferSize = 5000000;
const numberOfSplits = 10;
let animation, multiProgessCounter = 0;

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

window.addEventListener("load", () => {
    animation = bodymovin.loadAnimation({
        container: document.getElementById('lottie'),
        path: '/static/animation.json',
        renderer: 'svg',
        autoplay: false,
    });
});

function checkFileSize(file) {
    if (file.size > maxFileSize) {
        alert('File size is too big');
        return false;
    }
    if (file.name.match(/____/g)) {
        alert('File name is invalid');
        return false;
    }
    return true;
}

fileInput.addEventListener('change', (e) => {
    if (fileInput.files.length > 0) {
        const file = fileInput.files[0];
        if (checkFileSize(file)) {
            addMetaInfo(file);
        } else {
            fileInput.value = '';
        }
    }
});

function handleSingleFileTransfer(file, to) {
    const data = new FormData();
    data.append('file', file);
    data.append('to', to);

    const xhr = new XMLHttpRequest();
    xhr.addEventListener('loadend', () => {
        if (xhr.status === 200) {
            metaInfo.innerHTML = 'File uploaded successfully';
        }
    });
    xhr.upload.onprogress = (event) => {
        if (event.lengthComputable) {
            const percent = Math.round((event.loaded / event.total) * 100);
            metaInfo.innerHTML = `<div class="upload-stat">${percent}%</div>`;
            animation.goToAndStop(framsesInAnimation * percent / 100, true);
        }
    };
    xhr.onerror = (event) => {
        console.log(event);
    }

    xhr.open('POST', '/upload');
    xhr.send(data);
}

function handleMultipleFileTransfer(file, to) {
    const fileReader = new FileReader();
    let buffer;
    fileReader.onload = (e) => {
        buffer = new Uint8Array(e.target.result);
        splitAndSend(buffer, file, to);
    }
    fileReader.readAsArrayBuffer(file);
}

function splitAndSend(buffer, file, to) {
    const stepSize = Math.round(buffer.byteLength / numberOfSplits);
    for (let i = 0; i < numberOfSplits; i++) {
        const start = i * stepSize;
        const end = (i + 1) * stepSize;
        const blob = new Blob([buffer.subarray(start, end)]);
        const data = new FormData();
        data.append('file', blob, i + "____" + file.name);
        data.append('to', to);

        const xhr = new XMLHttpRequest();
        xhr.onprogress = (event) => {
            if (event.lengthComputable) {
                const percent = Math.round((event.loaded / event.total) * 100);
                multiProgessCounter += percent / numberOfSplits;
                metaInfo.innerHTML = `<div class="upload-stat">${multiProgessCounter.toFixed(2)}%</div>`;
                animation.goToAndStop(framsesInAnimation * multiProgessCounter / 100, true);
            }
        };
        xhr.onloadend = (event) => {
            multiProgessCounter = 0;
        };
        xhr.onerror = (event) => {
            multiProgessCounter = 0;
        };
        xhr.open('POST', '/multi');
        xhr.send(data);
    }
}

fileSubmit.addEventListener('click', async (e) => {
    e.preventDefault();

    const to = emailInput.value;
    if (!to || !to.match(/@/)) {
        alert('Email is invalid or missing');
        return;
    }

    if (fileInput.files.length > 0) {
        const file = fileInput.files[0];

        if (file.size <= maxSingleTransferSize) {
            handleSingleFileTransfer(file, to);
        } else {
            handleMultipleFileTransfer(file, to);
        }
    }
});