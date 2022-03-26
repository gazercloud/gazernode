export default async function Request(func, data = {}) {
    const formData = new FormData();
    formData.append('fn', func);
    formData.append('rj', JSON.stringify(data));

    // Default options are marked with *
    return await fetch('/api/request', {
        method: 'POST', // *GET, POST, PUT, DELETE, etc.
        mode: 'cors', // no-cors, *cors, same-origin
        cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
        credentials: 'same-origin', // include, *same-origin, omit
        headers: {
            //'Content-Type': 'multipart/form-data; boundary=CUSTOM'
            // 'Content-Type': 'application/x-www-form-urlencoded',
        },
        redirect: 'follow', // manual, *follow, error
        referrerPolicy: 'no-referrer', // no-referrer, *client
        body: formData // body data type must match "Content-Type" header
    }); // parses JSON response into native JavaScript objects
}

export function RequestFailed() {
    console.log("RequestFailed")
}
