const fileForm = document.querySelector('.file-form');
const fileInput = fileForm.querySelector('input[type="file"]');
const fileLabel = document.querySelector('.file-container label');
const senderEmailInput = document.querySelector('.sender-email');
const recipientEmailInput = document.querySelector('.recipient-email');
const fileSubmit = fileForm.querySelector('input[type="submit"]');
const framesInAnimation = 150;
const maxFileSize = 10737418240;
const maxSingleTransferSize = 5000000;
const numberOfSplits = 10;
let animation;

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

    fileLabel.innerHTML = `
    <div class="meta-info">
        <p>File name: ${fileName}</p>
        <p>File size: ${fileSizeFormatted}</p>
        <p>File type: ${fileType}</p>
        <p>File date: ${fileDateFormatted}</p>
    </div>
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
        const file = fileInput.files[fileInput.files.length - 1];
        if (checkFileSize(file)) {
            addMetaInfo(file);
            fileSubmit.classList.remove('locked');
        } else {
            fileInput.value = '';
        }
    }
});

function handleSingleFileTransfer(file, sender, recipient) {
    const data = new FormData();
    data.append('file', file);
    data.append('sender', sender);
    data.append('recipient', recipient);

    const xhr = new XMLHttpRequest();
    xhr.addEventListener('loadend', () => {
        if (xhr.status === 200) {
            fileLabel.innerHTML = 'File uploaded successfully';
        }
    });
    xhr.upload.onprogress = (event) => {
        if (event.lengthComputable) {
            const percent = Math.round((event.loaded / event.total) * 100);
            fileLabel.innerHTML = `<div class="upload-stat">${percent}%</div>`;
            animation.goToAndStop(framesInAnimation * percent / 100, true);
        }
    };
    xhr.onerror = (event) => {
        console.log(event);
    }

    xhr.open('POST', '/upload');
    xhr.send(data);
}

function handleMultipleFileTransfer(file, sender, recipient) {
    const fileReader = new FileReader();
    let buffer;
    fileReader.onload = (e) => {
        buffer = new Uint8Array(e.target.result);
        splitAndSend(buffer, file, sender, recipient);
    }
    fileReader.readAsArrayBuffer(file);
}

function splitAndSend(buffer, file, sender, recipient) {
    animation.goToAndStop(1, true);
    let multiProgessCounter = 0

    const stepSize = Math.round(buffer.byteLength / numberOfSplits);
    for (let i = 0; i < numberOfSplits; i++) {
        const start = i * stepSize;
        const end = (i + 1) * stepSize;
        const blob = new Blob([buffer.subarray(start, end)]);
        const data = new FormData();
        data.append('file', blob, i + "____" + file.name);
        data.append('sender', sender);
        data.append('recipient', recipient);

        const xhr = new XMLHttpRequest();
        xhr.onprogress = (event) => {
            if (event.lengthComputable) {
                // const percent = Math.round((event.loaded / event.total) * 100);
                // fileLabel.innerHTML = `<div class="upload-stat">${percent}%</div>`;
                // animation.goToAndStop(framesInAnimation * percent/100, true);
            }
        };
        xhr.onloadend = () => {
            if (xhr.status === 200) {
                multiProgessCounter += framesInAnimation / numberOfSplits;
                animation.goToAndStop(multiProgessCounter, true);
            }
        };
        xhr.onerror = () => {
            multiProgessCounter = 0;
        };
        xhr.open('POST', '/multi');
        xhr.send(data);
    }
}

fileSubmit.addEventListener('click', async (e) => {
    e.preventDefault();

    const sender = senderEmailInput.value;
    if (!sender || !sender.match(/@/)) {
        alert('Email is invalid or missing');
        return;
    }

    const recipient = recipientEmailInput.value;
    if (!recipient || !recipient.match(/@/)) {
        alert('Email is invalid or missing');
        return;
    }

    if (fileInput.files.length > 0) {
        const file = fileInput.files[0];

        if (file.size <= maxSingleTransferSize) {
            handleSingleFileTransfer(file, sender, recipient);
        } else {
            handleMultipleFileTransfer(file, sender, recipient);
        }
    }
});