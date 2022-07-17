const fileForm = document.querySelector('.file-form');
const fileInput = fileForm.querySelector('input[type="file"]');
const emailInput = document.querySelector('input[type="email"]');
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

window.addEventListener("load", () => {
    let animation = bodymovin.loadAnimation({
        container: document.getElementById('lottie'),
        path: '/static/animation.json',
        renderer: 'svg',
        autoplay: false,
    });
    
    console.log(animation)
});

function checkFileSize(file) {
    if (file.size > 10737418240) {
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

fileSubmit.addEventListener('click', async (e) => {
    e.preventDefault();

    const to = emailInput.value;
    if (!to || !to.match(/@/)) {
        alert('Email is invalid or missing');
        return;
    }

    if (fileInput.files.length > 0) {
        const file = fileInput.files[0];
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
                metaInfo.innerHTML = `Uploading ${percent}%`;
            }
        };
        xhr.onerror = (event) => {
            console.log(event);
        }

        xhr.open('POST', '/upload');
        xhr.send(data);
    }
});